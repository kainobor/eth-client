package blockchain

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/kainobor/eth-client/app/config"
	"github.com/kainobor/eth-client/app/helper"
)

type (
	// Client represents client of ethereum network
	Client struct {
		config *config.BlockchainConfig
		rpc    *rpc.Client
	}
)

const (
	sendTransactionMethod      = "eth_sendTransaction"
	getTransactionByHashMethod = "eth_getTransactionByHash"
	getBlockByNumberMethod     = "eth_getBlockByNumber"
	getBalanceMethod           = "eth_getBalance"
	getCurrentBlockMethod      = "eth_blockNumber"
)

// New client of ethereum network
func New(c *config.BlockchainConfig) *Client {
	return &Client{config: c}
}

// Init connections to network
func (cl *Client) Init() error {
	fullAddr := fmt.Sprintf("http://%s:%d", cl.config.IP, cl.config.Port)
	var err error

	if cl.rpc, err = rpc.Dial(fullAddr); err != nil {
		return fmt.Errorf("error while connecting to RPC: %v", err)
	}

	return nil
}

// SendTransaction sends unsigned transaction to network
func (cl *Client) SendTransaction(t *Transaction) (string, error) {
	var txID string
	if err := cl.rpc.Call(&txID, sendTransactionMethod, t); err != nil {
		return "", err
	}

	return txID, nil
}

// GetBalance returns balance by some address
func (cl *Client) GetBalance(addr string) (*big.Int, error) {
	var balanceHex string
	if err := cl.rpc.Call(&balanceHex, getBalanceMethod, addr, "latest"); err != nil {
		return nil, fmt.Errorf("error while getting balance: %v", err)
	}

	balance, ok := helper.HexToBig(balanceHex)
	if !ok {
		return nil, fmt.Errorf("can't parse `%s` as balance", balanceHex)
	}

	return balance, nil
}

// RenewTransaction renews transaction values from network
func (cl *Client) RenewTransaction(t *Transaction) error {
	err := cl.rpc.Call(&t, getTransactionByHashMethod, t.Hash())
	if err != nil {
		return fmt.Errorf("can't get transaction: %v", err)
	}

	return nil
}

// GetCurrentBlock returns most recent block from network
func (cl *Client) GetCurrentBlock() (*big.Int, error) {
	var blockNumHex string
	if err := cl.rpc.Call(&blockNumHex, getCurrentBlockMethod); err != nil {
		return nil, fmt.Errorf("error while getting current block: %v", err)
	}

	blockNum, ok := helper.HexToBig(blockNumHex)
	if !ok {
		return nil, fmt.Errorf("can't parse `%s` as block number", blockNumHex)
	}

	return blockNum, nil
}

// BlockExists checks that block with certain hash and number exists in network
func (cl *Client) BlockExists(blockNumber big.Int, blockHash string) (bool, error) {
	var blockData = make(map[string]interface{})

	err := cl.rpc.Call(&blockData, getBlockByNumberMethod, helper.BigToHex(blockNumber), false)
	if err != nil {
		return false, fmt.Errorf("can't get block by number: %v", err)
	}

	blockHashRaw, ok := blockData["hash"]
	if !ok {
		return false, fmt.Errorf("can't get hash from block data")
	} else if blockHashRaw.(string) == "" {
		return false, fmt.Errorf("empty block hash")
	}

	if blockHashRaw == blockHashRaw.(string) {
		return true, nil
	}

	return false, nil
}

// Close connection
func (cl *Client) Close() {
	cl.rpc.Close()
}
