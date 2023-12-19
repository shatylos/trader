package buyCheapSellHigh

import (
	"bitbucket.org/shatylos/trader/domain"
	"bitbucket.org/shatylos/trader/utils"
	"errors"
)

func (s *BuyCheapSellHigh) SetConfig(strategyConfig interface{}, domainConfig map[interface{}]interface{}) error {

	var err error

	configMap, ok := strategyConfig.(map[interface{}]interface{})
	if !ok {
		return utils.AppError{
			Message: "Config for the strategy buyCheapSellHigh is not valid",
		}
	}

	s.Id, err = utils.ToString(configMap["id"])
	if err != nil {
		return utils.AppError{
			Message:     "The field id is empty or contains not correct value type. Expects string value",
			ParentError: err,
		}
	}

	s.CoinPare, err = utils.ToString(configMap["coin_pare"])
	if err != nil {
		return utils.AppError{
			Message:     "The field coin_pare is empty or contains not correct value type",
			ParentError: err,
		}
	}

	s.MainCurrency, err = utils.ToString(configMap["main_currency"])
	if err != nil {
		return utils.AppError{
			Message:     "The field main_currency is empty or contains not correct value type",
			ParentError: err,
		}
	}

	s.TradeCurrency, err = utils.ToString(configMap["trade_currency"])
	if err != nil {
		return utils.AppError{
			Message:     "The field trade_currency is empty or contains not correct value type",
			ParentError: err,
		}
	}

	s.TimeoutSeconds, err = utils.ToTimeDuration(configMap["timeout_seconds"])
	if err != nil {
		return utils.AppError{
			Message:     "The field timeout_seconds is empty or contains not correct value",
			ParentError: err,
		}
	}

	s.Resolution, err = utils.ToString(configMap["resolution"])
	if err != nil {
		return utils.AppError{
			Message:     "The field resolution is empty or contains not correct value type",
			ParentError: err,
		}
	}

	s.LongTermMaxPrice, err = utils.ToFloat64(configMap["long_term_max_price"])
	if err != nil {
		return utils.AppError{
			Message:     "The field long_term_max_price is empty or contains not correct value type",
			ParentError: err,
		}
	}

	s.LongTermMinPrice, err = utils.ToFloat64(configMap["long_term_min_price"])
	if err != nil {
		return utils.AppError{
			Message:     "The field long_term_min_price is empty or contains not correct value type",
			ParentError: err,
		}
	}

	s.LongTermPercentBuffer, err = utils.ToFloat64(configMap["long_term_percent_buffer"])
	if err != nil {
		return utils.AppError{
			Message:     "The field long_term_percent_buffer is empty or contains not correct value type",
			ParentError: err,
		}
	}

	s.MainCurrencyPrecision, err = utils.ToInt64(configMap["main_currency_precision"])
	if err != nil {
		return utils.AppError{
			Message:     "The field main_currency_precision is empty or contains not correct value type. Expects int64 value",
			ParentError: err,
		}
	}

	s.PurchaseVolumePrecision, err = utils.ToInt64(configMap["purchase_volume_precision"])
	if err != nil {
		return utils.AppError{
			Message:     "The field purchase_volume_precision is empty or contains not correct value type. Expects int64 value",
			ParentError: err,
		}
	}

	s.PurchasePricePrecision, err = utils.ToInt64(configMap["purchase_price_precision"])
	if err != nil {
		return utils.AppError{
			Message:     "The field purchase_price_precision is empty or contains not correct value type. Expects int64 value",
			ParentError: err,
		}
	}

	s.MinutesToReducePriceRange, err = utils.ToInt64(configMap["minutes_to_reduce_price_range"])
	if err != nil {
		return utils.AppError{
			Message:     "The field minutes_to_reduce_price_range is empty or contains not correct value type. Expects int64 value",
			ParentError: err,
		}
	}

	s.AvgPriceCandleLimit, err = utils.ToInt64(configMap["avg_price_candle_limit"])
	if err != nil || s.AvgPriceCandleLimit < 1 {
		return utils.AppError{
			Message:     "The field avg_price_candle_limit is empty or contains not correct value type. Expects int64 value more than 1",
			ParentError: err,
		}
	}

	s.AvgPriceCandleOffset, err = utils.ToInt64(configMap["avg_price_candle_offset"])
	if err != nil || s.AvgPriceCandleOffset < 1 {
		return utils.AppError{
			Message:     "The field avg_price_candle_offset is empty or contains not correct value type. Expects int64 value more than 1",
			ParentError: err,
		}
	}

	s.ManualBuyPriceBeforeStart, err = utils.ToFloat64(configMap["manual_buy_price_before_start"])
	if err != nil && !errors.Is(err, utils.EmptyValueError) {
		return utils.AppError{
			Message:     "The field manual_buy_price_before_start contains not correct value type. Expects float64 value. Can be empty",
			ParentError: err,
		}
	}

	s.CostRanges, err = utils.ToInt64Slice(configMap["cost_ranges"])
	if err != nil {
		return utils.AppError{
			Message:     "The field cost_ranges is empty or contains not correct value type",
			ParentError: err,
		}
	}

	s.PercentRanges, err = utils.ToInt64Slice(configMap["percent_ranges"])
	if err != nil {
		return utils.AppError{
			Message:     "The field percent_ranges is empty or contains not correct value type",
			ParentError: err,
		}
	}

	err = s.applyDomainConfig(configMap, domainConfig)
	if err != nil {
		return err
	}

	return nil
}

func (s *BuyCheapSellHigh) applyDomainConfig(configMap map[interface{}]interface{}, domainConfig map[interface{}]interface{}) error {
	domainId, err := utils.ToString(configMap["domain"])
	if err != nil {
		return utils.AppError{
			Message:     "The field domain is empty or contains not correct value type",
			ParentError: err,
		}
	}
	domainConfigItem, ok := domainConfig[domainId].(map[interface{}]interface{})
	if !ok {
		return utils.AppError{
			Message: "The domain config is not valid. Domain value should be related to the domain config item",
		}
	}
	domainCode, err := utils.ToString(domainConfigItem["code"])
	if err != nil {
		return utils.AppError{
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
	s.Domain = domainItem

	return nil
}
