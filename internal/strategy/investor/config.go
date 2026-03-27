package investor

import (
	"errors"
	"github.com/shatylos/trader/internal/domain"
	"github.com/shatylos/trader/internal/strategy/investor/entity"
	"github.com/shatylos/trader/internal/strategy/investor/storage"
	_struct "github.com/shatylos/trader/internal/strategy/investor/struct"
	"github.com/shatylos/trader/tools/apperrors"
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
		err = apperrors.New("config for the strategy investor is not valid")
		return
	}

	i.Config.Id, err = _type.ToString(configMap["id"])
	if err != nil {
		err = apperrors.Wrap(err, "the field \"id\" is empty or contains not correct value type. Expects string value")
		return
	}

	var enabled int64
	enabled, err = _type.ToInt64(configMap["enabled"])
	if err != nil {
		err = apperrors.Wrap(err, "the field \"enabled\" is empty or contains not correct value type. Expects 1 or 0")
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
		err = apperrors.Wrap(err, "the field \"verbose\" in strategy Config is empty or contains not correct value type. Expects 1 or 0")
		return
	}
	if verbose == 1 {
		i.Config.Verbose = true
	} else {
		i.Config.Verbose = false
	}

	var tgNotifier int64
	tgNotifier, err = _type.ToInt64(configMap["telegram_notifier"])
	if err != nil && !errors.Is(err, _type.EmptyValueError) {
		err = apperrors.Wrap(err, "the field telegram_notifier in strategy Config contains not correct value type. Expects 1 or 0")
		return
	}
	if tgNotifier == 1 {
		i.Config.TelegramNotifier = true
	}

	i.Config.CoinPare, err = _type.ToString(configMap["coin_pare"])
	if err != nil {
		err = apperrors.Wrap(err, "the field coin_pare is empty or contains not correct value type")
		return
	}

	i.Config.MainCurrency, err = _type.ToString(configMap["main_currency"])
	if err != nil {
		err = apperrors.Wrap(err, "the field main_currency is empty or contains not correct value type")
		return
	}

	i.Config.TradeCurrency, err = _type.ToString(configMap["trade_currency"])
	if err != nil {
		err = apperrors.Wrap(err, "the field trade_currency is empty or contains not correct value type")
		return
	}

	i.Config.TimeoutDuration, err = _type.ToTimeDuration(configMap["timeout_seconds"])
	if err != nil {
		err = apperrors.Wrap(err, "the field timeout_seconds is empty or contains not correct value")
		return
	}
	i.Config.TimeoutDuration = i.Config.TimeoutDuration * time.Second

	i.Config.QtyPrecision, err = _type.ToInt64(configMap["qty_precision"])
	if err != nil {
		err = apperrors.Wrap(err, "the field qty_precision is empty or contains not correct value type. Expects int64 value")
		return
	}

	i.Config.PricePrecision, err = _type.ToInt64(configMap["price_precision"])
	if err != nil {
		err = apperrors.Wrap(err, "the field price_precision is empty or contains not correct value type. Expects int64 value")
		return
	}

	i.Config.CommissionBuy, err = _type.ToFloat64(configMap["commission_buy"])
	if err != nil {
		err = apperrors.Wrap(err, "the field commission_buy is empty")
		return
	}

	i.Config.CommissionSell, err = _type.ToFloat64(configMap["commission_sell"])
	if err != nil {
		err = apperrors.Wrap(err, "the field commission_sell is empty")
		return
	}

	i.Config.MinCoinReservePercent, err = _type.ToFloat64(configMap["min_coin_reserve_percent"])
	if err != nil {
		err = apperrors.Wrap(err, "the field min_coin_reserve_percent is empty")
		return
	}

	i.Config.MinQty, err = _type.ToFloat64(configMap["min_qty"])
	if err != nil {
		err = apperrors.Wrap(err, "empty value min_qty")
		return
	}

	var doIncreaseQtyToMinQty int64
	doIncreaseQtyToMinQty, err = _type.ToInt64(configMap["do_increase_qty_to_min_qty"])
	if err != nil && !errors.Is(err, _type.EmptyValueError) {
		err = apperrors.Wrap(err, "the field do_increase_qty_to_min_qty in strategy Config contains not correct value type. Expects 1 or 0")
		return
	}
	if doIncreaseQtyToMinQty == 1 {
		i.Config.DoIncreaseQtyToMinQty = true
	}

	i.Config.WithdrawPercent, err = _type.ToFloat64(configMap["withdraw_percent"])
	if err != nil {
		err = apperrors.Wrap(err, "the field withdraw_percent is empty")
		return
	}

	i.Config.RequestDelay, err = _type.ToTimeDuration(configMap["request_delay"])
	if err != nil {
		err = apperrors.Wrap(err, "the field request_delay is empty")
		return
	}
	i.Config.RequestDelay = i.Config.RequestDelay * time.Second

	err = i.applyTimeframeConfig(configMap["timeframes"])
	if err != nil {
		err = apperrors.Wrap(err, "error apply timeframes config")
		return
	}

	err = i.applyHeapConfig(configMap["heap_timeframe"])
	if err != nil {
		err = apperrors.Wrap(err, "error apply heap timeframe config")
		return
	}

	err = i.applyDomainConfig(configMap, domainConfig)
	if err != nil {
		err = apperrors.Wrap(err, "error apply domain config")
		return
	}

	i.Storage = storage.Storage{Id: i.Config.Id}
	err = i.Storage.InitStorage()
	if err != nil {
		err = apperrors.Wrap(err, "error init storage")
		return
	}

	return
}

