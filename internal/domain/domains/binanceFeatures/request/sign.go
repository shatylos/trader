package request

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	binanceStructs "github.com/shatylos/trader/internal/domain/domains/binanceFeatures/structs"
	"io"
	"net/url"
)

func getSignature(params binanceStructs.ApiParams, key string) (sign string) {
	queryParams := url.Values{}
	for k, value := range params {
		queryParams.Add(k, value)
	}

	queryStr := queryParams.Encode()
	h := hmac.New(sha256.New, []byte(key))
	_, err := io.WriteString(h, queryStr)
	if err != nil {
		return ""
	}
	sign = fmt.Sprintf("%x", h.Sum(nil))

	return
}
