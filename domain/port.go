package domain

import "tx-monitoring/model"

type TransactionRepo interface {
	Upsert(tx *model.Transaction) error
	BatchUpsert(txs []model.Transaction) error
	All() ([]model.Transaction, error)
	AllPending() ([]model.Transaction, error)
}

type BinanceApi interface {
	GetSymbols([]string) ([]model.BinanceSymbol, error)
}

type BroadcastApi interface {
	Broadcast(body *model.TxBroadcastReq) (*model.TxBroadcastRes, error)
	Check(txHash string) (*model.TxCheckRes, error)
}
