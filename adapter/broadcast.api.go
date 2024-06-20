package adapter

import (
	"fmt"
	"tx-monitoring/model"

	"github.com/imroc/req/v3"
)

type broadcastApi struct {
	url string
}

func NewBroadcastApi(url string) *broadcastApi {
	return &broadcastApi{url: url}
}

func (b *broadcastApi) Broadcast(body *model.TxBroadcastReq) (*model.TxBroadcastRes, error) {
	url := fmt.Sprintf("%s/broadcast", b.url)
	result := model.TxBroadcastRes{}

	if res, err := req.C().R().
		SetBody(body).
		SetSuccessResult(&result).
		Post(url); err != nil {
		return nil, err
	} else if res.IsErrorState() {
		return nil, fmt.Errorf("error: %s", res.String())
	}

	return &result, nil
}

func (b *broadcastApi) Check(txHash string) (*model.TxCheckRes, error) {
	url := fmt.Sprintf("%s/check/%s", b.url, txHash)
	result := model.TxCheckRes{}

	if res, err := req.C().R().
		SetSuccessResult(&result).
		Get(url); err != nil {
		return nil, err
	} else if res.IsErrorState() {
		return nil, fmt.Errorf("error: %s", res.String())
	}

	return &result, nil
}
