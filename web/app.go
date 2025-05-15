package web

import (
	"fmt"
	"github.com/shatylos/trader/internal/config"
	"github.com/shatylos/trader/internal/handlers"
	"github.com/shatylos/trader/tools/logger"
	"net/http"
	"os"
)

const defaultPort = "8080"

func StartWebApp() {
	port := defaultPort
	appConfig, err := config.GetConfig()
	if err != nil {
		logger.Error(fmt.Sprintf("Error getting app config %s", err.Error()))
	} else if appConfig.App["web_port"] != "" {
		port = appConfig.App["web_port"]
	}

	logger.Info(fmt.Sprintf("Starting web app http://127.0.0.1:%s/", port))
	mux := http.NewServeMux()
	mux.Handle("GET /", http.HandlerFunc(handlers.SetupListController))
	mux.Handle("GET /report/{setup_id}/{period}/", http.HandlerFunc(handlers.ReportPeriodHandler))
	mux.Handle("GET /report/{setup_id}/", http.HandlerFunc(handlers.ReportHandler))
	mux.Handle("GET /assets/{setup_id}/", http.HandlerFunc(handlers.AssetFormHandler))
	mux.Handle("POST /assets/{setup_id}/", http.HandlerFunc(handlers.AssetAddHandler))
	mux.Handle("GET /stats/", http.HandlerFunc(handlers.StatsHandler))

	err = http.ListenAndServe(fmt.Sprintf(":%s", port), mux)
	if err != nil {
		logger.Error(fmt.Sprintf("Error http ListenAndServe: %s", err.Error()))
		os.Exit(1)
	}
}
