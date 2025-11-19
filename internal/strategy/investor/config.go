package investor

import (
	"github.com/shatylos/trader/internal/domain"
	"github.com/shatylos/trader/internal/strategy/investor/entity"
	"github.com/shatylos/trader/internal/strategy/investor/storage"
	_struct "github.com/shatylos/trader/internal/strategy/investor/struct"
	"github.com/shatylos/trader/tools"
	"github.com/shatylos/trader/tools/logger"
	_type "github.com/shatylos/trader/tools/type"
	"time"
)

type Config struct {
	Id                    string
	Enabled               bool
	Verbose               bool
	TelegramNotifier      bool
	CoinPare              string
	MainCurrency          string
	TradeCurrency         string
	TimeoutDuration       time.Duration
	QtyPrecision          int64
	PricePrecision        int64
	CommissionBuy         float64
	CommissionSell        float64
	MinCoinReservePercent float64
	MinQty                float64
	DoIncreaseQtyToMinQty bool
	WithdrawPercent       float64
	RequestDelay          time.Duration
}

func (i *Investor) SetConfig(strategyConfig interface{}, domainConfig map[string]interface{}) (err error) {

	i.Config = Config{}

	configMap, ok := strategyConfig.(map[interface{}]interface{})
	if !ok {
		logger.Error("Config for the strategy investor is not valid")
		err = tools.AppError{
			Message: "Config for the strategy investor is not valid",
		}
		return
	}

	i.Config.Id, err = _type.ToString(configMap["id"])
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
		i.Config.Enabled = true
	} else {
		i.Config.Enabled = false
	}

	var verbose int64
	verbose, err = _type.ToInt64(configMap["verbose"])
	if err != nil {
		logger.Error("The field verbose in strategy Config is empty or contains not correct value type. Expects 1 or 0")
		err = tools.AppError{
			Message:     "The field verbose in strategy Config is empty or contains not correct value type. Expects 1 or 0",
			ParentError: err,
		}
		return
	}
	if verbose == 1 {
		i.Config.Verbose = true
	} else {
		i.Config.Verbose = false
	}

	if configMap["telegram_notifier"] != nil {
		var tgNotifier int64
		tgNotifier, err = _type.ToInt64(configMap["telegram_notifier"])
		if err != nil {
			logger.Error("The field telegram_notifier in strategy Config is empty or contains not correct value type. Expects 1 or 0")
			err = tools.AppError{
				Message:     "The field telegram_notifier in strategy Config is empty or contains not correct value type. Expects 1 or 0",
				ParentError: err,
			}
			return
		}
		if tgNotifier == 1 {
			i.Config.TelegramNotifier = true
		}
	}

	i.Config.CoinPare, err = _type.ToString(configMap["coin_pare"])
	if err != nil {
		err = tools.AppError{
			Message:     "The field coin_pare is empty or contains not correct value type",
			ParentError: err,
		}
		return
	}

	i.Config.MainCurrency, err = _type.ToString(configMap["main_currency"])
	if err != nil {
		return tools.AppError{
			Message:     "The field main_currency is empty or contains not correct value type",
			ParentError: err,
		}
	}

	i.Config.TradeCurrency, err = _type.ToString(configMap["trade_currency"])
	if err != nil {
		return tools.AppError{
			Message:     "The field trade_currency is empty or contains not correct value type",
			ParentError: err,
		}
	}

	i.Config.TimeoutDuration, err = _type.ToTimeDuration(configMap["timeout_seconds"])
	if err != nil {
		err = tools.AppError{
			Message:     "The field timeout_seconds is empty or contains not correct value",
			ParentError: err,
		}
		return
	}
	i.Config.TimeoutDuration = i.Config.TimeoutDuration * time.Second

	i.Config.QtyPrecision, err = _type.ToInt64(configMap["qty_precision"])
	if err != nil {
		return tools.AppError{
			Message:     "The field qty_precision is empty or contains not correct value type. Expects int64 value",
			ParentError: err,
		}
	}

	i.Config.PricePrecision, err = _type.ToInt64(configMap["price_precision"])
	if err != nil {
		return tools.AppError{
			Message:     "The field price_precision is empty or contains not correct value type. Expects int64 value",
			ParentError: err,
		}
	}

	i.Config.CommissionBuy, err = _type.ToFloat64(configMap["commission_buy"])
	if err != nil {
		return tools.AppError{
			Message:     "The field commission_buy is empty",
			ParentError: err,
		}
	}

	i.Config.CommissionSell, err = _type.ToFloat64(configMap["commission_sell"])
	if err != nil {
		return tools.AppError{
			Message:     "The field commission_sell is empty",
			ParentError: err,
		}
	}

	i.Config.MinCoinReservePercent, err = _type.ToFloat64(configMap["min_coin_reserve_percent"])
	if err != nil {
		return tools.AppError{
			Message:     "The field min_coin_reserve_percent is empty",
			ParentError: err,
		}
	}

	i.Config.MinQty, err = _type.ToFloat64(configMap["min_qty"])
	if err != nil {
		return tools.AppError{
			Message:     "Empty value min_qty",
			ParentError: err,
		}
	}

	var doIncreaseQtyToMinQty int64
	doIncreaseQtyToMinQty, _ = _type.ToInt64(configMap["do_increase_qty_to_min_qty"])
	if doIncreaseQtyToMinQty == 1 {
		i.Config.DoIncreaseQtyToMinQty = true
	}

	i.Config.WithdrawPercent, err = _type.ToFloat64(configMap["withdraw_percent"])
	if err != nil {
		return tools.AppError{
			Message:     "The field withdraw_percent is empty",
			ParentError: err,
		}
	}

	i.Config.RequestDelay, err = _type.ToTimeDuration(configMap["request_delay"])
	if err != nil {
		return tools.AppError{
			Message:     "The field request_delay is empty",
			ParentError: err,
		}
	}
	i.Config.RequestDelay = i.Config.RequestDelay * time.Second

	err = i.applyTimeframeConfig(configMap["timeframes"])
	if err != nil {
		return
	}

	err = i.applyHeapConfig(configMap["heap_timeframe"])
	if err != nil {
		return
	}

	err = i.applyDomainConfig(configMap, domainConfig)
	if err != nil {
		return
	}

	i.Storage = storage.Storage{Id: i.Config.Id}
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

	for _, tf := range rawTimeframes {
		tfMap, ok := tf.(map[interface{}]interface{})
		if !ok {
			logger.Error("Config for the strategy Investor is not valid. Item of \"timeframes\" must be map ok key value")
			err = tools.AppError{
				Message: "Config for the strategy Investor is not valid. Item of \"timeframes\" must be map ok key value",
			}
			return
		}

		timeframe := _struct.TimeframeItem{
			Config: _struct.TimeframeItemConfig{},
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

		timeframe.Config.CandleCacheDuration, err = _type.ToTimeDuration(tfMap["candle_cache_seconds"])
		if err != nil {
			return tools.AppError{
				Message:     "The field candle_cache_seconds is empty or contains not correct value type. Expects int64 value more than 1",
				ParentError: err,
			}
		}
		timeframe.Config.CandleCacheDuration = timeframe.Config.CandleCacheDuration * time.Second
		if timeframe.Config.CandleCacheDuration < time.Second {
			return tools.AppError{
				Message:     "The field candle_cache_seconds contains not correct value. Expects int64 value more than 1",
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

		timeframe.Config.MaxNumberOrdersToBuy, err = _type.ToInt64(tfMap["max_number_orders_to_buy"])
		if err != nil || timeframe.Config.MaxNumberOrdersToBuy == 0 {
			return tools.AppError{
				Message:     "Empty value max_number_orders_to_buy",
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

		timeframe.Config.MinPercentRangeToSell, err = _type.ToFloat64(tfMap["min_percent_range_to_sell"])
		if err != nil {
			return tools.AppError{
				Message:     "Empty value min_percent_range_to_sell",
				ParentError: err,
			}
		}

		timeframe.Config.MinPercentRangeToBuyMore, err = _type.ToFloat64(tfMap["min_percent_range_to_buy_more"])
		if err != nil {
			return tools.AppError{
				Message:     "Empty value min_percent_range_to_buy_more",
				ParentError: err,
			}
		}

		timeframe.Config.DurationToMoveToHeap, err = _type.ToTimeDuration(tfMap["hours_duration_to_move_to_heap"])
		if err != nil {
			return tools.AppError{
				Message:     "The field hours_duration_to_move_to_heap is empty",
				ParentError: err,
			}
		}
		timeframe.Config.DurationToMoveToHeap = timeframe.Config.DurationToMoveToHeap * time.Hour

		minHigherTFSlope, exists := tfMap["min_higher_tf_slope"]
		if exists {
			timeframe.Config.IsCheckHigherTF = true
			timeframe.Config.MinHigherTFSlope, err = _type.ToFloat64(minHigherTFSlope)
			if err != nil {
				return tools.AppError{
					Message:     "The field min_higher_tf_slope contains not correct value. Expected float64",
					ParentError: err,
				}
			}
		} else {
			timeframe.Config.IsCheckHigherTF = false
		}

		higherTFResolution, exists := tfMap["higher_tf_resolution"]
		if exists {
			timeframe.Config.HigherTFResolution, err = _type.ToString(higherTFResolution)
			if err != nil {
				return tools.AppError{
					Message:     "The field higher_tf_resolution is empty or contains not correct value type",
					ParentError: err,
				}
			}
		}

		i.Timeframes = append(i.Timeframes, timeframe)
	}

	return
}

func (i *Investor) applyHeapConfig(heapTimeframesConfig interface{}) (err error) {
	tfMap, ok := heapTimeframesConfig.(map[interface{}]interface{})
	if !ok {
		logger.Error("Config for the strategy Investor is not valid. Item of \"timeframes\" must be map ok key value")
		err = tools.AppError{
			Message: "Config for the strategy Investor is not valid. Item of \"timeframes\" must be map ok key value",
		}
		return
	}

	heapTimeframe := _struct.HeapTimeframe{
		Config:     _struct.HeapConfig{},
		HeapStatus: entity.HeapStatus{},
	}

	heapTimeframe.Config.Resolution, err = _type.ToString(tfMap["resolution"])
	if err != nil {
		return tools.AppError{
			Message:     "The field resolution is empty or contains not correct value type",
			ParentError: err,
		}
	}

	heapTimeframe.Config.QtyPercent, err = _type.ToFloat64(tfMap["qty_percent"])
	if err != nil || heapTimeframe.Config.QtyPercent == 0 {
		return tools.AppError{
			Message:     "Empty value qty_percent",
			ParentError: err,
		}
	}

	heapTimeframe.Config.CandleReview, err = _type.ToInt64(tfMap["candle_review"])
	if err != nil || heapTimeframe.Config.CandleReview < 10 {
		return tools.AppError{
			Message:     "The field candle_review is empty or contains not correct value type. Expects int64 value more than 10",
			ParentError: err,
		}
	}

	heapTimeframe.Config.CandleCacheDuration, err = _type.ToTimeDuration(tfMap["candle_cache_seconds"])
	if err != nil {
		return tools.AppError{
			Message:     "The field candle_cache_seconds is empty or contains not correct value type. Expects int64 value more than 1",
			ParentError: err,
		}
	}
	heapTimeframe.Config.CandleCacheDuration = heapTimeframe.Config.CandleCacheDuration * time.Second
	if heapTimeframe.Config.CandleCacheDuration < time.Second {
		return tools.AppError{
			Message:     "The field candle_cache_seconds contains not correct value. Expects int64 value more than 1",
			ParentError: err,
		}
	}

	heapTimeframe.Config.SidewaysMinCandlesAmount, err = _type.ToInt64(tfMap["sideways_min_candles_amount"])
	if err != nil || heapTimeframe.Config.SidewaysMinCandlesAmount == 0 {
		return tools.AppError{
			Message:     "Empty value sideways_min_candles_amount",
			ParentError: err,
		}
	}

	heapTimeframe.Config.MaxNumberOrdersToBuy, err = _type.ToInt64(tfMap["max_number_orders_to_buy"])
	if err != nil || heapTimeframe.Config.MaxNumberOrdersToBuy == 0 {
		return tools.AppError{
			Message:     "Empty value max_number_orders_to_buy",
			ParentError: err,
		}
	}

	heapTimeframe.Config.SidewaysPercentToPrice, err = _type.ToFloat64(tfMap["sideways_percent_to_price"])
	if err != nil || heapTimeframe.Config.SidewaysPercentToPrice == 0 {
		return tools.AppError{
			Message:     "Empty value sideways_percent_to_price",
			ParentError: err,
		}
	}

	heapTimeframe.Config.SidewaysPremiumCoefficient, err = _type.ToFloat64(tfMap["sideways_premium_coefficient"])
	if err != nil {
		return tools.AppError{
			Message:     "Empty value sideways_premium_coefficient",
			ParentError: err,
		}
	}

	heapTimeframe.Config.SidewaysDiscountCoefficient, err = _type.ToFloat64(tfMap["sideways_discount_coefficient"])
	if err != nil {
		return tools.AppError{
			Message:     "Empty value sideways_discount_coefficient",
			ParentError: err,
		}
	}

	heapTimeframe.Config.MinPercentRangeToSell, err = _type.ToFloat64(tfMap["min_percent_range_to_sell"])
	if err != nil {
		return tools.AppError{
			Message:     "Empty value min_percent_range_to_sell",
			ParentError: err,
		}
	}

	heapTimeframe.Config.MinPercentRangeToBuyMore, err = _type.ToFloat64(tfMap["min_percent_range_to_buy_more"])
	if err != nil {
		return tools.AppError{
			Message:     "Empty value min_percent_range_to_buy_more",
			ParentError: err,
		}
	}

	heapTimeframe.Config.QtyPercentOnMaxPrice, err = _type.ToFloat64(tfMap["qty_percent_on_max_price"])
	if err != nil {
		return tools.AppError{
			Message:     "Empty value qty_percent_on_max_price for heap config",
			ParentError: err,
		}
	}
	heapTimeframe.Config.QtyPercentOnMinPrice, err = _type.ToFloat64(tfMap["qty_percent_on_min_price"])
	if err != nil {
		return tools.AppError{
			Message:     "Empty value qty_percent_on_min_price for heap config",
			ParentError: err,
		}
	}

	i.HeapTimeframe = heapTimeframe
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
			Message: "The domain Config is not valid. Domain value should be related to the domain Config item",
		}
	}
	domainCode, err := _type.ToString(domainConfigItem["code"])
	if err != nil {
		return tools.AppError{
			Message:     "The field code is empty in domain Config or contains not correct value type",
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
