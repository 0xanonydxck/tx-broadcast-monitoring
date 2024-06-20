package model

import (
	"time"
)

type TransactionStatus string

const (
	TxStatusPending  TransactionStatus = "PENDING"
	TxStatusSuccess  TransactionStatus = "CONFIRMED"
	TxStatusFailed   TransactionStatus = "FAILED"
	TxStatusNotFound TransactionStatus = "DNE"
)

type Transaction struct {
	TxHash    string            `json:"tx_hash" gorm:"column:tx_hash;primaryKey"`
	Status    TransactionStatus `json:"status" gorm:"column:status"`
	CreatedAt *time.Time        `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt *time.Time        `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

func (t *Transaction) TableName() string {
	return "tbl_transactions"
}
