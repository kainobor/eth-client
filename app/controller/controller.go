package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/kainobor/eth-client/app/blockchain"
	"github.com/kainobor/eth-client/app/config"
	"github.com/kainobor/eth-client/app/handler"
	"github.com/kainobor/eth-client/app/helper"
	"github.com/kainobor/eth-client/app/logger"
	"github.com/kainobor/eth-client/app/storage"
)

type (
	// Controller for TCP requests
	Controller struct {
		bc  *blockchain.Client
		st  *storage.Storage
		h   *handler.Handler
		cc  *config.ConfirmationConfig
		log *logger.Logger
	}

	// ErrorResponse returns when something went wrong
	ErrorResponse struct {
		Error string `json:"error"`
	}

	// SuccessResponse returns when all is well
	SuccessResponse struct {
		Message string `json:"message"`
	}

	// LastTransaction is special representation of transaction for GetLast method's JSON
	LastTransaction struct {
		Date          string `json:"date"`
		Address       string `json:"address"`
		Amount        string `json:"amount"`
		Confirmations int64  `json:"confirmations"`
	}
)

const (
	fromSendArg   = "from"
	toSendArg     = "to"
	amountSendArg = "amount"
)

// New controller
func New(
	bc *blockchain.Client,
	st *storage.Storage,
	h *handler.Handler,
	cc *config.ConfirmationConfig,
	log *logger.Logger,
) *Controller {
	return &Controller{bc: bc, st: st, h: h, cc: cc, log: log}
}

// SendEth returns response for SendEth method
func (ctrl *Controller) SendEth(w http.ResponseWriter, r *http.Request) {
	var err error
	params := r.URL.Query()

	from := params.Get(fromSendArg)
	to := params.Get(toSendArg)
	amount := params.Get(amountSendArg)

	if err = ctrl.validateSendRequest(from, to, amount); err != nil {
		errMsg := "invalid request: " + err.Error()
		ctrl.sendError(w, errMsg, "request", r.Body, "error", err)
		return
	}

	var t *blockchain.Transaction

	t, err = blockchain.NewTransaction(from, to, amount)
	if err != nil {
		errMsg := "error while creating transaction"
		ctrl.sendError(w, errMsg, "request", r.Body, "error", err)
		return
	}

	go ctrl.processTransaction(t)

	ctrl.sendResponse(w, "transaction sent for processing", true)
}

// GetLast returns response for GetLast method
func (ctrl *Controller) GetLast(w http.ResponseWriter, r *http.Request) {
	txs, err := ctrl.st.LoadLastTransactions(ctrl.cc.ForLastConfirmationsAmount)
	if err != nil {
		errMsg := "error while loading last transaction"
		ctrl.sendError(w, errMsg, "error", err)
		return
	}

	response := make([]*LastTransaction, 0)
	for _, t := range txs {
		lastTransaction := &LastTransaction{
			Date:          t.CreatedAt().Format(time.RFC850),
			Address:       t.To(),
			Amount:        helper.BigToHex(t.Value()),
			Confirmations: t.Confirmations(),
		}

		response = append(response, lastTransaction)
	}

	respJSON, err := json.Marshal(response)
	if err != nil {
		errMsg := "error while marshaling response"
		ctrl.sendError(w, errMsg, "resp", response, "error", err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(200)
	w.Write(respJSON)

	if err := ctrl.st.TransactionsShowed(txs); err != nil {
		ctrl.log.Errorw("error while setting transactions as showed", "error", err)
	}
}

func (ctrl *Controller) sendError(w http.ResponseWriter, errMsg string, keysAndValues ...interface{}) {
	ctrl.log.Errorw(errMsg, keysAndValues...)
	ctrl.sendResponse(w, errMsg, false)
}

func (ctrl *Controller) sendResponse(w http.ResponseWriter, msg string, isSuccess bool) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(200) // success

	var resp interface{}
	if isSuccess {
		resp = &SuccessResponse{Message: msg}
	} else {
		resp = &ErrorResponse{Error: msg}
	}

	respJSON, err := json.Marshal(resp)
	if err != nil {
		ctrl.log.Errorw("error while response marshaling", "resp", resp, "err", err)
	}
	w.Write(respJSON)
}

func (ctrl *Controller) processTransaction(t *blockchain.Transaction) {
	var txHash string
	txHash, err := ctrl.bc.SendTransaction(t)
	if err != nil {
		ctrl.log.Errorw("error while sending transaction", "transaction", t, "error", err)
		return
	}
	t.SetHash(txHash)

	// For getting information about block
	if err = ctrl.bc.RenewTransaction(t); err != nil {
		ctrl.log.Errorw("error while renewing transaction", "transaction", t, "error", err)
		return
	}
	t.FixateCreatedAt()

	go ctrl.saveTransaction(t)
}

func (ctrl *Controller) saveTransaction(t *blockchain.Transaction) {
	if err := ctrl.st.SaveEntryTransaction(t); err != nil {
		ctrl.log.Errorw("error while saving entry transaction", "transaction", t, "error", err)
	}

	if err := ctrl.st.SaveWithdrawTransaction(t); err != nil {
		ctrl.log.Errorw("error while saving withdraw transaction", "transaction", t, "error", err)
	}

	ctrl.h.AddTransaction(t)
}

func (ctrl *Controller) validateSendRequest(from, to, amount string) error {
	switch {
	case !helper.IsHexAddress(from):
		return fmt.Errorf("wrong sender address in request")
	case !helper.IsHexAddress(to):
		return fmt.Errorf("wrong sender address in request")
	case !helper.IsHexString(amount):
		return fmt.Errorf("invalid amount format")
	}

	return nil
}
