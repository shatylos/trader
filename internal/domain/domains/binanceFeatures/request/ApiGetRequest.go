package request

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	binanceStructs "github.com/shatylos/trader/internal/domain/domains/binanceFeatures/structs"
	"github.com/shatylos/trader/tools"
	"github.com/shatylos/trader/tools/logger"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type ApiGetRequest struct {
	Uri       string
	ApiParams binanceStructs.ApiParams
	Secrets   binanceStructs.Secrets
}

var httpClient *http.Client

func (r *ApiGetRequest) DoRequest() (response binanceStructs.ApiResponse, err error) {
	if r.Secrets.Verbose {
		var paramsJsonBytes []byte
		paramsJsonBytes, err = json.Marshal(r.ApiParams)
		if err != nil {
			msg := fmt.Sprintf("Error marshalling apiParams to JSON: %s", err)
			logger.Error(msg)
			err = tools.AppError{
				Message:     msg,
				ParentError: err,
			}
			return
		}
		logger.Info(fmt.Sprintf("Send GET request to: URI: %s. Params: %s", r.Uri, paramsJsonBytes))
	}

	var requestUrl string
	requestUrl = r.getRequestData()
	requestBody := bytes.NewBuffer(nil)

	var request *http.Request
	request, err = http.NewRequest("GET", requestUrl, requestBody)
	if err != nil {
		msg := fmt.Sprintf("Error creating newRequest: %s", err)
		logger.Error(msg)
		err = tools.AppError{
			Message:     msg,
			ParentError: err,
		}
		return
	}

	request.Header.Add("Content-Type", "application/json")

	httpClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{},
		},
		Timeout: 10 * time.Second,
	}

	resp, err := httpClient.Do(request)
	if err != nil {
		msg := fmt.Sprintf("Binance Features API error do request: %s", err)
		logger.Error(msg)
		err = tools.AppError{
			Message:     msg,
			ParentError: err,
		}
		return
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			logger.Error(fmt.Sprintf("Binance Features API error closing response body: %s", err))
		}
	}(resp.Body)

	if resp.Status != "200 OK" {
		msg := fmt.Sprintf("Binance Features API error. http status: %s", resp.Status)
		logger.Error(msg)
		err = tools.AppError{
			Message: msg,
		}
		return
	}

	body := bytes.Buffer{}
	_, err = io.Copy(&body, resp.Body)
	if err != nil {
		msg := fmt.Sprintf("Binance Features API error read body: %s", err)
		logger.Error(msg)
		err = tools.AppError{
			Message:     msg,
			ParentError: err,
		}
		return
	}
	response = body.Bytes()
	body.Reset()
	return
}

func (r *ApiGetRequest) getRequestData() (requestUrl string) {
	queryParams := url.Values{}
	for key, value := range r.ApiParams {
		queryParams.Add(key, value)
	}

	var builder strings.Builder
	builder.WriteString(r.Secrets.ApiEndpoint)
	builder.WriteString(r.Uri)
	builder.WriteByte('?')
	builder.WriteString(queryParams.Encode())
	requestUrl = builder.String()
	return
}
