package fibonacci

import (
	"github.com/shatylos/trader/internal/domain"
	"github.com/shatylos/trader/tools"
	"github.com/shatylos/trader/tools/logger"
	_type "github.com/shatylos/trader/tools/type"
	"time"
)

type Config struct {
	Enabled               bool
	Verbose               bool
	TelegramNotifier      bool
	MaxCandleReview       int64
	MinCandleReview       int64
	Code                  string
	CoinPare              string
	MainCurrency          string
	TradeCurrency         string
	MinDepo               float64
	Leverage              int64
	Id                    string
	Resolution            string
	ResolutionMins        int64
	LongTrendResolution   string
	LongTrendCandleReview int64
	TimeoutSeconds        time.Duration
	QtyPrecision          int64
	MinQty                float64
	PricePrecision        int64
	FibEntryPoint1        float64
	FibEntryPoint2        float64
	FibEntryPoint3        float64
	FibStopLoss           float64
	FibTakeProfit1        float64
	FibTakeProfit2        float64
	FibTakeProfit3        float64
	RiskPercent           float64
	EP1ToFullQtyPercent   float64
	EP2ToFullQtyPercent   float64
	EP3ToFullQtyPercent   float64
	HoursToReduceTP1      int64
	HoursToReduceTP2      int64
	HoursToReduceTP3      int64
	PercentToReduceTP     int64
	WithdrawPercent       float64
}

var resolutions = map[string]int64{
	"1":   1,
	"3":   3,
	"5":   5,
	"15":  15,
	"30":  30,
	"60":  60,
	"120": 120,
	"240": 240,
	"360": 360,
	"720": 720,
	"D":   86400,
	"W":   604800,
}

