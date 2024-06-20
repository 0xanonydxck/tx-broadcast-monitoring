package adapter

import (
	"tx-monitoring/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type transactionGorm struct {
	db *gorm.DB
}

func NewTransactionGorm(db *gorm.DB) *transactionGorm {
	return &transactionGorm{db: db}
}

func (t *transactionGorm) Upsert(tx *model.Transaction) error {
	return t.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(tx).Error
}

func (t *transactionGorm) BatchUpsert(txs []model.Transaction) error {
	return t.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&txs).Error
}

func (t *transactionGorm) All() ([]model.Transaction, error) {
	var txs []model.Transaction
	err := t.db.Order("created_at DESC").Order("status DESC").Find(&txs).Error
	return txs, err
}

func (t *transactionGorm) AllPending() ([]model.Transaction, error) {
	var txs []model.Transaction
	err := t.db.Where("status = ?", model.TxStatusPending).Find(&txs).Error
	return txs, err
}
