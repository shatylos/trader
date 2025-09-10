package investor

import (
	"github.com/shatylos/trader/internal/domain"
	"github.com/shatylos/trader/internal/strategy/investor/storage"
	"github.com/shatylos/trader/tools"
	"github.com/shatylos/trader/tools/logger"
	_type "github.com/shatylos/trader/tools/type"
	"time"
)

type Config struct {
	Id               string
	Enabled          bool
	Verbose          bool
	TelegramNotifier bool
	CoinPare         string
	MainCurrency     string
	TradeCurrency    string
	TimeoutSeconds   time.Duration
	QtyPrecision     int64
	PricePrecision   int64
	WithdrawPercent  float64
	RequestDelay     time.Duration
}

type TimeframeConfig struct {
	Resolution                  string
	QtyPercent                  float64
	CandleReview                int64
	SidewaysMinCandlesAmount    int64
	SidewaysPercentToPrice      float64
	SidewaysPremiumCoefficient  float64
	SidewaysDiscountCoefficient float64
	IsHeap                      bool
}

func (i *Investor) SetConfig(strategyConfig interface{}, domainConfig map[string]interface{}) (err error) {

	i.config = Config{}

	configMap, ok := strategyConfig.(map[interface{}]interface{})
	if !ok {
		logger.Error("Config for the strategy investor is not valid")
		err = tools.AppError{
			Message: "Config for the strategy investor is not valid",
		}
		return
	}

	i.config.Id, err = _type.ToString(configMap["id"])
	if err != nil {
		logger.Error("The field id is empty or contains not correct value type. Expects string value")
		err = tools.AppError{
			Message:     "The field id is empty or contains not correct value type. Expects string value",
			ParentError: err,
		}
		return
	}

	var enabled int64
	enabled, err = _type.ToInt64(configMap["enabled"])
	if err != nil {
		logger.Error("The field enabled is empty or contains not correct value type. Expects 1 or 0")
		err = tools.AppError{
			Message:     "The field enabled is empty or contains not correct value type. Expects 1 or 0",
			ParentError: err,
		}
		return
	}
	if enabled == 1 {
		i.config.Enabled = true
	} else {
		i.config.Enabled = false
	}

	var verbose int64
	verbose, err = _type.ToInt64(configMap["verbose"])
	if err != nil {
		logger.Error("The field verbose in strategy config is empty or contains not correct value type. Expects 1 or 0")
		err = tools.AppError{
			Message:     "The field verbose in strategy config is empty or contains not correct value type. Expects 1 or 0",
			ParentError: err,
		}
		return
	}
	if verbose == 1 {
		i.config.Verbose = true
	} else {
		i.config.Verbose = false
	}

	if configMap["telegram_notifier"] != nil {
		var tgNotifier int64
		tgNotifier, err = _type.ToInt64(configMap["telegram_notifier"])
		if err != nil {
			logger.Error("The field telegram_notifier in strategy config is empty or contains not correct value type. Expects 1 or 0")
			err = tools.AppError{
				Message:     "The field telegram_notifier in strategy config is empty or contains not correct value type. Expects 1 or 0",
				ParentError: err,
			}
			return
		}
		if tgNotifier == 1 {
			i.config.TelegramNotifier = true
		}
	}

	i.config.CoinPare, err = _type.ToString(configMap["coin_pare"])
	if err != nil {
		err = tools.AppError{
			Message:     "The field coin_pare is empty or contains not correct value type",
			ParentError: err,
		}
		return
	}

	i.config.MainCurrency, err = _type.ToString(configMap["main_currency"])
	if err != nil {
		return tools.AppError{
			Message:     "The field main_currency is empty or contains not correct value type",
			ParentError: err,
		}
	}

	i.config.TradeCurrency, err = _type.ToString(configMap["trade_currency"])
	if err != nil {
		return tools.AppError{
			Message:     "The field trade_currency is empty or contains not correct value type",
			ParentError: err,
		}
	}

	i.config.TimeoutSeconds, err = _type.ToTimeDuration(configMap["timeout_seconds"])
	if err != nil {
		err = tools.AppError{
			Message:     "The field timeout_seconds is empty or contains not correct value",
			ParentError: err,
		}
		return
	}

	i.config.QtyPrecision, err = _type.ToInt64(configMap["qty_precision"])
	if err != nil {
		return tools.AppError{
			Message:     "The field qty_precision is empty or contains not correct value type. Expects int64 value",
			ParentError: err,
		}
	}

	i.config.PricePrecision, err = _type.ToInt64(configMap["price_precision"])
	if err != nil {
		return tools.AppError{
			Message:     "The field price_precision is empty or contains not correct value type. Expects int64 value",
			ParentError: err,
		}
	}

	i.config.WithdrawPercent, err = _type.ToFloat64(configMap["withdraw_percent"])
	if err != nil {
		return tools.AppError{
			Message:     "The field withdraw_percent is empty",
			ParentError: err,
		}
	}

	i.config.RequestDelay, err = _type.ToTimeDuration(configMap["request_delay"])
	if err != nil {
		return tools.AppError{
			Message:     "The field request_delay is empty",
			ParentError: err,
		}
	}

	err = i.applyTimeframeConfig(configMap["timeframes"])
	if err != nil {
		return
	}

	err = i.applyDomainConfig(configMap, domainConfig)
	if err != nil {
		return
	}

	i.Storage = storage.Storage{Id: i.config.Id}
	err = i.Storage.InitStorage()
	if err != nil {
		return
	}

	return
}

