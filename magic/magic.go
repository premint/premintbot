package magic

import (
	"time"

	"github.com/magiclabs/magic-admin-go"
	"github.com/magiclabs/magic-admin-go/client"
	"github.com/premint/premintbot/config"
)

func ProvideMagic(cfg config.Config) *client.API {
	cl := magic.NewClientWithRetry(5, time.Second, 10*time.Second)

	return client.New(cfg.MagicSecretKey, cl)
}

var Options = ProvideMagic
