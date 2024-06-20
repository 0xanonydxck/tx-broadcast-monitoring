package adapter

import (
	"fmt"
	"strings"
	"tx-monitoring/model"

	"github.com/imroc/req/v3"
)

type binanceApi struct {
	url       string
	apiKey    string
	secretKey string
}

func NewBinanceApi(url, apiKey, secretKey string) *binanceApi {
	return &binanceApi{
		url:       url,
		apiKey:    apiKey,
		secretKey: secretKey,
	}
}

func (b *binanceApi) GetSymbols(symbols []string) ([]model.BinanceSymbol, error) {
	url := fmt.Sprintf("%s/v3/ticker/price", b.url)
	symbolQuery := fmt.Sprintf("[\"%v\"]", strings.Join(symbols, "\",\""))
	result := []model.BinanceSymbol{}

	if res, err := req.C().R().
		SetQueryParam("symbols", symbolQuery).
		SetHeader("apiKey", b.apiKey).
		SetHeader("secretKey", b.secretKey).
		SetSuccessResult(&result).Get(url); err != nil {
		return nil, err
	} else if res.IsErrorState() {
		return nil, fmt.Errorf("error: %s", res.Status)
	}

	return result, nil
}
