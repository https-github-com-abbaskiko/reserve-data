package settings

import (
	"log"
	"os"
)

const (
	tokenDBFileName             string = "token.db"
	TOKEN_DEFAULT_JSON_PATH     string = "/go/src/github.com/KyberNetwork/reserve-data/cmd/mainnet_setting.json"
	TOKEN_DEFAULT_JSON_SIM_PATH string = "/go/src/github.com/KyberNetwork/reserve-data/cmd/shared/deployment_dev.json"
)

type Settings struct {
	Tokens *TokenSetting
}

var setting Settings

func NewSetting(tokenSetting *TokenSetting) *Settings {
	setting = Settings{tokenSetting}
	return &setting
}

func (setting *Settings) HandleEmptyToken(normalPath, simPath string) {
	allToks, err := setting.GetAllTokens()
	if err != nil || len(allToks) < 1 {
		if err != nil {
			log.Printf("Setting Init: Token DB is faulty (%s), attempt to load token from file", err)
		} else {
			log.Printf("Setting Init: Token DB is empty, attempt to load token from file")
		}
		tokenPath := normalPath
		if os.Getenv("KYBER_ENV") == "simulation" {
			tokenPath = simPath
		}

		if err = setting.LoadTokenFromFile(tokenPath); err != nil {
			log.Printf("Setting Init: Can not load Token from file: %s, Token DB is needed to be updated manually", err)
		}
	}
}
