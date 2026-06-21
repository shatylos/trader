package vwapReversion

import (
	"github.com/shatylos/trader/internal/domain"
	"github.com/shatylos/trader/internal/trading/constant"
	"github.com/shatylos/trader/tools/apperrors"
	"github.com/shatylos/trader/tools/logger"
	_type "github.com/shatylos/trader/tools/type"
	"time"
)

type Config struct {
	Enabled               bool
	Verbose               bool
	TelegramNotifier      bool
	Code                  string
	Id                    string
	CoinPare              string
	MainCurrency          string
	TradeCurrency         string
	MinDepo               float64
	Leverage              int64
	TimeoutSeconds        time.Duration
	Resolution            string
	ResolutionMins        int64
	VwapPeriod            int64
	LongTrendResolution   string
	LongTrendCandleReview int64
	PricePrecision        int64
	QtyPrecision          int64
	MinQty                float64
	EntrySigmaMult        float64
	SlSigmaMult           float64
	RequireConfirmation   bool
	RiskPercent           float64
	WithdrawPercent       float64
	IncreaseRiskToMinQty  bool
	LTCandlesCacheTTL     time.Duration
}

var resolutions = map[string]int64{
	constant.Resol1m:  1,
	constant.Resol5m:  5,
	constant.Resol15m: 15,
	constant.Resol30m: 30,
	constant.Resol1h:  60,
	constant.Resol2h:  120,
	constant.Resol4h:  240,
	constant.Resol1d:  86400,
	constant.Resol1w:  604800,
}