func (i *Investor) applyTimeframeConfig(timeframesConfig interface{}) (err error) {

	rawTimeframes, ok := timeframesConfig.([]interface{})
	if !ok {
		logger.Error("Config for the strategy Investor is not valid. \"timeframes\" must be array")
		err = tools.AppError{
			Message: "Config for the strategy Investor is not valid. \"timeframes\" must be array",
		}
		return
	}

	isHeapSet := false

	for _, tf := range rawTimeframes {
		tfMap, ok := tf.(map[interface{}]interface{})
		if !ok {
			logger.Error("Config for the strategy Investor is not valid. Item of \"timeframes\" must be map ok key value")
			err = tools.AppError{
				Message: "Config for the strategy Investor is not valid. Item of \"timeframes\" must be map ok key value",
			}
			return
		}

		timeframe := Timeframe{
			Config: TimeframeConfig{},
		}

		timeframe.Config.Resolution, err = _type.ToString(tfMap["resolution"])
		if err != nil {
			return tools.AppError{
				Message:     "The field resolution is empty or contains not correct value type",
				ParentError: err,
			}
		}

		timeframe.Config.QtyPercent, err = _type.ToFloat64(tfMap["qty_percent"])
		if err != nil || timeframe.Config.QtyPercent == 0 {
			return tools.AppError{
				Message:     "Empty value qty_percent",
				ParentError: err,
			}
		}

		timeframe.Config.CandleReview, err = _type.ToInt64(tfMap["candle_review"])
		if err != nil || timeframe.Config.CandleReview < 10 {
			return tools.AppError{
				Message:     "The field candle_review is empty or contains not correct value type. Expects int64 value more than 10",
				ParentError: err,
			}
		}

		timeframe.Config.SidewaysMinCandlesAmount, err = _type.ToInt64(tfMap["sideways_min_candles_amount"])
		if err != nil || timeframe.Config.SidewaysMinCandlesAmount == 0 {
			return tools.AppError{
				Message:     "Empty value sideways_min_candles_amount",
				ParentError: err,
			}
		}

		timeframe.Config.SidewaysPercentToPrice, err = _type.ToFloat64(tfMap["sideways_percent_to_price"])
		if err != nil || timeframe.Config.SidewaysPercentToPrice == 0 {
			return tools.AppError{
				Message:     "Empty value sideways_percent_to_price",
				ParentError: err,
			}
		}

		timeframe.Config.SidewaysPremiumCoefficient, err = _type.ToFloat64(tfMap["sideways_premium_coefficient"])
		if err != nil {
			return tools.AppError{
				Message:     "Empty value sideways_premium_coefficient",
				ParentError: err,
			}
		}

		timeframe.Config.SidewaysDiscountCoefficient, err = _type.ToFloat64(tfMap["sideways_discount_coefficient"])
		if err != nil {
			return tools.AppError{
				Message:     "Empty value sideways_discount_coefficient",
				ParentError: err,
			}
		}

		var heap int64
		heap, _ = _type.ToInt64(tfMap["heap"])
		if heap == 1 {
			if isHeapSet {
				message := "Heap must be set only for one timeframe"
				logger.Error(message)
				err = tools.AppError{
					Message: message,
				}
				return
			}
			timeframe.Config.IsHeap = true
			isHeapSet = true
		}

		i.Timeframes = append(i.Timeframes, timeframe)
	}

	if !isHeapSet {
		message := "Heap timeframe must be set"
		logger.Error(message)
		err = tools.AppError{
			Message: message,
		}
		return
	}

	return
}

func (i *Investor) applyDomainConfig(configMap map[interface{}]interface{}, domainConfig map[string]interface{}) error {
	domainId, err := _type.ToString(configMap["domain"])
	if err != nil {
		return tools.AppError{
			Message:     "The field domain is empty or contains not correct value type",
			ParentError: err,
		}
	}
	domainConfigItem, ok := domainConfig[domainId].(map[interface{}]interface{})
	if !ok {
		return tools.AppError{
			Message: "The domain config is not valid. Domain value should be related to the domain config item",
		}
	}
	domainCode, err := _type.ToString(domainConfigItem["code"])
	if err != nil {
		return tools.AppError{
			Message:     "The field code is empty in domain config or contains not correct value type",
			ParentError: err,
		}
	}
	domainItem, err := domain.GetSpotDomain(domainCode)
	if err != nil {
		return err
	}

	err = domainItem.SetConfig(domainConfigItem)
	if err != nil {
		return err
	}
	i.provider = domainItem

	return nil
}
