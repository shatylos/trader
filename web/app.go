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

	http.Handle("GET /", http.HandlerFunc(handlers.SetupListController))
	http.Handle("GET /report/{setup_id}/{period}/", http.HandlerFunc(handlers.ReportController))

	err = http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		logger.Error(fmt.Sprintf("Error http ListenAndServe: %s", err.Error()))
		os.Exit(1)
	}
}
