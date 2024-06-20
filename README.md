# Transaction Broadcasting and Monitoring Client
![demo gif](/asset/demo.gif)

## Overview
The transaction broadcasting and monitoring client which build into console application which separate following hexagonal architecture, which benefit for make modifying and testing easier.

## Features
- **Broadcast**: send HTTP request with `POST` method for make a transaction with the server and store the `TxHash` into `sqlite` database.
- **Show Transaction**: query and show the `Transaction History` with the transaction status(`CONFIRMED`,`PENDING`,`FAILED` and `DNE`)
- **Check Transaction**: polling the `Transaction Status` from server and update into database

## Installation
1. Clone "Transaction Broadcast and Monitoring Client" repository project into your machine.
```bash
git clone https://github.com/dxckboi/tx-broadcast-monitoring.git
```

2. configuration the config file or clone the `example.yaml` to `config.yaml`.
```yaml
binance:
  url: https://testnet.binance.vision/api # binance api url
  api_key: KsjQtk5UH6nXIzaZo5mYORAgLbPXxJgr40pQ¡h×9HdeSN7fnShfYk4yFcqHC1bUo # binance api key
  secret_key: nT04uwYxeRKXCNB3MTQm6AIbS71H0TFEtlAwIpHMWOJM7MMyzT31bGnkTUKmST # binance secret key

broadcast:
  url: https://mock-node-wgqbnxruha-as.a.run.app # broadcast url
  check_status_interval: 10 # seconds
```

3. Then run the following command for install `Go` dependencies.
```bash
go mod download
```

4. Finally run the following command for start the application.
```bash
go run main.go
```

## Usage
Here’s a simple example script that demonstrates how to use the Client Module to handle a transaction:

### Handling Transaction Statuses
The methodology for checking/handling the transaction statuses, I used `Long-Running Operations` in the client site for polling the result from the `Broadcast Server` every `10` seconds by using `github.com/go-co-op/gocron` dependency for schedule the running task (you can config the polling range in config field `check_status_interval`) so you can see the handle method named `func (h *handler) CheckTransaction()`.

```go
type handler struct {
	service   domain.Service
	validator *validator.Validate
}

func NewHandler(service domain.Service) *handler {
	return &handler{
		service:   service,
		validator: validator.New(),
	}
}

func (h *handler) Run() {
	for {
		h.AskMenu()
	}
}

func (h *handler) AskMenu() {
	const (
		BROADCAST         = "Broadcast"
		SHOW_TRANSACTIONS = "Show Transactions"
		QUIT              = "Quit"
	)

	menuItems := []string{
		BROADCAST,
		SHOW_TRANSACTIONS,
		QUIT,
	}

	answer, err := prompt.New().Ask("Menu").Choose(menuItems)
	if err != nil {
		log.Panic().Err(err).Msg("Handler::AskMenu(): failed to retrieve answer")
	}

	switch answer {
	case BROADCAST:
		h.Broadcast()
	case SHOW_TRANSACTIONS:
		h.ShowTransactions()
	case QUIT:
		h.Quit()
	default:
		log.Error().Msg("Handler::AskMenu(): invalid answer")
		h.AskMenu()
	}
}

/* ...other function... */

func (h *handler) CheckTransaction() {
	conf := config.Get().Broadcast
	scheduler := gocron.NewScheduler(time.Local)
	scheduler.Every(conf.CheckStatusInterval).Seconds().Do(h.service.CheckTxPending)
	scheduler.StartAsync()
}

func (h *handler) Quit() {
	log.Info().Msg("Handler::Quit(): quitting")
	os.Exit(0)
}
```
