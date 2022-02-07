package trader

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go/v4"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
)

func createAuthToken() string {
	accessKey := fmt.Sprintf("%v", viper.Get("upbit.credentials.access_key"))
	secretKey := fmt.Sprintf("%v", viper.Get("upbit.credentials.secret_key"))

	payload := jwt.MapClaims{
		"access_key": accessKey,
		"nonce":      uuid.New().String(),
	}
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err := claims.SignedString([]byte(secretKey))
	if err != nil {
		log.Errorf("*** Failed to generate jwt token for request to upbit: %e", err)
		panic("Failed to generate jwt token for request to upbit")
	}
	return token
}

func newHttpRequest(method, path string) (*http.Request, error) {
	url := fmt.Sprintf("%v/%v", viper.Get("upbit.url"), path)
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Errorf("Failed to create a new http request: %e", err)
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", createAuthToken()))
	return req, nil
}

type Trader struct {
	httpClient *http.Client
}

func NewTrader() *Trader {
	return &Trader{
		httpClient: &http.Client{},
	}
}

type Account struct {
	Currency               string `json:currency`
	Balance                string `json:balance`
	Locked                 string `json:locked`
	Avg_buy_price          string `json:avg_buy_price`
	Avg_buy_price_modified string `json:avg_buy_price_modified`
	Unit_currency          string `json:unit_currency`
}

type Accounts []Account

func (t *Trader) GetAccounts() Accounts {
	req, err := newHttpRequest(http.MethodGet, "accounts")
	if err != nil {
		log.Errorf("Occurred an error when create a new http request for getting account's data, %e", err)
	}
	res, _ := t.httpClient.Do(req)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var accounts Accounts
	json.Unmarshal(body, &accounts)
	fmt.Println(accounts)
	return accounts
}
