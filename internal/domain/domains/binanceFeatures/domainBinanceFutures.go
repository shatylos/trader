package binanceFeatures

import (
	binanceStructs "github.com/shatylos/trader/internal/domain/domains/binanceFeatures/structs"
	providerStructs "github.com/shatylos/trader/internal/domain/structs"
	"github.com/shatylos/trader/tools"
	"github.com/shatylos/trader/tools/logger"
	_type "github.com/shatylos/trader/tools/type"
)

type DomainBinanceFutures struct {
	code    string
	secrets binanceStructs.Secrets
}

func (d *DomainBinanceFutures) GetCode() string {
	return d.code
}

func (d *DomainBinanceFutures) SetConfig(config map[interface{}]interface{}) (err error) {
	domainCode, err := _type.ToString(config["code"])
	if err != nil {
		return tools.AppError{
			Message: "The field code is empty or contains not correct value type. In DomainBinanceFutures config. Expects a string",
		}
	}
	d.code = domainCode

	secretMap, ok := config["secrets"].(map[interface{}]interface{})
	if !ok {
		return tools.AppError{
			Message: "The field secrets is empty or contains not correct value type. In DomainBinanceFutures config. Expects a map with \"key\", \"pass\" and \"endpoint\" keys",
		}
	}

	secrets := binanceStructs.Secrets{}

	apiEndpoint, err := _type.ToString(secretMap["endpoint"])
	if err != nil {
		return tools.AppError{
			Message: "The field secrets.endpoint is empty or contains not correct value type. In DomainBinanceFutures config. Expects a string",
		}
	}
	secrets.ApiEndpoint = apiEndpoint

	key, err := _type.ToString(secretMap["key"])
	if err != nil {
		return tools.AppError{
			Message: "The field secrets.key is empty or contains not correct value type. In DomainBinanceFutures config. Expects a string",
		}
	}
	secrets.Key = key

	pass, err := _type.ToString(secretMap["pass"])
	if err != nil {
		return tools.AppError{
			Message: "The field secrets.pass is empty or contains not correct value type. In DomainBinanceFutures config. Expects a string",
		}
	}
	secrets.Pass = pass

	var verbose int64
	verbose, err = _type.ToInt64(secretMap["verbose"])
	if err != nil {
		logger.Error("The field verbose in domain config is empty or contains not correct value type. Expects 1 or 0")
		err = tools.AppError{
			Message:     "The field verbose in domain config is empty or contains not correct value type. Expects 1 or 0",
			ParentError: err,
		}
		return err
	}
	if verbose == 1 {
		secrets.Verbose = true
	} else {
		secrets.Verbose = false
	}

	d.secrets = secrets

	return
}

func (d *DomainBinanceFutures) OpenPosition(positionRequest providerStructs.DomainPositionRequest) (string, error) {
	panic("Not implemented")
}

func (d *DomainBinanceFutures) ModifyTpSl(tpSlRequest providerStructs.TpSlRequest) error {
	panic("Not implemented")
}
