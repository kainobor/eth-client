package args

import (
	"flag"
	"fmt"
)

type (
	// Args struct that implements arguments of cli-app
	Args struct {
		Env          string
		ConfigPath   string
		ConfigPrefix string
	}
)

const (
	// EnvDev argument value development environment
	EnvDev = "dev"
	// EnvProd argument value for production environment
	EnvProd = "prod"

	configPathFlag   = "cp"
	configPrefixFlag = "cn"
	envFlag          = "e"

	defaultConfigPath   = "./config"
	defaultConfigPrefix = "config"
	defaultEnv          = EnvDev
)

// New arguments
func New() *Args {
	return new(Args)
}

// Init flags and parsing them
func (a *Args) Init() {
	a.ConfigPath = *flag.String(configPathFlag, defaultConfigPath, "Path to folder with config")
	a.ConfigPrefix = *flag.String(configPrefixFlag, defaultConfigPrefix, "Prefix of config filename")
	a.Env = *flag.String(envFlag, defaultEnv, fmt.Sprintf("Environment alias: %s or %s", EnvDev, EnvProd))

	flag.Parse()
}

// Validate current values
func (a *Args) Validate() error {
	if a.Env != EnvDev && a.Env != EnvProd {
		return fmt.Errorf("wrong environment")
	}

	return nil
}
