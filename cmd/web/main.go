package main

import (
	"context"
	"net/http"
	"os"

	"github.com/hernanhrm/budget-forge/cmd/web/routes"
	"github.com/hernanhrm/budget-forge/pkg/database"
	"github.com/hernanhrm/budget-forge/pkg/localconfig"
	"github.com/hernanhrm/budget-forge/pkg/logger"
	"github.com/hernanhrm/budget-forge/pkg/server"
)

func main() {
	log := logger.NewDevelopment()

	cfg, err := localconfig.GetConfig(log)
	if err != nil {
		log.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	db, err := database.NewConnection(context.Background(), cfg.Database.URL, log)
	if err != nil {
		log.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	srv := server.NewServer(server.Config{
		Port:  cfg.Service.Port(),
		Debug: cfg.Service.Debug,
	}, log, func(mux *http.ServeMux) {
		routes.Setup(mux, db)
	})

	if err := srv.Start(context.Background()); err != nil {
		log.Error("server error", "error", err)
		os.Exit(1)
	}
}
