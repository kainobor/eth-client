package main

import (
	"fmt"
	"os"

	"github.com/kainobor/eth-client/app/args"
	"github.com/kainobor/eth-client/app/blockchain"
	"github.com/kainobor/eth-client/app/config"
	"github.com/kainobor/eth-client/app/controller"
	"github.com/kainobor/eth-client/app/handler"
	"github.com/kainobor/eth-client/app/logger"
	"github.com/kainobor/eth-client/app/server"
	"github.com/kainobor/eth-client/app/storage"
	_ "github.com/lib/pq"
)

const (
	argsErrorCode   = 6
	configErrorCode = 7
)

func main() {
	a := args.New()
	a.Init()
	if err := a.Validate(); err != nil {
		fmt.Printf("Error with arguments validating: %v", err)
		os.Exit(argsErrorCode)
	}

	c, err := config.Init(a)
	if err != nil {
		fmt.Printf("Error with config initiating: %v", err)
		os.Exit(configErrorCode)
	}

	log := logger.New()
	log.Init(a.Env, c.Logger)

	bc := blockchain.New(c.Blockchain)
	bc.Init()
	defer bc.Close()

	st := storage.New(c.Storage)
	if err := st.Connect(); err != nil {
		log.Fatalw("error while connecting storage", "config", c.Storage, "error", err)
	}
	defer st.Close()

	h := handler.New(c.Handler, bc, st, log)
	if err := h.Handle(c.Confirmation); err != nil {
		log.Fatalw("error while starting handling", "error", err)
	}

	ctrl := controller.New(bc, st, h, c.Confirmation, log)

	srv := server.New(c.Server, bc, log)
	srv.RegisterRoutes(ctrl)

	log.Info("Starting server")
	if err := srv.Start(); err != nil {
		log.Fatalw("server working failed", "error", err)
	}

}
