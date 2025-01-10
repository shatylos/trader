package request

import (
	"github.com/shatylos/trader/utils"
	"strconv"
)

type UserInfoStruct struct {
	Uid        int64
	ServerDate int64
	Balances   map[string]float64
	Reserved   map[string]float64
}

func GetUserInfo() (*UserInfoStruct, error) {
	params := make(ApiParams, 0)
	queryResp, er := apiQuery("user_info", params)
	if er != nil {
		return nil, er
	}
	return mapUserInfo(queryResp)
}

func mapUserInfo(source map[string]interface{}) (*UserInfoStruct, error) {

	result := UserInfoStruct{
		Balances: map[string]float64{},
		Reserved: map[string]float64{},
	}

	uid, ok := source["uid"].(float64)
	if !ok {
		return nil, utils.AppError{Message: "Can not parse broker response. Bad uid."}
	}
	result.Uid = int64(uid)

	serverDate, ok := source["server_date"].(float64)
	if !ok {
		return nil, utils.AppError{Message: "Can not parse broker response. Bad server_date."}
	}
	result.ServerDate = int64(serverDate)

	balances, ok := source["balances"].(map[string]interface{})
	if !ok {
		return nil, utils.AppError{Message: "Can not parse broker response. Bad balances."}
	}
	for key, val := range balances {
		balance, ok := val.(string)
		if !ok {
			return nil, utils.AppError{Message: "Can not parse broker response. Bad balance value."}
		}
		balanceFloat, er := strconv.ParseFloat(balance, 8)
		if er != nil {
			return nil, er
		}
		result.Balances[key] = balanceFloat
	}

	reserved, ok := source["reserved"].(map[string]interface{})
	if !ok {
		return nil, utils.AppError{Message: "Can not parse broker response. Bad reserved."}
	}
	for key, val := range reserved {
		reservedItem, ok := val.(string)
		if !ok {
			return nil, utils.AppError{Message: "Can not parse broker response. Bad reserved value."}
		}
		reserveFloat, er := strconv.ParseFloat(reservedItem, 8)
		if er != nil {
			return nil, er
		}
		result.Reserved[key] = reserveFloat
	}

	return &result, nil
}
