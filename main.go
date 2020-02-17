package main

import (
	"email-acceptor/endpoints/rest"
	"email-acceptor/helpers"
	"email-acceptor/pkg/graceful"
	"email-acceptor/pkg/logger"
	"email-acceptor/pkg/middlewares"
	"email-acceptor/store"
	"github.com/Zetkolink/email-sender/models"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"time"
)

func main() {
	cfg := helpers.InitConfig()

	db := helpers.InitDb(cfg)

	rb := helpers.InitRabbit(cfg)

	ch := rb.InitChannel(cfg)

	lg := logger.New(os.Stderr, cfg.LogLevel, cfg.LogFormat)

	ms := store.NewMessageStore(db, ch, cfg.Rb.Exchange)
	db.AutoMigrate(&models.Message{})
	ms.Migrate()

	restHandler := rest.New(lg, ms)

	srv := setupServer(cfg, lg, restHandler)
	lg.Infof("listening for requests on %s...", cfg.Addr)
	if err := srv.ListenAndServe(); err != nil {
		lg.Fatalf("http server exited: %s", err)
	}
}

func setupServer(cfg helpers.Config, lg logger.Logger, rest http.Handler) *graceful.Server {
	router := mux.NewRouter()
	router.PathPrefix("/").Handler(rest)

	handler := middlewares.WithRecovery(lg, router)

	srv := graceful.NewServer(handler, 20*time.Second, os.Interrupt)
	srv.Log = lg.Errorf
	srv.Addr = cfg.Addr
	return srv
}
