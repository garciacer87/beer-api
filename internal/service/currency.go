package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/garciacer87/beer-api/internal/contract"
)

//CurrencyClient abstraction used to implement a currency helper
type CurrencyClient interface {
	GetBoxPrice(beer contract.Beer, to string, quantity int) (float64, error)
}

type currencyClient struct {
	baseURL  string
	apiToken string
}

//NewCurrencyClient retrieves a new currency client object
func NewCurrencyClient(token, baseURL string) CurrencyClient {
	return &currencyClient{apiToken: token, baseURL: baseURL}
}

//GetBoxPrice retrieves conversion rate according to currencies provided, then calculates and returns the total amount according to quantity provided
func (cc *currencyClient) GetBoxPrice(beer contract.Beer, to string, quantity int) (float64, error) {
	if cc.apiToken == "" {
		return 0, errors.New("currency API token not defined in environment var CURRENCY_API_TOKEN")
	}

	//Default value for currency to
	if strings.TrimSpace(to) == "" {
		to = "CLP"
	}

	q := strings.ToUpper(fmt.Sprintf("%s_%s", beer.Currency, to))

	url := fmt.Sprintf("%s?q=%s&compact=ultra&apiKey=%s", cc.baseURL, q, cc.apiToken)

	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}

	if resp.StatusCode == http.StatusOK {
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return 0, fmt.Errorf("error getting response body")
		}

		var convResp map[string]interface{}
		json.Unmarshal(respBody, &convResp)

		convRate, ok := convResp[q].(float64)

		if !ok {
			return 0, errors.New("currency rate not valid")
		}

		result := beer.Price * convRate * float64(quantity)

		return result, nil
	}

	return 0, errors.New("response code was not 200")

}
