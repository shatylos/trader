package main

import (
	"fmt"
	"github.com/shatylos/trader/internal/storage"
	"github.com/shatylos/trader/tools/logger"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	go func() {
		logger.Info("Starting pprof listen: http://127.0.0.1:6061/debug/pprof/")
		err := http.ListenAndServe(":6061", nil)
		logger.Error(fmt.Sprintf("Error listening pprof %s", err.Error()))
	}()

	err := storage.InitStorage()
	if err != nil {
		logger.Error(fmt.Sprintf("Error initialization storage: %s", err.Error()))
	}
}
