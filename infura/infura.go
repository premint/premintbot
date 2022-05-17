package infura

import (
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/premint/premintbot/config"
	ens "github.com/wealdtech/go-ens/v3"
	"go.uber.org/zap"
)

type InfuraClient struct {
	Client *ethclient.Client
	logger *zap.SugaredLogger
}

// ProvideInfura provides an infura client
func ProvideInfura(cfg config.Config, logger *zap.SugaredLogger) *InfuraClient {
	client, _ := ethclient.Dial(fmt.Sprintf("https://mainnet.infura.io/v3/%s", cfg.InfuraKey))
	return &InfuraClient{
		Client: client,
		logger: logger,
	}
}

var Options = ProvideInfura

func (i *InfuraClient) GetAddressFromENSName(ensName string) string {
	address, err := ens.Resolve(i.Client, ensName)
	if err != nil {
		i.logger.Error(err)
		return ""
	}

	return address.Hex()
}
