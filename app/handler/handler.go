package handler

import (
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/kainobor/eth-client/app/blockchain"
	"github.com/kainobor/eth-client/app/config"
	"github.com/kainobor/eth-client/app/helper"
	"github.com/kainobor/eth-client/app/logger"
	"github.com/kainobor/eth-client/app/storage"
)

type (
	// Handler of current application data
	Handler struct {
		config         *config.HandlerConfig
		bc             *blockchain.Client
		st             *storage.Storage
		transactions   map[string]*blockchain.Transaction
		curBlockNumber big.Int
		log            *logger.Logger
		sync.RWMutex
	}
)

// New handler
func New(c *config.HandlerConfig, bc *blockchain.Client, st *storage.Storage, log *logger.Logger) *Handler {
	transactions := make(map[string]*blockchain.Transaction)

	return &Handler{config: c, bc: bc, st: st, transactions: transactions, log: log}
}

// Handle app data and gets it from DB before starting
func (h *Handler) Handle(cc *config.ConfirmationConfig) error {
	var err error
	var curBlockNum *big.Int
	if curBlockNum, err = h.bc.GetCurrentBlock(); err != nil {
		return fmt.Errorf("can't get current block: %v", err)
	}

	h.SetCurBlockNum(*curBlockNum)

	if h.transactions, err = h.st.LoadTransactionsByStatus(blockchain.PendingStatus); err != nil {
		return fmt.Errorf("can't load pending transactions: %v", err)
	}

	go func() {
		tcr := time.NewTicker(h.config.TransactionInterval)
		for range tcr.C {
			h.handleTransactions(cc.SuccessConfirmationsAmount)
		}
	}()

	go func() {
		tcr := time.NewTicker(h.config.CurBlockInterval)
		for range tcr.C {
			h.handleCurrentBlock()
		}
	}()

	return nil
}

// AddTransaction adds one transaction to handling queue
func (h *Handler) AddTransaction(t *blockchain.Transaction) {
	h.Lock()
	h.transactions[t.Hash()] = t
	h.Unlock()
}

// SetCurBlockNum is synchronous setter
func (h *Handler) SetCurBlockNum(bn big.Int) {
	h.Lock()
	h.curBlockNumber = bn
	h.Unlock()
}

// CurBlockNum is synchronous getter
func (h *Handler) CurBlockNum() big.Int {
	h.RLock()
	defer h.RUnlock()

	return h.curBlockNumber

}

func (h *Handler) handleTransactions(confirmationsForSuccess int64) {
	var existBlocks = make(map[string]bool)

	// Remember, that copy have the same pointers!
	copyMap := h.copyTransactions()

	for hash, t := range copyMap {
		exist, err := h.checkBlockExistense(existBlocks, t)
		if err != nil {
			h.log.Errorw(err.Error(), "transaction", t)
		}
		if !exist {
			continue
		}

		confirmations := h.currentTransactionConfirmations(t)
		if confirmations != t.Confirmations() {
			if err := h.st.UpdateConfirmations(t.ID(), confirmations); err != nil {
				h.log.Errorw("can't update confirmation", "error", err)
				continue
			}

			t.SetConfirmations(confirmations)
			// Balance may to change if some block before current was cancelled
			h.updateBalances(t)
		}

		if confirmations > confirmationsForSuccess {
			if err := h.st.UpdateTransactionStatus(t.ID(), blockchain.SuccessStatus); err != nil {
				h.log.Errorw("can't update transaction status", "error", err)
				continue
			}

			t.SetStatus(blockchain.SuccessStatus)
			h.delTransaction(hash)
		}
	}
}

func (h *Handler) handleCurrentBlock() {
	num, err := h.bc.GetCurrentBlock()
	if err != nil {
		h.log.Errorw("error while getting current block number", "error", err)
		return
	}

	h.SetCurBlockNum(*num)
}

// updateBalances gets sender and receiver balances from network and saves it to DB
func (h *Handler) updateBalances(t *blockchain.Transaction) error {
	fromBal, err := h.bc.GetBalance(t.From())
	if err != nil {
		return fmt.Errorf("can't get sender balance: %v", err)
	}

	toBal, err := h.bc.GetBalance(t.To())
	if err != nil {
		return fmt.Errorf("can't get receiver balance: %v", err)
	}

	if err := h.st.UpsertBalance(t.From(), helper.BigToHex(*fromBal)); err != nil {
		return fmt.Errorf("can't update sender balance: %v", err)
	}

	if err := h.st.UpsertBalance(t.To(), helper.BigToHex(*toBal)); err != nil {
		return fmt.Errorf("can't update receiver balance: %v", err)
	}

	return nil
}

func (h *Handler) copyTransactions() map[string]*blockchain.Transaction {
	m := make(map[string]*blockchain.Transaction)

	h.RLock()
	for hash, t := range h.transactions {
		m[hash] = t
	}
	h.RUnlock()

	return m
}

func (h *Handler) delTransaction(hash string) {
	h.Lock()
	delete(h.transactions, hash)
	h.Unlock()
}

func (h *Handler) checkBlockExistense(existBlocks map[string]bool, t *blockchain.Transaction) (bool, error) {
	var err error
	blockExist, ok := existBlocks[t.BlockHash()]

	if !ok {
		if blockExist, err = h.bc.BlockExists(t.BlockNumber(), t.BlockHash()); err != nil {
			return false, fmt.Errorf("can't check is block exists: %v", err)
		}
		existBlocks[t.BlockHash()] = blockExist
	}

	if !blockExist {
		if err := h.st.UpdateTransactionStatus(t.ID(), blockchain.FailStatus); err != nil {
			return false, fmt.Errorf("can't set transaction failure: %v", err)
		}
		t.SetStatus(blockchain.FailStatus)
		h.delTransaction(t.Hash())

		// Balance may to change if block is cancelled
		h.updateBalances(t)
		return false, nil
	}

	return blockExist, nil
}

func (h *Handler) currentTransactionConfirmations(t *blockchain.Transaction) int64 {
	curConfirmationsBig := big.NewInt(0)
	transBlock := t.BlockNumber()
	curBlock := h.CurBlockNum()

	return curConfirmationsBig.Sub(&curBlock, &transBlock).Int64()
}
