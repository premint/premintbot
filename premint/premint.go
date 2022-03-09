package premint

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/mager/premintbot/config"
)

type Data struct {
	WalletAddress string `json:"wallet_address"`
}

type GetWalletAddressesResp struct {
	Data []Data `json:"data"`
}

type CheckPremintStatusResp struct {
	Registered bool `json:"registered"`
}

type PremintClient struct {
	httpClient *http.Client
	APIKey     string
}

// ProvidePremint provides an HTTP client
func ProvidePremint(cfg config.Config) *PremintClient {
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}

	return &PremintClient{
		httpClient: &http.Client{
			Transport: tr,
		},
	}
}

var Options = ProvidePremint

func CheckPremintStatus(apiKey, userID string) (bool, error) {
	url := fmt.Sprintf("https://www.premint.xyz/api/%s/entry/%s", apiKey, userID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
		return false, nil
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
		return false, nil
	}

	defer resp.Body.Close()

	var (
		r CheckPremintStatusResp
	)

	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		log.Fatal(err)
		return false, nil
	}

	return r.Registered, nil
}
