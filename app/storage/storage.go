package storage

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/kainobor/eth-client/app/blockchain"
	"github.com/kainobor/eth-client/app/config"
	"github.com/kainobor/eth-client/app/helper"
)

type (
	// Storage is client for database connection
	Storage struct {
		config *config.StorageConfig
		db     *sql.DB
	}
)

// New storage
func New(config *config.StorageConfig) *Storage {
	return &Storage{config: config}
}

// Connect to DB
func (st *Storage) Connect() error {
	var err error
	if st.db, err = sql.Open("postgres", connectString(st.config)); err != nil {
		return fmt.Errorf("storage connectiong error: %v", err)
	}

	if err = st.db.Ping(); err != nil {
		return fmt.Errorf("storage is not responding: %v", err)
	}

	return nil
}

// UpsertBalance inserts or updates balance by some address
func (st *Storage) UpsertBalance(addr, balance string) error {
	if _, err := st.db.Exec(UpsertBalanceSQL, strings.ToLower(addr), balance); err != nil {
		return err
	}

	return nil
}

// SaveEntryTransaction inserts entry transaction to DB and renews transaction ID
func (st *Storage) SaveEntryTransaction(t *blockchain.Transaction) error {
	blockNum := t.BlockNumber()
	var insertedID int64
	err := st.db.QueryRow(
		InsertEntryTransactionSQL,
		t.Hash(),
		t.BlockHash(),
		(&blockNum).Int64(),
		t.From(),
		t.To(),
		t.CreatedAt(),
		helper.BigToHex(t.Value()),
	).Scan(&insertedID)
	if err != nil {
		return fmt.Errorf("transaction `%s` not inserted: %v", t.Hash(), err)
	}

	t.SetID(insertedID)

	return nil
}

// SaveWithdrawTransaction saves transaction as withdraw
func (st *Storage) SaveWithdrawTransaction(t *blockchain.Transaction) error {
	if _, err := st.db.Exec(InsertWithdrawTransactionSQL, t.Hash(), t.From(), t.To(), helper.BigToHex(t.Value())); err != nil {
		return err
	}

	return nil
}

// UpdateConfirmations updates confirmations amount by ID and checks that row was updated
func (st *Storage) UpdateConfirmations(id, confirmations int64) error {
	res, err := st.db.Exec(UpdateConfirmationsSQL, confirmations, id)
	if err != nil {
		return fmt.Errorf("error while executing confirmations updating: %v", err)
	}

	if affected, err := res.RowsAffected(); err != nil {
		return fmt.Errorf("error while getting affected rows: %v", err)
	} else if affected == 0 {
		return fmt.Errorf("transaction #%d not updated", id)
	}

	return nil
}

// UpdateTransactionStatus updates status by ID and checks that row was updated
func (st *Storage) UpdateTransactionStatus(id int64, status string) error {
	res, err := st.db.Exec(UpdateTransactionStatusSQL, status, id)
	if err != nil {
		return fmt.Errorf("error while executing status updating: %v", err)
	}

	if affected, err := res.RowsAffected(); err != nil {
		return fmt.Errorf("error while getting affected rows: %v", err)
	} else if affected == 0 {
		return fmt.Errorf("transaction #%d not updated", id)
	}

	return nil
}

// TransactionsShowed updates all transactions as showed page by page
func (st *Storage) TransactionsShowed(txs map[string]*blockchain.Transaction) error {
	var ids []string
	var i int
	var rest = len(txs)
	for _, tx := range txs {
		i++
		ids = append(ids, fmt.Sprintf("%d", tx.ID()))

		if i < st.config.PageSize && i < rest {
			continue
		}

		i = 0
		rest = rest - st.config.PageSize
		inCondition := strings.Join(ids, ",")
		_, err := st.db.Exec(fmt.Sprintf(UpdateTransactionsShowedSQL, inCondition))
		if err != nil {
			return fmt.Errorf("error while setting transactions as showed: %v", err)
		}

		// clear slice
		ids = ids[:0]
	}

	return nil
}

// LoadTransactionsByStatus returns all transactions with some status
func (st *Storage) LoadTransactionsByStatus(status string) (map[string]*blockchain.Transaction, error) {
	return st.loadTransactions(SelectTransactionsByStatusSQL, status)
}

// LoadLastTransactions returns all not shown transactions
// with amount of confirmations less than some value
func (st *Storage) LoadLastTransactions(lastConfirmations int64) (map[string]*blockchain.Transaction, error) {
	return st.loadTransactions(SelectLastTransactionsSQL, lastConfirmations)
}

// Close DB connection
func (st *Storage) Close() error {
	if err := st.db.Close(); err != nil {
		return fmt.Errorf("storage closing error: %v", err)
	}

	return nil
}

func (st *Storage) loadTransactions(query string, args ...interface{}) (map[string]*blockchain.Transaction, error) {
	txs := make(map[string]*blockchain.Transaction)

	rows, err := st.db.Query(query, args...)
	if err != nil {
		return txs, fmt.Errorf("error while selecting transactions: %v", err)
	}

	for rows.Next() {
		var dbTx = new(blockchain.DBTransaction)
		err = rows.Scan(&dbTx.ID, &dbTx.Hash, &dbTx.BlockHash, &dbTx.BlockNumber, &dbTx.From, &dbTx.To, &dbTx.Confirmations, &dbTx.Amount, &dbTx.Status, &dbTx.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("error while scanning transaction: %v", err)
		}

		tx := new(blockchain.Transaction)
		if err := tx.FillFromDB(dbTx); err != nil {
			return nil, fmt.Errorf("error while filling transaction: %v", tx)
		}

		txs[tx.Hash()] = tx
	}

	return txs, nil
}

func connectString(config *config.StorageConfig) string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.IP, config.Port, config.User, config.Password, config.DBName,
	)
}
