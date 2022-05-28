package premint

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.uber.org/zap"
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
func ProvidePremint() *PremintClient {
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

// CheckPremintStatusForUser checks if the user is registered on the Premint list
// Example:
// https://api.premint.xyz/v1/6c808f0c6964bd664854f051b438ca8b3aafec70938e5deca6ab49662a789cac/discord/360541062839926785
func CheckPremintStatusForUser(logger *zap.SugaredLogger, apiKey, userID string) (CheckPremintStatusResp, error) {
	r := CheckPremintStatusResp{}
	url := fmt.Sprintf("https://api.premint.xyz/v1/%s/discord/%s", apiKey, userID)

	logger.Infow("Calling Premint API with discord user ID", "url", url)

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

// CheckPremintStatusForAddress checks if the address is registered on the Premint list
// Example: https://api.premint.xyz/v1/6c808f0c6964bd664854f051b438ca8b3aafec70938e5deca6ab49662a789cac/wallet/0x064DcA21b1377D1655AC3CA3e95282D9494E5611
func CheckPremintStatusForAddress(logger *zap.SugaredLogger, apiKey, address string) (CheckPremintStatusResp, error) {
	r := CheckPremintStatusResp{}
	url := fmt.Sprintf("https://api.premint.xyz/v1/%s/wallet/%s", apiKey, address)

	logger.Infow("Calling Premint API with wallet address", "url", url)

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