func (i *Investor) applyTimeframeConfig(timeframesConfig interface{}) (err error) {

	rawTimeframes, ok := timeframesConfig.([]interface{})
	if !ok {
		err = apperrors.New("config for the strategy Investor is not valid. \"timeframes\" must be array")
		return
	}

	for _, tf := range rawTimeframes {
		tfMap, ok := tf.(map[interface{}]interface{})
		if !ok {
			err = apperrors.New("config for the strategy Investor is not valid. Item of \"timeframes\" must be map ok key value")
			return
		}

		timeframe := _struct.TimeframeItem{
			Config: _struct.TimeframeItemConfig{},
		}

		timeframe.Config.Resolution, err = _type.ToString(tfMap["resolution"])
		if err != nil {
			err = apperrors.Wrap(err, "the field resolution is empty or contains not correct value type")
			return
		}

		var canOpenNewOrder int64
		canOpenNewOrder, err = _type.ToInt64(tfMap["can_open_new_order"])
		if err != nil {
			err = apperrors.Wrap(err, "the field \"can_open_new_order\" is empty or contains not correct value type. Expects 1 or 0")
			return
		}
		if canOpenNewOrder == 1 {
			timeframe.Config.CanOpenNewOrder = true
		} else {
			timeframe.Config.CanOpenNewOrder = false
		}

		timeframe.Config.QtyPercent, err = _type.ToFloat64(tfMap["qty_percent"])
		if err != nil || timeframe.Config.QtyPercent == 0 {
			err = apperrors.Wrap(err, "empty value qty_percent")
			return
		}

		timeframe.Config.CandleReview, err = _type.ToInt64(tfMap["candle_review"])
		if err != nil || timeframe.Config.CandleReview < 10 {
			err = apperrors.Wrap(err, "the field candle_review is empty or contains not correct value type. Expects int64 value more than 10")
			return
		}

		timeframe.Config.CandleCacheDuration, err = _type.ToTimeDuration(tfMap["candle_cache_seconds"])
		if err != nil {
			err = apperrors.Wrap(err, "the field candle_cache_seconds is empty or contains not correct value type. Expects int64 value more than 1")
			return
		}
		timeframe.Config.CandleCacheDuration = timeframe.Config.CandleCacheDuration * time.Second
		if timeframe.Config.CandleCacheDuration < time.Second {
			err = apperrors.Wrap(err, "the field candle_cache_seconds contains not correct value. Expects int64 value more than 1")
			return
		}

		timeframe.Config.SidewaysMinCandlesAmount, err = _type.ToInt64(tfMap["sideways_min_candles_amount"])
		if err != nil || timeframe.Config.SidewaysMinCandlesAmount == 0 {
			err = apperrors.Wrap(err, "empty value sideways_min_candles_amount")
			return
		}

		timeframe.Config.MaxNumberOrdersToBuy, err = _type.ToInt64(tfMap["max_number_orders_to_buy"])
		if err != nil || timeframe.Config.MaxNumberOrdersToBuy == 0 {
			err = apperrors.Wrap(err, "empty value max_number_orders_to_buy")
			return
		}

		timeframe.Config.SidewaysPercentToPrice, err = _type.ToFloat64(tfMap["sideways_percent_to_price"])
		if err != nil || timeframe.Config.SidewaysPercentToPrice == 0 {
			err = apperrors.Wrap(err, "empty value sideways_percent_to_price")
			return
		}

		timeframe.Config.SidewaysPremiumCoefficient, err = _type.ToFloat64(tfMap["sideways_premium_coefficient"])
		if err != nil {
			err = apperrors.Wrap(err, "empty value sideways_premium_coefficient")
			return
		}

		timeframe.Config.SidewaysDiscountCoefficient, err = _type.ToFloat64(tfMap["sideways_discount_coefficient"])
		if err != nil {
			err = apperrors.Wrap(err, "empty value sideways_discount_coefficient")
			return
		}

		timeframe.Config.MinPercentRangeToSell, err = _type.ToFloat64(tfMap["min_percent_range_to_sell"])
		if err != nil {
			err = apperrors.Wrap(err, "empty value min_percent_range_to_sell")
			return
		}

		timeframe.Config.MinPercentRangeToBuyMore, err = _type.ToFloat64(tfMap["min_percent_range_to_buy_more"])
		if err != nil {
			err = apperrors.Wrap(err, "empty value min_percent_range_to_buy_more")
			return
		}

		timeframe.Config.DurationToMoveToHeap, err = _type.ToTimeDuration(tfMap["hours_duration_to_move_to_heap"])
		if err != nil {
			err = apperrors.Wrap(err, "the field hours_duration_to_move_to_heap is empty")
			return
		}
		timeframe.Config.DurationToMoveToHeap = timeframe.Config.DurationToMoveToHeap * time.Hour

		minHigherTFSlope, exists := tfMap["min_higher_tf_slope"]
		if exists {
			timeframe.Config.IsCheckHigherTF = true
			timeframe.Config.MinHigherTFSlope, err = _type.ToFloat64(minHigherTFSlope)
			if err != nil {
				err = apperrors.Wrap(err, "the field min_higher_tf_slope contains not correct value. Expected float64")
				return
			}
		} else {
			timeframe.Config.IsCheckHigherTF = false
		}

		higherTFResolution, exists := tfMap["higher_tf_resolution"]
		if exists {
			timeframe.Config.HigherTFResolution, err = _type.ToString(higherTFResolution)
			if err != nil {
				err = apperrors.Wrap(err, "the field higher_tf_resolution is empty or contains not correct value type")
				return
			}
		}

		i.Timeframes = append(i.Timeframes, timeframe)
	}

	return
}

