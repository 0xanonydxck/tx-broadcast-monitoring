package main

import (
	"tx-monitoring/adapter"
	"tx-monitoring/config"
	"tx-monitoring/domain"
	"tx-monitoring/handler"
	"tx-monitoring/infra"
	"tx-monitoring/logger"
)

func init() {
	logger.Init()
	config.Init()
	infra.InitDB()
}

func main() {
	conf := config.Get()
	binanceConf := conf.Binance
	broadcastConf := conf.Broadcast

	hdl := handler.NewHandler(
		domain.NewService(
			adapter.NewTransactionGorm(infra.GetDB()),
			adapter.NewBinanceApi(binanceConf.URL, binanceConf.ApiKey, binanceConf.SecretKey),
			adapter.NewBroadcastApi(broadcastConf.URL),
		),
	)

	hdl.CheckTransaction()
	hdl.Run()
}
