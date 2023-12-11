package utils

import (
	"strconv"
	"time"
)

func ToInt64(value interface{}) (int64, error) {
	switch value.(type) {
	case int:
		return int64(value.(int)), nil
	case int64:
		return value.(int64), nil
	case float64:
		return int64(value.(float64)), nil
	case string:
		return strconv.ParseInt(value.(string), 10, 64)
	}
	return 0, AppError{
		Message: "The value can not be converted to int64",
	}
}

func ToFloat64(value interface{}) (float64, error) {
	if value == nil {
		return 0, EmptyValueError
	}
	switch value.(type) {
	case int:
		return float64(value.(int)), nil
	case int64:
		return float64(value.(int64)), nil
	case float64:
		return value.(float64), nil
	case string:
		return strconv.ParseFloat(value.(string), 64)
	}
	return 0, AppError{
		Message: "The value can not be converted to float64",
	}
}

func ToString(value interface{}) (string, error) {
	switch value.(type) {
	case string:
		return value.(string), nil
	case int:
		return strconv.Itoa(value.(int)), nil
	case int64:
		return strconv.FormatInt(value.(int64), 10), nil
	case float64:
		return strconv.FormatFloat(value.(float64), 'f', -1, 64), nil
	}

	return "", AppError{
		Message: "The value can not be converted to string",
	}
}

func ToTimeDuration(value interface{}) (time.Duration, error) {
	switch value.(type) {
	case int:
		return time.Duration(value.(int)), nil
	case int64:
		return time.Duration(value.(int64)), nil
	case float64:
		return time.Duration(value.(float64)), nil
	case string:
		return time.ParseDuration(value.(string))
	}
	return 0, AppError{
		Message: "The value can not be converted to time.Duration",
	}
}

func ToInt64Slice(value interface{}) ([]int64, error) {
	switch value.(type) {
	case []interface{}:
		var result []int64
		for _, item := range value.([]interface{}) {
			itemInt64, err := ToInt64(item)
			if err != nil {
				return nil, err
			}
			result = append(result, itemInt64)
		}
		return result, nil
	}
	return nil, AppError{
		Message: "The value can not be converted to []int64",
	}
}
