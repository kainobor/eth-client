package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kainobor/eth-client/app/blockchain"
	"github.com/kainobor/eth-client/app/config"
	"github.com/kainobor/eth-client/app/controller"
	"github.com/kainobor/eth-client/app/logger"
)

const (
	sendEthRoute = "/SendEth"
	getLastRoute = "/GetLast"
)

type (
	// Server implementation of IServer
	Server struct {
		config *config.ServerConfig
		router *mux.Router
		bc     *blockchain.Client
		log    *logger.Logger
	}
)

// New server
func New(c *config.ServerConfig, bc *blockchain.Client, l *logger.Logger) *Server {
	return &Server{config: c, router: mux.NewRouter(), bc: bc, log: l}
}

// RegisterRoutes registers all available routes in server router
func (srv *Server) RegisterRoutes(ctrl *controller.Controller) {
	srv.router.HandleFunc(sendEthRoute, ctrl.SendEth).Methods("GET")
	srv.router.HandleFunc(getLastRoute, ctrl.GetLast).Methods("GET")
}

// Start listening of TCP-connections
func (srv *Server) Start() error {
	err := http.ListenAndServe(fmt.Sprintf(":%d", srv.config.Port), srv.router)

	return fmt.Errorf("starting was ended with error: %v", err)
}
