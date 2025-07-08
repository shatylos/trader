package request

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	binanceStructs "github.com/shatylos/trader/internal/domain/domains/binanceFeatures/structs"
	"github.com/shatylos/trader/tools/logger"
	"io"
	"net/url"
)

func getSignature(params binanceStructs.ApiParams, key string) (sign string) {
	queryParams := url.Values{}
	for k, value := range params {
		valueStr, ok := value.(string)
		if ok {
			queryParams.Add(k, valueStr)
		} else {
			valueArr, ok := value.([]binanceStructs.ApiParams)
			if ok {
				valueArrStr, err := json.Marshal(valueArr)
				if err != nil {
					logger.Warning(fmt.Sprintf("can not Marshal value to send in Binance request (%s)", valueArr))
				}
				queryParams.Add(k, string(valueArrStr))
			}
		}
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
