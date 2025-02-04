package request

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/tls"
	"encoding/json"
	"fmt"
	bybitStructs "github.com/shatylos/trader/domain/domains/bybit/structs"
	"github.com/shatylos/trader/utils"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

type ApiParams map[string]interface{}

var httpClient *http.Client

func getClient() *http.Client {
	if httpClient == nil {
		config := tls.Config{}
		if utils.AppConfig("INSECURE_REQUEST") == "yes" {
			config = tls.Config{
				InsecureSkipVerify: true,
			}
		}
		httpClient = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &config,
			},
			Timeout: 10 * time.Second,
		}
	}
	return httpClient
}

func apiQueryGet(uri string, params ApiParams, secrets bybitStructs.Secrets) (interface{}, *utils.AppError) {
	return apiQuery(uri, params, secrets, "GET",
		func(params ApiParams, apiEndpoint string, uri string) (string, *bytes.Buffer, error) {
			queryParams := url.Values{}
			for key, value := range params {
				v, ok := value.(string)
				if !ok {
					return "", nil, utils.AppError{Message: "ByBit request parameter is not a string"}
				}
				queryParams.Add(key, v)
			}

			var builder strings.Builder
			builder.WriteString(apiEndpoint)
			builder.WriteString(uri)
			builder.WriteByte('?')
			builder.WriteString(queryParams.Encode())
			requestUrl := builder.String()

			queryParams = nil
			buffer := bytes.NewBuffer(nil)
			return requestUrl, buffer, nil
		})
}

func apiQueryPost(uri string, params ApiParams, secrets bybitStructs.Secrets) (interface{}, *utils.AppError) {
	return apiQuery(uri, params, secrets, "POST",
		func(params ApiParams, apiEndpoint string, uri string) (string, *bytes.Buffer, error) {
			postContent, err := json.Marshal(params)
			if err != nil {
				return "", nil, err
			}
			buffer := bytes.NewBuffer(postContent)

			var builder strings.Builder
			builder.WriteString(apiEndpoint)
			builder.WriteString(uri)
			builder.WriteByte('?')
			requestUrl := builder.String()

			return requestUrl, buffer, nil
		})
}

func apiQuery(uri string, params ApiParams, secrets bybitStructs.Secrets, method string, getRequestData func(params ApiParams, apiEndpoint string, uri string) (string, *bytes.Buffer, error)) (interface{}, *utils.AppError) {

	params["api_key"] = secrets.Key
	params["timestamp"] = fmt.Sprintf("%d", time.Now().UnixMilli())
	params["sign"] = getSignature(params, secrets.Pass)

	requestUrl, requestBody, err := getRequestData(params, secrets.ApiEndpoint, uri)
	if err != nil {
		return nil, &utils.AppError{
			Message:     fmt.Sprintf("ByBit API error prepare request query: %s", err.Error()),
			ParentError: err,
		}
	}

	req, _ := http.NewRequest(method, requestUrl, requestBody)
	req.Header.Add("Content-Type", "application/json")

	client := getClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, &utils.AppError{
			Message:     fmt.Sprintf("ByBit API error get request: %s", err.Error()),
			ParentError: err,
		}
	}
	defer resp.Body.Close()

	if resp.Status != "200 OK" {
		return nil, &utils.AppError{
			Message: fmt.Sprintf("ByBit API error. http status: %s", resp.Status),
		}
	}

	body := &bytes.Buffer{}
	_, err1 := io.Copy(body, resp.Body)
	if err1 != nil {
		return nil, &utils.AppError{
			Message:     fmt.Sprintf("ByBit API error read body: %s", err1.Error()),
			ParentError: err,
		}
	}

	var dat map[string]interface{}
	err2 := json.Unmarshal(body.Bytes(), &dat)
	body.Reset()
	if err2 != nil {
		return nil, &utils.AppError{
			Message:     fmt.Sprintf("ByBit API error unmarshal data: %s", err2.Error()),
			ParentError: err,
		}
	}

	if retCode, ok := dat["retCode"]; ok && retCode.(float64) != 0 {
		return nil, &utils.AppError{
			Message: fmt.Sprintf("ByBit API error: %s", dat["retMsg"].(string)),
			Code:    retCode.(float64),
		}
	}

	return dat["result"].(interface{}), nil
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