func (f *Fibonacci) SetConfig(strategyConfig interface{}, domainConfig map[string]interface{}) (err error) {

	f.config = Config{}

	configMap, ok := strategyConfig.(map[interface{}]interface{})
	if !ok {
		logger.Error("Config for the strategy fibonacci is not valid")
		err = tools.AppError{
			Message: "Config for the strategy fibonacci is not valid",
		}
		return
	}

	f.config.Id, err = _type.ToString(configMap["id"])
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
		f.config.Enabled = true
	} else {
		f.config.Enabled = false
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
		f.config.Verbose = true
	} else {
		f.config.Verbose = false
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
			f.config.TelegramNotifier = true
		}
	}

	f.config.CoinPare, err = _type.ToString(configMap["coin_pare"])
	if err != nil {
		err = tools.AppError{
			Message:     "The field coin_pare is empty or contains not correct value type",
			ParentError: err,
		}
		return
	}

	f.config.MainCurrency, err = _type.ToString(configMap["main_currency"])
	if err != nil {
		return tools.AppError{
			Message:     "The field main_currency is empty or contains not correct value type",
			ParentError: err,
		}
	}

	f.config.TradeCurrency, err = _type.ToString(configMap["trade_currency"])
	if err != nil {
		return tools.AppError{
			Message:     "The field trade_currency is empty or contains not correct value type",
			ParentError: err,
		}
	}

	f.config.MinDepo, err = _type.ToFloat64(configMap["min_depo"])
	if err != nil {
		return tools.AppError{
			Message:     "Wrong value min_depo",
			ParentError: err,
		}
	}

	f.config.Leverage, err = _type.ToInt64(configMap["leverage"])
	if err != nil || f.config.Leverage == 0 {
		return tools.AppError{
			Message:     "Wrong value leverage",
			ParentError: err,
		}
	}

	f.config.TimeoutSeconds, err = _type.ToTimeDuration(configMap["timeout_seconds"])
	if err != nil {
		err = tools.AppError{
			Message:     "The field timeout_seconds is empty or contains not correct value",
			ParentError: err,
		}
		return
	}

	f.config.Resolution, err = _type.ToString(configMap["resolution"])
	if err != nil {
		return tools.AppError{
			Message:     "The field resolution is empty or contains not correct value type",
			ParentError: err,
		}
	}
	f.config.ResolutionMins, ok = resolutions[f.config.Resolution]
	if !ok {
		return tools.AppError{
			Message: "The field contains not correct value",
		}
	}

	f.config.MaxCandleReview, err = _type.ToInt64(configMap["max_candle_review"])
	if err != nil || f.config.MaxCandleReview < 10 {
		return tools.AppError{
			Message:     "The field max_candle_review is empty or contains not correct value type. Expects int64 value more than 10",
			ParentError: err,
		}
	}

	f.config.LongTrendResolution, err = _type.ToString(configMap["long_trend_resolution"])
	if err != nil {
		return tools.AppError{
			Message:     "The field long_trend_resolution is empty or contains not correct value type",
			ParentError: err,
		}
	}

	f.config.LongTrendCandleReview, err = _type.ToInt64(configMap["long_trend_candle_review"])
	if err != nil || f.config.MaxCandleReview < 10 {
		return tools.AppError{
			Message:     "The field long_trend_candle_review is empty or contains not correct value type. Expects int64 value more than 10",
			ParentError: err,
		}
	}

	f.config.MinCandleReview, err = _type.ToInt64(configMap["min_candle_review"])
	if err != nil || f.config.MinCandleReview < 1 {
		return tools.AppError{
			Message:     "The field min_candle_review is empty or contains not correct value type. Expects int64 value more than 1",
			ParentError: err,
		}
	}

	f.config.FibEntryPoint1, err = _type.ToFloat64(configMap["fib_entry_point_1"])
	if err != nil {
		return tools.AppError{
			Message:     "Wrong value fib_entry_point_1",
			ParentError: err,
		}
	}
	f.config.FibEntryPoint2, err = _type.ToFloat64(configMap["fib_entry_point_2"])
	if err != nil {
		return tools.AppError{
			Message:     "Wrong value fib_entry_point_2",
			ParentError: err,
		}
	}
	f.config.FibEntryPoint3, err = _type.ToFloat64(configMap["fib_entry_point_3"])
	if err != nil {
		return tools.AppError{
			Message:     "Wrong value fib_entry_point_3",
			ParentError: err,
		}
	}
	f.config.FibStopLoss, err = _type.ToFloat64(configMap["fib_stop_loss"])
	if err != nil {
		return tools.AppError{
			Message:     "Wrong value fib_stop_loss",
			ParentError: err,
		}
	}
	f.config.FibTakeProfit1, err = _type.ToFloat64(configMap["fib_take_profit_1"])
	if err != nil {
		return tools.AppError{
			Message:     "Wrong value fib_take_profit_1",
			ParentError: err,
		}
	}
	f.config.FibTakeProfit2, err = _type.ToFloat64(configMap["fib_take_profit_2"])
	if err != nil {
		return tools.AppError{
			Message:     "Wrong value fib_take_profit_2",
			ParentError: err,
		}
	}
	f.config.FibTakeProfit3, err = _type.ToFloat64(configMap["fib_take_profit_3"])
	if err != nil {
		return tools.AppError{
			Message:     "Wrong value fib_take_profit_3",
			ParentError: err,
		}
	}

	f.config.PricePrecision, err = _type.ToInt64(configMap["price_precision"])
	if err != nil {
		return tools.AppError{
			Message:     "The field price_precision is empty or contains not correct value type. Expects int64 value",
			ParentError: err,
		}
	}

	f.config.QtyPrecision, err = _type.ToInt64(configMap["qty_precision"])
	if err != nil {
		return tools.AppError{
			Message:     "The field qty_precision is empty or contains not correct value type. Expects int64 value",
			ParentError: err,
		}
	}

	f.config.MinQty, err = _type.ToFloat64(configMap["min_qty"])
	if err != nil {
		return tools.AppError{
			Message:     "The field min_qty is empty or contains not correct value type. Expects ToFloat64 value",
			ParentError: err,
		}
	}

	f.config.RiskPercent, err = _type.ToFloat64(configMap["risk_percent"])
	if err != nil || f.config.RiskPercent == 0 {
		return tools.AppError{
			Message:     "Empty value risk_percent",
			ParentError: err,
		}
	}
	f.config.EP1ToFullQtyPercent, err = _type.ToFloat64(configMap["ep_1_to_full_qty_percent"])
	if err != nil || f.config.EP1ToFullQtyPercent == 0 {
		return tools.AppError{
			Message:     "Empty value ep_1_to_full_qty_percent",
			ParentError: err,
		}
	}
	f.config.EP2ToFullQtyPercent, err = _type.ToFloat64(configMap["ep_2_to_full_qty_percent"])
	if err != nil || f.config.EP2ToFullQtyPercent == 0 {
		return tools.AppError{
			Message:     "Empty value ep_2_to_full_qty_percent",
			ParentError: err,
		}
	}
	f.config.EP3ToFullQtyPercent, err = _type.ToFloat64(configMap["ep_3_to_full_qty_percent"])
	if err != nil || f.config.EP3ToFullQtyPercent == 0 {
		return tools.AppError{
			Message:     "Empty value ep_3_to_full_qty_percent",
			ParentError: err,
		}
	}

	f.config.HoursToReduceTP1, err = _type.ToInt64(configMap["hours_to_reduce_tp_1"])
	if err != nil {
		return tools.AppError{
			Message:     "The field hours_to_reduce_tp_1 is empty or contains not correct value type. Expects int64 value",
			ParentError: err,
		}
	}

	f.config.HoursToReduceTP2, err = _type.ToInt64(configMap["hours_to_reduce_tp_2"])
	if err != nil {
		return tools.AppError{
			Message:     "The field hours_to_reduce_tp_2 is empty or contains not correct value type. Expects int64 value",
			ParentError: err,
		}
	}

	f.config.HoursToReduceTP3, err = _type.ToInt64(configMap["hours_to_reduce_tp_3"])
	if err != nil {
		return tools.AppError{
			Message:     "The field hours_to_reduce_tp_3 is empty or contains not correct value type. Expects int64 value",
			ParentError: err,
		}
	}
	f.config.PercentToReduceTP, err = _type.ToInt64(configMap["percent_to_reduce_tp"])
	if err != nil {
		return tools.AppError{
			Message:     "The field percent_to_reduce_tp is empty or contains not correct value type. Expects int64 value",
			ParentError: err,
		}
	}

	f.config.WithdrawPercent, err = _type.ToFloat64(configMap["withdraw_percent"])
	if err != nil || f.config.WithdrawPercent == 0 {
		return tools.AppError{
			Message:     "Empty value withdraw_percent",
			ParentError: err,
		}
	}

	err = f.applyDomainConfig(configMap, domainConfig)
	if err != nil {
		return
	}

	return
}

func (f *Fibonacci) applyDomainConfig(configMap map[interface{}]interface{}, domainConfig map[string]interface{}) error {
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
	domainItem, err := domain.GetDomainInterface(domainCode)
	if err != nil {
		return err
	}

	err = domainItem.SetConfig(domainConfigItem)
	if err != nil {
		return err
	}
	f.provider = domainItem

	return nil
}
