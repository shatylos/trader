package request

import (
	"bitbucket.org/shatylos/trader/utils"
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

type ApiParams map[string]string

func apiQuery(method string, params ApiParams) (map[string]interface{}, error) {
	key := os.Getenv("TRADER_EXMO_KEY")
	secret := os.Getenv("TRADER_EXMO_SECRET")

	if key == "" || secret == "" {
		return nil, utils.AppError{Message: "Exmo env variables TRADER_EXMO_KEY, TRADER_EXMO_SECRET are not set"}
	}

	postParams := url.Values{}
	postParams.Add("nonce", nonce())
	if params != nil {
		for key, value := range params {
			postParams.Add(key, value)
		}
	}
	postContent := postParams.Encode()

	sign := doSign(postContent, secret)

	req, _ := http.NewRequest("POST", "https://api.exmo.com/v1.1/"+method, bytes.NewBuffer([]byte(postContent)))
	req.Header.Set("Key", key)
	req.Header.Set("Sign", sign)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(postContent)))

	client := &http.Client{}
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

	if result, ok := dat["result"]; ok && result.(bool) != true {
		return nil, errors.New(dat["error"].(string))
	}

	return dat, nil
}

func nonce() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func doSign(message string, secret string) string {
	mac := hmac.New(sha512.New, []byte(secret))
	mac.Write([]byte(message))
	return fmt.Sprintf("%x", mac.Sum(nil))
}
