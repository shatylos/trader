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
	"strconv"
	"strings"
	"time"
)

type ApiPostRequest struct {
	Uri       string
	ApiParams binanceStructs.ApiParams
	Secrets   binanceStructs.Secrets
}

func (r *ApiPostRequest) DoRequest() (response binanceStructs.ApiResponse, err error) {
	r.ApiParams["timestamp"] = strconv.FormatInt(time.Now().UnixMilli(), 10)

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
		logger.Info(fmt.Sprintf("Send POST request to: URI: %s. Params: %s", r.Uri, paramsJsonBytes))
	}

	var requestUrl string
	requestUrl = r.getRequestData()
	requestBody := bytes.NewBuffer(nil)

	var request *http.Request
	request, err = http.NewRequest("POST", requestUrl, requestBody)
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
	request.Header.Add("X-MBX-APIKEY", r.Secrets.Key)

	var httpClient http.Client
	config := tls.Config{}
	transport := http.Transport{
		TLSClientConfig: &config,
	}
	httpClient = http.Client{
		Transport: &transport,
		Timeout:   10 * time.Second,
	}

	var resp *http.Response
	resp, err = httpClient.Do(request)
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
		errDef := Body.Close()
		if errDef != nil {
			logger.Error(fmt.Sprintf("Binance Features API error closing response body: %s", err))
		}
	}(resp.Body)

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

	if resp.Status != "200 OK" {
		msg := fmt.Sprintf("Binance Features API error. http status: %s. Body: %s", resp.Status, response)
		logger.Error(msg)
		err = tools.AppError{
			Message: msg,
		}
		return
	}

	return
}

func (r *ApiPostRequest) getRequestData() (requestUrl string) {
	queryParams := url.Values{}
	for key, value := range r.ApiParams {
		valueStr, ok := value.(string)
		if ok {
			queryParams.Add(key, valueStr)
		} else {
			valueArr, ok := value.([]binanceStructs.ApiParams)
			if ok {
				valueArrStr, err := json.Marshal(valueArr)
				if err != nil {
					logger.Warning(fmt.Sprintf("can not Marshal value to send in Binance request (%s)", valueArr))
				}
				queryParams.Add(key, string(valueArrStr))
			}
		}
	}

	var builder strings.Builder
	builder.WriteString(r.Secrets.ApiEndpoint)
	builder.WriteString(r.Uri)
	builder.WriteByte('?')
	builder.WriteString(queryParams.Encode())
	builder.WriteString("&signature=")
	builder.WriteString(getSignature(r.ApiParams, r.Secrets.Pass))
	requestUrl = builder.String()
	return
}
