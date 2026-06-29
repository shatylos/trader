package request

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/tls"
	"encoding/json"
	"fmt"
	bybitStructs "github.com/shatylos/trader/internal/domain/domains/bybit/structs"
	"github.com/shatylos/trader/tools"
	"github.com/shatylos/trader/tools/apperrors"
	"github.com/shatylos/trader/tools/logger"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

type ApiParams map[string]interface{}

var httpClient *http.Client

var RequestApiError = apperrors.New("request api error")
var LeverageNotModifiedApiError = apperrors.New("leverage has not been modified")

func getClient() *http.Client {
	if httpClient == nil {
		config := tls.Config{}
		httpClient = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &config,
			},
			Timeout: 10 * time.Second,
		}
	}
	return httpClient
}

func apiQueryGet(uri string, params ApiParams, secrets bybitStructs.Secrets) (result interface{}, err error) {
	result, err = apiQuery(uri, params, secrets, "GET",
		func(params ApiParams, apiEndpoint string, uri string) (requestUrl string, buffer *bytes.Buffer, err error) {
			queryParams := url.Values{}
			for key, value := range params {
				v, ok := value.(string)
				if !ok {
					err = apperrors.New("ByBit request parameter is not a string, key: %s, value: %s", key, value)
					return
				}
				queryParams.Add(key, v)
			}

			var builder strings.Builder
			builder.WriteString(apiEndpoint)
			builder.WriteString(uri)
			builder.WriteByte('?')
			builder.WriteString(queryParams.Encode())
			requestUrl = builder.String()

			queryParams = nil
			buffer = bytes.NewBuffer(nil)
			return
		})
	if err != nil {
		err = apperrors.Wrap(err, "error sending bybit GET request, uri: %s, params: %s", uri, params)
		return
	}
	return
}

func apiQueryPost(uri string, params ApiParams, secrets bybitStructs.Secrets) (result interface{}, err error) {
	result, err = apiQuery(uri, params, secrets, "POST",
		func(params ApiParams, apiEndpoint string, uri string) (requestUrl string, buffer *bytes.Buffer, err error) {
			var postContent []byte
			postContent, err = json.Marshal(params)
			if err != nil {
				err = apperrors.Wrap(err, "error marshal params: %s", params)
				return
			}
			buffer = bytes.NewBuffer(postContent)

			var builder strings.Builder
			builder.WriteString(apiEndpoint)
			builder.WriteString(uri)
			builder.WriteByte('?')
			requestUrl = builder.String()

			return
		})
	if err != nil {
		err = apperrors.Wrap(err, "error sending bybit POST request, uri: %s, params: %s", uri, params)
		return
	}
	return
}

func apiQuery(uri string, params ApiParams, secrets bybitStructs.Secrets, method string, getRequestData func(params ApiParams, apiEndpoint string, uri string) (string, *bytes.Buffer, error)) (result interface{}, err error) {

	if secrets.Verbose {
		var paramsJsonBytes []byte
		paramsJsonBytes, err = json.Marshal(params)
		if err != nil {
			err = apperrors.Wrap(err, "error marshalling params to JSON, params: %s", params)
			return
		}
		logger.Info(fmt.Sprintf("Send %s request to: URI: %s. Params: %s", method, uri, paramsJsonBytes))
	}

	params["api_key"] = secrets.Key
	params["timestamp"] = fmt.Sprintf("%d", time.Now().UnixMilli())
	params["sign"] = getSignature(params, secrets.Pass)

	var requestUrl string
	var requestBody *bytes.Buffer
	requestUrl, requestBody, err = getRequestData(params, secrets.ApiEndpoint, uri)
	if err != nil {
		err = apperrors.Wrap(err, "error get request data, uri: %s, params: %s", uri, params)
		return
	}

	req, _ := http.NewRequest(method, requestUrl, requestBody)
	req.Header.Add("Content-Type", "application/json")

	client := getClient()
	resp, err := client.Do(req)
	if err != nil {
		err = apperrors.New("ByBit API error get request")
		return
	}
	defer resp.Body.Close()

	if resp.Status != "200 OK" {
		return nil, tools.AppError{
			Message: fmt.Sprintf("ByBit API error. http status: %s", resp.Status),
		}
	}

	body := &bytes.Buffer{}
	_, err1 := io.Copy(body, resp.Body)
	if err1 != nil {
		return nil, tools.AppError{
			Message:     fmt.Sprintf("ByBit API error read body: %s", err1.Error()),
			ParentError: err,
		}
	}

	var dat map[string]interface{}
	err2 := json.Unmarshal(body.Bytes(), &dat)
	body.Reset()
	if err2 != nil {
		return nil, tools.AppError{
			Message:     fmt.Sprintf("ByBit API error unmarshal data: %s", err2.Error()),
			ParentError: err,
		}
	}

	if retCode, ok := dat["retCode"]; ok && retCode.(float64) != 0 {
		switch retCode {
		case 110043.0:
			err = apperrors.Wrap(LeverageNotModifiedApiError, "ByBit API error: %s", dat["retMsg"])
		default:
			err = apperrors.WrapExcuse(RequestApiError, "ByBit API error: %s", dat["retMsg"])
		}
		return nil, err
	}

	result = dat["result"].(interface{})
	return
}

func getSignature(params ApiParams, key string) string {
	keys := make([]string, len(params))
	i := 0
	_val := ""
	for k, _ := range params {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	for _, k := range keys {
		strValue, isString := params[k].(string)
		if isString {
			_val += k + "=" + strValue + "&"
		} else {
			boolValue, _ := params[k].(bool)
			if boolValue {
				_val += k + "=true&"
			} else {
				_val += k + "=false&"
			}
		}
	}
	_val = _val[0 : len(_val)-1]
	h := hmac.New(sha256.New, []byte(key))
	_, err := io.WriteString(h, _val)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}
