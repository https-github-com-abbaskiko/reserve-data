package exchange

import (
	"github.com/KyberNetwork/reserve-data/common"
	"github.com/KyberNetwork/reserve-data/settings"
)

type Setting interface {
	GetInternalTokenByID(tokenID string) (common.Token, error)
	GetAllTokens() ([]common.Token, error)
	GetTokenByID(tokenID string) (common.Token, error)
	GetFee(ex settings.ExchangeName) (common.ExchangeFees, error)
	GetMinDeposit(ex settings.ExchangeName) (common.ExchangesMinDeposit, error)
	GetDepositAddresses(ex settings.ExchangeName) (common.ExchangeAddresses, error)
	UpdateDepositAddress(name settings.ExchangeName, addrs common.ExchangeAddresses, timestamp uint64) error
	GetExchangeInfo(ex settings.ExchangeName) (common.ExchangeInfo, error)
	UpdateExchangeInfo(ex settings.ExchangeName, exInfo common.ExchangeInfo, timestamp uint64) error
	ETHToken() common.Token
}