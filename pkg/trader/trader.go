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

type Bid struct {
	Currency  string `json:currency`
	PriceUnit string `json:price_unit`
	MinTotal  int    `json:min_total`
}

type Ask struct {
	Currency  string `json:currency`
	PriceUnit string `json:price_unit`
	MinTotal  int    `json:min_total`
	MaxTotal  int    `json:max_total`
}

type BidAccount struct {
	Currency            string `json:currency`
	Balance             string `json:balance`
	Locked              string `json:locked`
	AvgBuyPrice         string `json:avg_buy_price`
	AvgBuyPriceModified bool   `json:avg_buy_price_modified`
	UnitCurrency        string `json:unit_currency`
}

type AskAccount struct {
	Currency            string `json:currency`
	Balance             string `json:balance`
	Locked              string `json:locked`
	AvgBuyPrice         string `json:avg_buy_price`
	AvgBuyPriceModified bool   `json:avg_buy_price_modified`
	UnitCurrency        string `json:unit_currency`
}

type Market struct {
	Id         string     `json:id`
	Name       string     `json:name`
	OrderTypes string[]   `json:order_types`
	OrderSides string[]   `json:order_sides`
	Bid        Bid        `json:bid`
	BidAccount BidAccount `json:bid_account`
	AskAccount AskAccount `json:ask_account`
}

type OrderChance struct {
	BidFee string `json:bid_fee`
	AskFee string `json:ask_fee`
	Market Market `json:market`
	Ask    Ask    `json:ask`
	State  string `json:state`
}

func (t *Trader) Buy() OrderChance {
	req, err := newHttpRequest(http.MethodGet, "orders/chance")
	if err != nil {
		log.Errorf("Occurred an erro when create a new http request to get an info of orders change, %e", err)
	}
	res, _ := t.httpClient.Do(req)
	defer req.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	var orderChance OrderChance

	json.Unmarshal(body, &orderChance)
	return orderChance
}
