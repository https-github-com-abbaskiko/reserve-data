package http

import (
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"

	"github.com/KyberNetwork/reserve-data/common"
	"github.com/KyberNetwork/reserve-data/http/httputil"
)

const (
	pricingOPAddressName      = "pricing_operator"
	depositOPAddressName      = "deposit_operator"
	intermediateOPAddressName = "intermediate_operator"
	wrapper                   = "wrapper"
	network                   = "network"
	conversionAddr            = "conversion_rate"
	reserveAddr               = "reserve"
)

// GetAddresses get address config from core
func (s *Server) GetAddresses(c *gin.Context) {
	var (
		addresses = make(map[string]ethereum.Address)
	)
	addresses[pricingOPAddressName] = s.blockchain.GetPricingOPAddress()
	addresses[depositOPAddressName] = s.blockchain.GetDepositOPAddress()
	addresses[intermediateOPAddressName] = s.blockchain.GetIntermediatorOPAddress()
	addresses[wrapper] = s.addrs.Wrapper
	addresses[network] = s.addrs.Proxy
	addresses[conversionAddr] = s.addrs.Pricing
	addresses[reserveAddr] = s.addrs.Reserve
	result := common.NewAddressesResponse(addresses)
	httputil.ResponseSuccess(c, httputil.WithData(result))
}
