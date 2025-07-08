package request

import (
	"encoding/json"
	"github.com/shatylos/trader/internal/domain/domains/binanceFeatures/structs"
	"github.com/shatylos/trader/tools/logger"
	_type "github.com/shatylos/trader/tools/type"
)

func IsErrorCode(rawResponse structs.ApiResponse, expectErrorCode int64) (isError bool) {
	var jsonResponse interface{}
	jsonErr := json.Unmarshal(rawResponse, &jsonResponse)
	if jsonErr != nil {
		logger.Error(jsonErr.Error())
		return
	}

	queryRespMap, ok := jsonResponse.(map[string]interface{})
	if !ok {
		return
	}
	codeInterface, ok := queryRespMap["code"].(interface{})
	if !ok {
		return
	}
	code, toIntErr := _type.ToInt64(codeInterface)
	if toIntErr != nil {
		return
	}
	if expectErrorCode == code {
		isError = true
	}
	return
}