func (i *Investor) applyHeapConfig(heapTimeframesConfig interface{}) (err error) {
	tfMap, ok := heapTimeframesConfig.(map[interface{}]interface{})
	if !ok {
		err = apperrors.New("config for the strategy Investor is not valid. Item of \"timeframes\" must be map ok key value")
		return
	}

	heapTimeframe := _struct.HeapTimeframe{
		Config:     _struct.HeapConfig{},
		HeapStatus: entity.HeapStatus{},
	}

	heapTimeframe.Config.Resolution, err = _type.ToString(tfMap["resolution"])
	if err != nil {
		err = apperrors.Wrap(err, "the field resolution is empty or contains not correct value type")
		return
	}

	heapTimeframe.Config.QtyPercent, err = _type.ToFloat64(tfMap["qty_percent"])
	if err != nil || heapTimeframe.Config.QtyPercent == 0 {
		err = apperrors.Wrap(err, "empty value qty_percent")
		return
	}

	heapTimeframe.Config.CandleReview, err = _type.ToInt64(tfMap["candle_review"])
	if err != nil || heapTimeframe.Config.CandleReview < 10 {
		err = apperrors.Wrap(err, "the field candle_review is empty or contains not correct value type. Expects int64 value more than 10")
		return
	}

	heapTimeframe.Config.CandleCacheDuration, err = _type.ToTimeDuration(tfMap["candle_cache_seconds"])
	if err != nil {
		err = apperrors.Wrap(err, "the field candle_cache_seconds is empty or contains not correct value type. Expects int64 value more than 1")
		return
	}
	heapTimeframe.Config.CandleCacheDuration = heapTimeframe.Config.CandleCacheDuration * time.Second
	if heapTimeframe.Config.CandleCacheDuration < time.Second {
		err = apperrors.Wrap(err, "the field candle_cache_seconds contains not correct value. Expects int64 value more than 1")
		return
	}

	heapTimeframe.Config.SidewaysMinCandlesAmount, err = _type.ToInt64(tfMap["sideways_min_candles_amount"])
	if err != nil || heapTimeframe.Config.SidewaysMinCandlesAmount == 0 {
		err = apperrors.Wrap(err, "empty value sideways_min_candles_amount")
		return
	}

	heapTimeframe.Config.MaxNumberOrdersToBuy, err = _type.ToInt64(tfMap["max_number_orders_to_buy"])
	if err != nil || heapTimeframe.Config.MaxNumberOrdersToBuy == 0 {
		err = apperrors.Wrap(err, "empty value max_number_orders_to_buy")
		return
	}

	heapTimeframe.Config.SidewaysPercentToPrice, err = _type.ToFloat64(tfMap["sideways_percent_to_price"])
	if err != nil || heapTimeframe.Config.SidewaysPercentToPrice == 0 {
		err = apperrors.Wrap(err, "empty value sideways_percent_to_price")
		return
	}

	heapTimeframe.Config.SidewaysPremiumCoefficient, err = _type.ToFloat64(tfMap["sideways_premium_coefficient"])
	if err != nil {
		err = apperrors.Wrap(err, "empty value sideways_premium_coefficient")
		return
	}

	heapTimeframe.Config.SidewaysDiscountCoefficient, err = _type.ToFloat64(tfMap["sideways_discount_coefficient"])
	if err != nil {
		err = apperrors.Wrap(err, "empty value sideways_discount_coefficient")
		return
	}

	heapTimeframe.Config.MinPercentRangeToSell, err = _type.ToFloat64(tfMap["min_percent_range_to_sell"])
	if err != nil {
		err = apperrors.Wrap(err, "empty value min_percent_range_to_sell")
		return
	}

	heapTimeframe.Config.MinPercentRangeToBuyMore, err = _type.ToFloat64(tfMap["min_percent_range_to_buy_more"])
	if err != nil {
		err = apperrors.Wrap(err, "empty value min_percent_range_to_buy_more")
		return
	}

	heapTimeframe.Config.QtyPercentOnMaxPrice, err = _type.ToFloat64(tfMap["qty_percent_on_max_price"])
	if err != nil {
		err = apperrors.Wrap(err, "empty value qty_percent_on_max_price for heap config")
		return
	}
	heapTimeframe.Config.QtyPercentOnMinPrice, err = _type.ToFloat64(tfMap["qty_percent_on_min_price"])
	if err != nil {
		err = apperrors.Wrap(err, "empty value qty_percent_on_min_price for heap config")
		return
	}

	i.HeapTimeframe = heapTimeframe
	return
}

func (i *Investor) applyDomainConfig(configMap map[interface{}]interface{}, domainConfig map[string]interface{}) (err error) {
	var domainId string
	domainId, err = _type.ToString(configMap["domain"])
	if err != nil {
		err = apperrors.Wrap(err, "the field domain is empty or contains not correct value type")
		return
	}
	domainConfigItem, ok := domainConfig[domainId].(map[interface{}]interface{})
	if !ok {
		err = apperrors.New("the domain Config is not valid. Domain value should be related to the domain Config item")
		return
	}
	var domainCode string
	domainCode, err = _type.ToString(domainConfigItem["code"])
	if err != nil {
		err = apperrors.Wrap(err, "the field code is empty in domain Config or contains not correct value type")
		return
	}
	var domainItem domain.SpotDomainInterface
	domainItem, err = domain.GetSpotDomain(domainCode)
	if err != nil {
		err = apperrors.Wrap(err, "error get spot domain with code %s", domainCode)
		return err
	}

	err = domainItem.SetConfig(domainConfigItem)
	if err != nil {
		err = apperrors.Wrap(err, "error set config with code %s", domainCode)
		return err
	}
	i.provider = domainItem

	return nil
}
