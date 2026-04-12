package investor

import (
	"errors"
	"github.com/shatylos/trader/internal/domain"
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

		timeframe.Config.FullAmountPercent, err = _type.ToFloat64(tfMap["full_amount_percent"])
		if err != nil || timeframe.Config.FullAmountPercent == 0 {
			err = apperrors.Wrap(err, "empty value full_amount_percent")
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

		timeframe.Config.SidewaysPercentToPrice, err = _type.ToFloat64(tfMap["sideways_percent_to_price"])
		if err != nil || timeframe.Config.SidewaysPercentToPrice == 0 {
			err = apperrors.Wrap(err, "empty value sideways_percent_to_price")
			return
		}

		var equalAllOrders int64
		equalAllOrders, err = _type.ToInt64(tfMap["equal_all_orders"])
		if err != nil {
			err = apperrors.Wrap(err, "the field \"equal_all_orders\" is empty or contains not correct value type. Expects 1 or 0")
			return
		}
		if equalAllOrders == 1 {
			timeframe.Config.IsEqualAllOrders = true
		}

		rawVwapDeviationsBuy, ok := tfMap["vwap_deviations_buy"].([]interface{})
		if !ok {
			err = apperrors.New("config is not valid. \"vwap_deviations_buy\" must be array")
			return
		}
		timeframe.Config.VwapDeviationsBuy = make([]float64, len(rawVwapDeviationsBuy))
		for i, val := range rawVwapDeviationsBuy {
			var deviation float64
			deviation, err = _type.ToFloat64(val)
			if err != nil {
				err = apperrors.Wrap(err, "config is not valid. \"vwap_deviations_buy\" must be array of float64. Given value: %s", val)
				return
			}
			timeframe.Config.VwapDeviationsBuy[i] = deviation
		}

		rawVwapDeviationsSell, ok := tfMap["vwap_deviations_sell"].([]interface{})
		if !ok {
			err = apperrors.New("config is not valid. \"vwap_deviations_sell\" must be array")
			return
		}
		timeframe.Config.VwapDeviationsSell = make([]float64, len(rawVwapDeviationsSell))
		for i, val := range rawVwapDeviationsSell {
			var deviation float64
			deviation, err = _type.ToFloat64(val)
			if err != nil {
				err = apperrors.Wrap(err, "config is not valid. \"vwap_deviations_sell\" must be array of float64. Given value: %s", val)
				return
			}
			timeframe.Config.VwapDeviationsSell[i] = deviation
		}

		i.Timeframes = append(i.Timeframes, timeframe)
	}

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
