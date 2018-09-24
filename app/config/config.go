package config

import (
	"fmt"
	"time"

	"github.com/kainobor/eth-client/app/args"

	"github.com/spf13/viper"
)

type (
	// Config is wrapper for configurations af all modules
	Config struct {
		Server       *ServerConfig
		Blockchain   *BlockchainConfig
		Storage      *StorageConfig
		Handler      *HandlerConfig
		Confirmation *ConfirmationConfig
		Logger       *LoggerConfig
	}

	// ServerConfig is config for TCP-server
	ServerConfig struct {
		Port int
	}

	// BlockchainConfig is config for blockchain network client
	BlockchainConfig struct {
		IP   string
		Port int
	}

	// StorageConfig is config for DB-connection
	StorageConfig struct {
		IP       string
		Port     int
		User     string
		Password string
		DBName   string
		PageSize int
	}

	//HandlerConfig is config for handling application data
	HandlerConfig struct {
		TransactionInterval time.Duration
		CurBlockInterval    time.Duration
	}

	// ConfirmationConfig that contains data about acceptance of confirmations
	ConfirmationConfig struct {
		SuccessConfirmationsAmount int64
		ForLastConfirmationsAmount int64
	}

	// LoggerConfig is config for logger
	LoggerConfig struct {
		InfoPaths []string // Where to write standard logs
		ErrPaths  []string // Where to write errors
	}
)

// Init configuration and read from file
func Init(a *args.Args) (*Config, error) {
	c := new(Config)

	vpr := viper.New()
	vpr.AddConfigPath(a.ConfigPath)
	vpr.SetConfigName(a.ConfigPrefix + "_" + a.Env)
	if err := vpr.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error while reading config: %v", err)
	}
	if err := vpr.Unmarshal(c); err != nil {
		return nil, fmt.Errorf("error while unmarshaling config: %v", err)
	}

	return c, nil
}
