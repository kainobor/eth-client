package blockchain

import (
	"encoding/json"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/kainobor/eth-client/app/helper"
)

type (
	// Transaction represents entry transaction
	Transaction struct {
		id            int64
		hash          string
		from          string
		to            string
		value         big.Int
		confirmations int64
		block         Block
		status        string
		createdAt     time.Time
		sync.RWMutex
	}

	// Block represents blockchain block
	Block struct {
		number big.Int
		hash   string
	}

	// DBTransaction represents transaction data, that saved in DB
	DBTransaction struct {
		ID            int64
		Hash          string
		BlockHash     string
		BlockNumber   int64
		From          string
		To            string
		Confirmations int64
		Amount        string
		Status        string
		CreatedAt     time.Time
	}
)

const (
	// PendingStatus is status for pending transaction
	PendingStatus = "pending"
	// SuccessStatus is status for confirmed transaction
	SuccessStatus = "success"
	// FailStatus is status for transaction that is not in network already
	FailStatus = "fail"
)

// NewTransaction is constructor for transactions
func NewTransaction(from, to, value string) (*Transaction, error) {
	bigValue, ok := helper.HexToBig(value)
	if !ok {
		return nil, fmt.Errorf("can't parsing `%s` to big.Int", value)
	}

	return &Transaction{from: from, to: to, value: *bigValue}, nil
}

// MarshalJSON implements the json.Unmarshaler interface
func (t *Transaction) MarshalJSON() ([]byte, error) {
	t.RLock()
	params := map[string]interface{}{
		"from":  t.from,
		"to":    t.to,
		"value": helper.BigToHex(t.value),
	}
	t.RUnlock()

	return json.Marshal(params)
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (t *Transaction) UnmarshalJSON(data []byte) error {
	params := struct {
		BlockHash   string `json:"blockHash"`
		BlockNumber string `json:"blockNumber"`
		From        string `json:"from"`
		To          string `json:"to"`
		Hash        string `json:"hash"`
		Value       string `json:"value"`
	}{}

	if err := json.Unmarshal(data, &params); err != nil {
		return fmt.Errorf("error while unmarshaling transaction: %v", err)
	}

	var ok bool
	var blockNumber *big.Int

	t.Lock()
	t.hash = params.Hash

	val, ok := helper.HexToBig(params.Value)
	if !ok {
		t.Unlock()
		return fmt.Errorf("wrong transaction value: %s", params.Value)
	}
	t.value = *val

	block := Block{}
	block.hash = params.BlockHash
	if blockNumber, ok = helper.HexToBig(params.BlockNumber); !ok {
		t.Unlock()
		return fmt.Errorf("wrong block number: %s", params.BlockNumber)
	}
	block.number = *blockNumber

	t.block = block
	t.to = params.To
	t.from = params.From

	t.Unlock()
	return nil
}

// FillFromDB gets values from DB and updates transaction with them
func (t *Transaction) FillFromDB(dbt *DBTransaction) error {
	t.Lock()
	defer t.Unlock()

	val, ok := helper.HexToBig(dbt.Amount)
	if !ok {
		return fmt.Errorf("can't parse amount `%s` from DB", dbt.Amount)
	}
	t.value = *val

	block := new(Block)
	block.hash = dbt.Hash
	block.number = *big.NewInt(dbt.BlockNumber)
	t.block = *block

	t.id = dbt.ID
	t.hash = dbt.Hash
	t.from = dbt.From
	t.to = dbt.To
	t.confirmations = dbt.Confirmations
	t.status = dbt.Status
	t.createdAt = dbt.CreatedAt

	return nil
}

// FixateCreatedAt saved current time as "createdAt"
func (t *Transaction) FixateCreatedAt() {
	t.Lock()
	if t.createdAt.Equal(time.Time{}) {
		t.createdAt = time.Now()
	}
	t.Unlock()
}

// BlockHash is synchronous getter
func (t *Transaction) BlockHash() string {
	t.RLock()
	defer t.RUnlock()

	return t.block.hash
}

// BlockNumber is synchronous getter
func (t *Transaction) BlockNumber() big.Int {
	t.RLock()
	defer t.RUnlock()

	return t.block.number
}

// Confirmations is synchronous getter
func (t *Transaction) Confirmations() int64 {
	t.RLock()
	defer t.RUnlock()

	return t.confirmations
}

// SetConfirmations is synchronous setter
func (t *Transaction) SetConfirmations(cnf int64) {
	t.Lock()
	t.confirmations = cnf
	t.Unlock()
}

// Status is synchronous getter
func (t *Transaction) Status() string {
	t.RLock()
	defer t.RUnlock()

	return t.status
}

// SetStatus is synchronous setter
func (t *Transaction) SetStatus(st string) {
	t.Lock()
	t.status = st
	t.Unlock()
}

// ID is synchronous getter
func (t *Transaction) ID() int64 {
	t.RLock()
	defer t.RUnlock()

	return t.id
}

// SetID is synchronous getter
func (t *Transaction) SetID(id int64) {
	t.Lock()
	t.id = id
	t.Unlock()
}

// From is synchronous getter
func (t *Transaction) From() string {
	t.RLock()
	defer t.RUnlock()

	return t.from
}

// To is synchronous getter
func (t *Transaction) To() string {
	t.RLock()
	defer t.RUnlock()

	return t.to
}

// Value is synchronous getter
func (t *Transaction) Value() big.Int {
	t.RLock()
	defer t.RUnlock()

	return t.value
}

// Hash is synchronous getter
func (t *Transaction) Hash() string {
	t.RLock()
	defer t.RUnlock()

	return t.hash
}

// SetHash is synchronous setter
func (t *Transaction) SetHash(hash string) {
	t.Lock()
	t.hash = hash
	t.Unlock()
}

// CreatedAt is synchronous getter
func (t *Transaction) CreatedAt() time.Time {
	t.RLock()
	defer t.RUnlock()

	return t.createdAt
}
