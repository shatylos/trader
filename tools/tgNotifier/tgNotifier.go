package tgNotifier

import (
	"bytes"
	"fmt"
	"github.com/shatylos/trader/tools/logger"
	"net/http"
)

type Notifier struct {
	IsInit bool
	Url    string
}

var notifier Notifier

func Init(url string) {
	notifier.IsInit = true
	notifier.Url = url
}

func Notify(message string) {
	go func() {
		if !notifier.IsInit {
			logger.Error("Notifier was not init.")
			return
		}
		if notifier.Url == "" {
			logger.Error("Notifier url was not set.")
			return
		}

		data := []byte(message)

		resp, err := http.Post(notifier.Url, "text/plain", bytes.NewBuffer(data))
		if err != nil {
			logger.Error(fmt.Sprintf("Error sending tg bot notification: %s", err.Error()))
			return
		}
		defer resp.Body.Close()
		logger.Info("Sent notification to tg bot.")
	}()
}
