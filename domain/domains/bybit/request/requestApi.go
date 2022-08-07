package request

import (
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

type ApiParams map[string]string

const testEndpoint = "https://api-testnet.bybit.com"
const liveEndpoint = "https://api.bybit.com"

func apiQueryGet(method string, params ApiParams, isDemo bool) (interface{}, error) {
	key := utils.AppConfig("TRADER_BYBIT_API_KEY")
	secret := utils.AppConfig("TRADER_BYBIT_SECRET_API_KEY")

	if key == "" || secret == "" {
		return nil, utils.AppError{Message: "Bybit env variables TRADER_BYBIT_API_KEY, TRADER_BYBIT_SECRET_API_KEY are not set"}
	}

	params["api_key"] = key
	params["timestamp"] = fmt.Sprintf("%d", time.Now().UnixMilli())
	params["sign"] = getSignature(params, secret)

	queryParams := url.Values{}
	for key, value := range params {
		queryParams.Add(key, value)
	}
	queryString := queryParams.Encode()

	endpoint := liveEndpoint
	if isDemo {
		endpoint = testEndpoint
	}

	req, _ := http.NewRequest("GET", endpoint+method+"?"+queryString, bytes.NewBuffer(nil))

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

	if retCode, ok := dat["ret_code"]; ok && retCode.(float64) != 0 {
		return nil, errors.New(dat["ret_msg"].(string))
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
		_val += k + "=" + params[k] + "&"
	}
	_val = _val[0 : len(_val)-1]
	h := hmac.New(sha256.New, []byte(key))
	_, err := io.WriteString(h, _val)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}
