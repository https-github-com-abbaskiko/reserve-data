package http

import (
	"encoding/json"
	"fmt"

	"github.com/KyberNetwork/reserve-data/common"
	"github.com/KyberNetwork/reserve-data/http/httputil"
	"github.com/gin-gonic/gin"
)

// GetPWIEquationV2 returns the current PWI equations.
func (s *Server) GetPWIEquationV2(c *gin.Context) {
	_, ok := s.Authenticated(c, []string{}, []Permission{ReadOnlyPermission, RebalancePermission, ConfigurePermission, ConfirmConfPermission})
	if !ok {
		return
	}
	data, err := s.metric.GetPWIEquationV2()
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	httputil.ResponseSuccess(c, httputil.WithData(data))
}

// SetPWIEquationV2 stores the given PWI equations to pending for later evaluation.
func (s *Server) SetPWIEquationV2(c *gin.Context) {
	const dataPostFormKey = "data"

	postForm, ok := s.Authenticated(c, []string{dataPostFormKey}, []Permission{ConfigurePermission})
	if !ok {
		return
	}

	data := []byte(postForm.Get(dataPostFormKey))
	if len(data) > maxDataSize {
		httputil.ResponseFailure(c, httputil.WithError(errDataSizeExceed))
		return
	}

	var input common.PWIEquationRequestV2
	if err := json.Unmarshal(data, &input); err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	for tokenID := range input {
		if _, err := s.setting.GetInternalTokenByID(tokenID); err != nil {
			httputil.ResponseFailure(c, httputil.WithReason(fmt.Sprintf("Token %s is unsupported", tokenID)))
		}
	}

	if err := s.metric.StorePendingPWIEquationV2(data); err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	httputil.ResponseSuccess(c)
}

// GetPendingPWIEquationV2 returns the pending PWI equations.
func (s *Server) GetPendingPWIEquationV2(c *gin.Context) {
	_, ok := s.Authenticated(c, []string{}, []Permission{ReadOnlyPermission, RebalancePermission, ConfigurePermission, ConfirmConfPermission})
	if !ok {
		return
	}
	data, err := s.metric.GetPendingPWIEquationV2()
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	httputil.ResponseSuccess(c, httputil.WithData(data))
}

// ConfirmPWIEquationV2 accepts the pending PWI equations and remove it from pending bucket.
func (s *Server) ConfirmPWIEquationV2(c *gin.Context) {
	const dataPostFormKey = "data"

	postForm, ok := s.Authenticated(c, []string{dataPostFormKey}, []Permission{ConfirmConfPermission})
	if !ok {
		return
	}
	postData := postForm.Get(dataPostFormKey)
	err := s.metric.StorePWIEquationV2(postData)
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	httputil.ResponseSuccess(c)
}

// RejectPWIEquationV2 rejects the PWI equations request and removes
// it from pending storage.
func (s *Server) RejectPWIEquationV2(c *gin.Context) {
	_, ok := s.Authenticated(c, []string{}, []Permission{ConfirmConfPermission})
	if !ok {
		return
	}

	if err := s.metric.RemovePendingPWIEquationV2(); err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	httputil.ResponseSuccess(c)
}