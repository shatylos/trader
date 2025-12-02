package tgNotifier

import (
	"bytes"
	"github.com/shatylos/trader/tools/apperrors"
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
			err = apperrors.Wrap(err, "Error sending tg bot notification")
			logger.PrintError(err)
			return
		}
		defer resp.Body.Close()
		logger.Info("Sent notification to tg bot.")
	}()
}
