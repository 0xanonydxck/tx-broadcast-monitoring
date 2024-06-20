package domain

import "tx-monitoring/model"

type Service interface {
	Symbols() ([]model.BinanceSymbol, error)
	Transactions() ([]model.Transaction, error)
	Broadcast(body *model.TxBroadcastReq) (*model.TxBroadcastRes, error)
	CheckTxPending() error
}
