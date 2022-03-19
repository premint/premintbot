package premint

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/kr/pretty"
	"github.com/mager/premintbot/config"
)

type Data struct {
	WalletAddress string `json:"wallet_address"`
}

type GetWalletAddressesResp struct {
	Data []Data `json:"data"`
}

type CheckPremintStatusResp struct {
	Registered    bool   `json:"registered"`
	DiscordID     int    `json:"discord_id"`
	WalletAddress string `json:"wallet_address"`
	ProjectName   string `json:"project_name"`
	ProjectURL    string `json:"project_url"`
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

func CheckPremintStatusForUser(apiKey, userID string) (CheckPremintStatusResp, error) {
	r := CheckPremintStatusResp{}
	url := fmt.Sprintf("https://www.premint.xyz/api/%s/discord/%s", apiKey, userID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
		return r, nil
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
		return r, nil
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		log.Fatal(err)
		return r, nil
	}

	return r, nil
}

func CheckPremintStatusForAddress(apiKey, address string) (CheckPremintStatusResp, error) {
	r := CheckPremintStatusResp{}
	url := fmt.Sprintf("https://www.premint.xyz/api/%s/wallet/%s", apiKey, address)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
		return r, nil
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
		return r, nil
	}
	pretty.Print(url)

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		log.Fatal(err)
		return r, nil
	}

	return r, nil
}
