package request

import (
	bybitStructs "bitbucket.org/shatylos/trader/domain/domains/bybit/structs"
	"bitbucket.org/shatylos/trader/utils"
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"time"
)

type ApiParams map[string]interface{}

func apiQueryGet(uri string, params ApiParams, secrets bybitStructs.Secrets) (interface{}, error) {
	return apiQuery(uri, params, secrets, "GET",
		func(params ApiParams) (string, *bytes.Buffer, error) {
			queryParams := url.Values{}
			for key, value := range params {
				v, ok := value.(string)
				if !ok {
					return "", nil, utils.AppError{Message: "ByBit request parameter is not a string"}
				}
				queryParams.Add(key, v)
			}
			queryString := "?" + queryParams.Encode()
			buffer := bytes.NewBuffer(nil)
			return queryString, buffer, nil
		})
}

func apiQueryPost(uri string, params ApiParams, secrets bybitStructs.Secrets) (interface{}, error) {
	return apiQuery(uri, params, secrets, "POST",
		func(params ApiParams) (string, *bytes.Buffer, error) {
			postContent, err := json.Marshal(params)
			if err != nil {
				return "", nil, err
			}
			buffer := bytes.NewBuffer(postContent)
			return "", buffer, nil
		})
}

func apiQuery(uri string, params ApiParams, secrets bybitStructs.Secrets, method string, getRequestData func(params ApiParams) (string, *bytes.Buffer, error)) (interface{}, error) {

	params["api_key"] = secrets.Key
	params["timestamp"] = fmt.Sprintf("%d", time.Now().UnixMilli())
	params["sign"] = getSignature(params, secrets.Pass)

	queryString, requestBody, err := getRequestData(params)
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest(method, secrets.ApiEndpoint+uri+queryString, requestBody)
	req.Header.Add("Content-Type", "application/json")

	config := tls.Config{}
	if utils.AppConfig("INSECURE_REQUEST") == "yes" {
		config = tls.Config{
			InsecureSkipVerify: true,
		}
	}
	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &config,
	}}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.Status != "200 OK" {
		return nil, errors.New("http status: " + resp.Status)
	}

	body, err1 := ioutil.ReadAll(resp.Body)
	if err1 != nil {
		return nil, err1
	}

	var dat map[string]interface{}
	err2 := json.Unmarshal([]byte(body), &dat)
	if err2 != nil {
		return nil, err2
	}

	if retCode, ok := dat["retCode"]; ok && retCode.(float64) != 0 {
		return nil, errors.New(dat["retMsg"].(string))
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