func (v *VwapReversion) SetConfig(strategyConfig interface{}, domainConfig map[string]interface{}) (err error) {

	v.config = Config{}

	configMap, ok := strategyConfig.(map[interface{}]interface{})
	if !ok {
		logger.Error("Config for the strategy vwap_reversion is not valid")
		return apperrors.New("config for the strategy vwap_reversion is not valid")
	}

	v.config.Id, err = _type.ToString(configMap["id"])
	if err != nil {
		return apperrors.Wrap(err, "the field id is empty or contains not correct value type. Expects string value")
	}

	var enabled int64
	enabled, err = _type.ToInt64(configMap["enabled"])
	if err != nil {
		return apperrors.Wrap(err, "the field enabled is empty or contains not correct value type. Expects 1 or 0")
	}
	v.config.Enabled = enabled == 1

	var verbose int64
	verbose, err = _type.ToInt64(configMap["verbose"])
	if err != nil {
		return apperrors.Wrap(err, "the field verbose is empty or contains not correct value type. Expects 1 or 0")
	}
	v.config.Verbose = verbose == 1

	if configMap["telegram_notifier"] != nil {
		var tgNotifier int64
		tgNotifier, err = _type.ToInt64(configMap["telegram_notifier"])
		if err != nil {
			return apperrors.Wrap(err, "the field telegram_notifier is empty or contains not correct value type. Expects 1 or 0")
		}
		v.config.TelegramNotifier = tgNotifier == 1
	}

	v.config.CoinPare, err = _type.ToString(configMap["coin_pare"])
	if err != nil {
		return apperrors.Wrap(err, "the field coin_pare is empty or contains not correct value type")
	}

	v.config.MainCurrency, err = _type.ToString(configMap["main_currency"])
	if err != nil {
		return apperrors.Wrap(err, "the field main_currency is empty or contains not correct value type")
	}

	v.config.TradeCurrency, err = _type.ToString(configMap["trade_currency"])
	if err != nil {
		return apperrors.Wrap(err, "the field trade_currency is empty or contains not correct value type")
	}

	v.config.MinDepo, err = _type.ToFloat64(configMap["min_depo"])
	if err != nil {
		return apperrors.Wrap(err, "wrong value min_depo")
	}

	v.config.Leverage, err = _type.ToInt64(configMap["leverage"])
	if err != nil || v.config.Leverage == 0 {
		return apperrors.Wrap(err, "wrong value leverage")
	}

	v.config.TimeoutSeconds, err = _type.ToTimeDuration(configMap["timeout_seconds"])
	if err != nil {
		return apperrors.Wrap(err, "the field timeout_seconds is empty or contains not correct value")
	}

	v.config.Resolution, err = _type.ToString(configMap["resolution"])
	if err != nil {
		return apperrors.Wrap(err, "the field resolution is empty or contains not correct value type")
	}
	v.config.ResolutionMins, ok = resolutions[v.config.Resolution]
	if !ok {
		return apperrors.New("the field resolution contains not correct value: %s", configMap["resolution"])
	}

	v.config.VwapPeriod, err = _type.ToInt64(configMap["vwap_period"])
	if err != nil || v.config.VwapPeriod < 10 {
		return apperrors.Wrap(err, "the field vwap_period is empty or contains not correct value type. Expects int64 value more than 10")
	}

	v.config.LongTrendResolution, err = _type.ToString(configMap["long_trend_resolution"])
	if err != nil {
		return apperrors.Wrap(err, "the field long_trend_resolution is empty or contains not correct value type")
	}

	v.config.LongTrendCandleReview, err = _type.ToInt64(configMap["long_trend_candle_review"])
	if err != nil || v.config.LongTrendCandleReview < 10 {
		return apperrors.Wrap(err, "the field long_trend_candle_review is empty or contains not correct value type. Expects int64 value more than 10")
	}

	v.config.PricePrecision, err = _type.ToInt64(configMap["price_precision"])
	if err != nil {
		return apperrors.Wrap(err, "the field price_precision is empty or contains not correct value type. Expects int64 value")
	}

	v.config.QtyPrecision, err = _type.ToInt64(configMap["qty_precision"])
	if err != nil {
		return apperrors.Wrap(err, "the field qty_precision is empty or contains not correct value type. Expects int64 value")
	}

	v.config.MinQty, err = _type.ToFloat64(configMap["min_qty"])
	if err != nil {
		return apperrors.Wrap(err, "the field min_qty is empty or contains not correct value type. Expects float64 value")
	}

	v.config.EntrySigmaMult, err = _type.ToFloat64(configMap["entry_sigma_mult"])
	if err != nil || v.config.EntrySigmaMult <= 0 {
		return apperrors.Wrap(err, "the field entry_sigma_mult is empty or contains not correct value type. Expects float64 value more than 0")
	}

	v.config.SlSigmaMult, err = _type.ToFloat64(configMap["sl_sigma_mult"])
	if err != nil || v.config.SlSigmaMult <= v.config.EntrySigmaMult {
		return apperrors.Wrap(err, "the field sl_sigma_mult is empty or must be greater than entry_sigma_mult")
	}

	var requireConfirmation int64
	requireConfirmation, err = _type.ToInt64(configMap["require_confirmation"])
	if err != nil {
		return apperrors.Wrap(err, "the field require_confirmation is empty or contains not correct value type. Expects 1 or 0")
	}
	v.config.RequireConfirmation = requireConfirmation == 1

	v.config.RiskPercent, err = _type.ToFloat64(configMap["risk_percent"])
	if err != nil || v.config.RiskPercent == 0 {
		return apperrors.Wrap(err, "empty value risk_percent")
	}

	v.config.WithdrawPercent, err = _type.ToFloat64(configMap["withdraw_percent"])
	if err != nil {
		return apperrors.Wrap(err, "the field withdraw_percent is empty")
	}

	if configMap["increase_resk_to_min_qty"] != nil {
		var increaseRiskToMinQty int64
		increaseRiskToMinQty, err = _type.ToInt64(configMap["increase_resk_to_min_qty"])
		if err != nil {
			return apperrors.Wrap(err, "the field increase_resk_to_min_qty is empty or contains not correct value type. Expects 1 or 0")
		}
		v.config.IncreaseRiskToMinQty = increaseRiskToMinQty == 1
	}

	v.config.LTCandlesCacheTTL, err = _type.ToTimeDuration(configMap["long_trend_candles_cache_ttl"])
	if err != nil {
		return apperrors.Wrap(err, "the field long_trend_candles_cache_ttl is empty or contains not correct value")
	}
	v.config.LTCandlesCacheTTL = v.config.LTCandlesCacheTTL * time.Second

	err = v.applyDomainConfig(configMap, domainConfig)
	if err != nil {
		return apperrors.Wrap(err, "error apply domain config")
	}

	return
}

func (v *VwapReversion) applyDomainConfig(configMap map[interface{}]interface{}, domainConfig map[string]interface{}) error {
	domainId, err := _type.ToString(configMap["domain"])
	if err != nil {
		return apperrors.Wrap(err, "the field domain is empty or contains not correct value type")
	}
	domainConfigItem, ok := domainConfig[domainId].(map[interface{}]interface{})
	if !ok {
		return apperrors.New("the domain config is not valid. Domain value should be related to the domain config item")
	}
	domainCode, err := _type.ToString(domainConfigItem["code"])
	if err != nil {
		return apperrors.Wrap(err, "the field code is empty in domain config or contains not correct value type")
	}
	domainItem, err := domain.GetFuturesDomain(domainCode)
	if err != nil {
		return apperrors.Wrap(err, "error get futures domain")
	}

	err = domainItem.SetConfig(domainConfigItem)
	if err != nil {
		return apperrors.Wrap(err, "error set domain config")
	}
	v.provider = domainItem

	return nil
}
