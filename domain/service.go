package domain

import (
	"tx-monitoring/model"

	"github.com/rs/zerolog/log"
)

type service struct {
	transactionRepo TransactionRepo
	binanceApi      BinanceApi
	broadcastApi    BroadcastApi
}

func NewService(
	transactionRepo TransactionRepo,
	binanceApi BinanceApi,
	broadcastApi BroadcastApi,
) *service {
	return &service{
		transactionRepo: transactionRepo,
		binanceApi:      binanceApi,
		broadcastApi:    broadcastApi,
	}
}

func (s *service) Symbols() ([]model.BinanceSymbol, error) {
	symbols := []string{"BTCUSDT", "ETHUSDT", "BNBUSDT", "SOLUSDT"}
	result, err := s.binanceApi.GetSymbols(symbols)
	if err != nil {
		log.Error().Err(err).Msg("Service::Symbols(): failed to get symbols")
		return nil, err
	}

	return result, err
}

func (s *service) Transactions() ([]model.Transaction, error) {
	txs, err := s.transactionRepo.All()
	if err != nil {
		log.Error().Err(err).Msg("Service::Transactions(): failed to get transactions")
		return nil, err
	}

	return txs, err
}

func (s *service) Broadcast(body *model.TxBroadcastReq) (*model.TxBroadcastRes, error) {
	result, err := s.broadcastApi.Broadcast(body)
	if err != nil {
		log.Error().Err(err).Msg("Service::Broadcast(): failed to broadcast")
		return nil, err
	}

	tx := model.Transaction{
		TxHash: result.TxHash,
		Status: model.TxStatusPending,
	}

	if err := s.transactionRepo.Upsert(&tx); err != nil {
		log.Error().Err(err).Str("txHash", tx.TxHash).Msg("Service::Broadcast(): failed to upsert transaction")
		return nil, err
	}

	log.Info().Str("txHash", tx.TxHash).Msg("Service::Broadcast(): broadcast successfully")
	return result, err
}

func (s *service) CheckTxPending() error {
	txs, err := s.transactionRepo.AllPending()
	if err != nil {
		log.Error().Err(err).Msg("Service::CheckTxPending(): failed to get pending transactions")
		return err
	} else if len(txs) == 0 {
		return nil
	}

	for i, tx := range txs {
		result, err := s.broadcastApi.Check(tx.TxHash)
		if err != nil {
			log.Error().Err(err).Str("txHash", tx.TxHash).Msg("Service::CheckTxPending(): failed to check transaction")
			continue
		}

		if result.Status != model.TxStatusPending {
			txs[i].Status = result.Status
		}
	}

	if err := s.transactionRepo.BatchUpsert(txs); err != nil {
		log.Error().Err(err).Msg("Service::CheckTxPending(): failed to update transactions")
		return err
	}

	return nil
}
