package storage

const (
	// UpsertBalanceSQL inserts new balance or updates if have balance with the same address
	UpsertBalanceSQL = `INSERT INTO eth_client.eth_balance (address, balance) VALUES ($1, $2) ON CONFLICT (address) DO UPDATE SET balance = $2;`
	// InsertEntryTransactionSQL inserts new entry transaction
	InsertEntryTransactionSQL = `INSERT INTO eth_client.transactions_entry (hash, block_hash, block_number, from_addr, to_addr, created_at, amount, confirmations) VALUES ($1, $2, $3, $4, $5, $6, $7, 0) RETURNING id;`
	// InsertWithdrawTransactionSQL inserts new withdraw transaction
	InsertWithdrawTransactionSQL = `INSERT INTO eth_client.transactions_withdraw (hash, from_addr, to_addr, amount, created_at) VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP);`
	// SelectTransactionsByStatusSQL selects all transactions with some status
	SelectTransactionsByStatusSQL = `SELECT id, hash, block_hash, block_number, from_addr, to_addr, confirmations, amount, status, created_at FROM eth_client.transactions_entry WHERE status = $1`
	// SelectLastTransactionsSQL selects all transactions that are not showed and with confirmations less that some value
	SelectLastTransactionsSQL = `SELECT id, hash, block_hash, block_number, from_addr, to_addr, confirmations, amount, status, created_at FROM eth_client.transactions_entry WHERE showed = FALSE OR confirmations < $1`
	// UpdateConfirmationsSQL update confirmation value for some entry transaction
	UpdateConfirmationsSQL = `UPDATE eth_client.transactions_entry SET confirmations = $1 WHERE id = $2`
	// UpdateTransactionStatusSQL update status for some entry transaction
	UpdateTransactionStatusSQL = `UPDATE eth_client.transactions_entry SET status = $1 WHERE id = $2`
	// UpdateTransactionsShowedSQL is batch updating of showed value
	UpdateTransactionsShowedSQL = `UPDATE eth_client.transactions_entry SET showed = true WHERE id IN (%s)`
)
