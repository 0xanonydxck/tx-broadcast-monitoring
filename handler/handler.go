package handler

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"tx-monitoring/config"
	"tx-monitoring/domain"
	"tx-monitoring/model"

	"github.com/cqroot/prompt"
	"github.com/cqroot/prompt/input"
	"github.com/go-co-op/gocron"
	"github.com/go-playground/validator/v10"
	"github.com/jedib0t/go-pretty/table"
	"github.com/rs/zerolog/log"
)

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

func (h *handler) Broadcast() {
	symbols, err := h.service.Symbols()
	if err != nil {
		return
	}

	req := new(model.TxBroadcastReq)
	choices := make([]string, len(symbols))
	for i, symbol := range symbols {
		choices[i] = fmt.Sprintf("%v - %v", symbol.Symbol, symbol.Price)
	}

	symbolWithPrice, err := prompt.New().Ask("Choose Symbol").Choose(choices)
	if err != nil {
		log.Panic().Err(err).Msg("Handler::Broadcast(): failed to retrieve answer")
	}

	parts := strings.Split(symbolWithPrice, " - ")
	req.Symbol = strings.TrimSpace(parts[0])
	symbolPrice, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err != nil {
		log.Panic().Err(err).Msg("Handler::Broadcast(): failed to parse price")
	}

	priceStr, err := prompt.New().Ask("Enter Price (uint)").
		Input("1",
			input.WithInputMode(input.InputInteger),
			input.WithValidateFunc(func(s string) error {
				if err := h.validator.Var(s, "required,gt=0"); err != nil {
					return err
				}

				inputPrice, err := strconv.ParseFloat(s, 64)
				if err != nil {
					return err
				}

				if inputPrice < symbolPrice {
					return fmt.Errorf("price must be greater than or equal to %v", symbolPrice)
				}

				return nil
			}))

	if err != nil {
		log.Panic().Err(err).Msg("Handler::Broadcast(): failed to retrieve answer")
	}

	price, err := strconv.ParseUint(priceStr, 10, 64)
	if err != nil {
		log.Panic().Err(err).Msg("Handler::Broadcast(): failed to parse price")
	}
	req.Price = price
	req.Timestamp = time.Now().Unix()

	_, err = h.service.Broadcast(req)
	if err != nil {
		return
	}
}

func (h *handler) ShowTransactions() {
	txs, err := h.service.Transactions()
	if err != nil {
		return
	}

	if len(txs) == 0 {
		log.Info().Msg("Handler::ShowTransactions(): no transactions found")
		return
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"TxHash", "Status", "Created At", "Updated At"})
	for _, tx := range txs {
		t.AppendRow(table.Row{
			tx.TxHash,
			tx.Status,
			tx.CreatedAt.Format("2006-01-02 15:04:05"),
			tx.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	t.Render()
}

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
