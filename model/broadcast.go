package model

type TxBroadcastReq struct {
	Symbol    string `json:"symbol"`
	Price     uint64 `json:"price"`
	Timestamp int64  `json:"timestamp"`
}

type TxBroadcastRes struct {
	TxHash string `json:"tx_hash"`
}

type TxCheckRes struct {
	Status TransactionStatus `json:"tx_status"`
}
