package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/kainobor/eth-client/app/handler"
	"github.com/kainobor/eth-client/app/helper"
	"github.com/kainobor/eth-client/app/logger"
	"github.com/kainobor/eth-client/app/storage"

	"github.com/kainobor/eth-client/app/blockchain"
	"github.com/kainobor/eth-client/app/config"
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

	if err = ctrl.addBalanceIfNeed(from); err != nil {
		errMsg := "saving sender balance error"
		ctrl.sendError(w, errMsg, "addr", from, "error", err)
		return
	}
	if err = ctrl.addBalanceIfNeed(to); err != nil {
		errMsg := "saving recipient balance error"
		ctrl.sendError(w, errMsg, "addr", to, "error", err)
		return
	}

	var t *blockchain.Transaction

	t, err = blockchain.NewTransaction(from, to, amount)
	if err != nil {
		errMsg := "error while creating transaction"
		ctrl.sendError(w, errMsg, "request", r.Body, "error", err)
		return
	}

	var txHash string
	txHash, err = ctrl.bc.SendTransaction(t)
	if err != nil {
		errMsg := "error while sending transaction"
		ctrl.sendError(w, errMsg, "error", err)
		return
	}
	t.SetHash(txHash)

	if err = ctrl.bc.RenewTransaction(t); err != nil {
		errMsg := "error while renewing transaction"
		ctrl.sendError(w, errMsg, "error", err)
		return
	}
	t.FixateCreatedAt()

	if err != nil {
		errMsg := "error while sending transaction"
		ctrl.sendError(w, errMsg, "transaction", t, "error", err)
		return
	}

	if err = ctrl.st.SaveEntryTransaction(t); err != nil {
		errMsg := "error while saving entry transaction"
		ctrl.sendError(w, errMsg, "transaction", t, "error", err)
		return
	}

	if err = ctrl.st.SaveWithdrawTransaction(t); err != nil {
		errMsg := "error while saving withdraw transaction"
		ctrl.sendError(w, errMsg, "transaction", t, "error", err)
		return
	}

	ctrl.h.AddTransaction(t)

	ctrl.sendSuccess([]byte(txHash), w)

	return
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

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(400)

	resp := &ErrorResponse{Error: errMsg}
	respJSON, err := json.Marshal(resp)
	if err != nil {
		ctrl.log.Errorw("error while error response marshaling", "resp", resp, "err", err)
	}
	w.Write(respJSON)
}

func (ctrl *Controller) sendSuccess(data []byte, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(200) // success

	w.Write(data)
}

func (ctrl *Controller) addBalanceIfNeed(addr string) error {
	bal, err := ctrl.bc.GetBalance(addr)
	if err != nil {
		return fmt.Errorf("can't get `%s` balance: %v", addr, err)
	}

	if err := ctrl.st.UpsertBalance(addr, helper.BigToHex(*bal)); err != nil {
		return fmt.Errorf("error while saving `%s` balance: %v", addr, err)
	}

	return nil
}
