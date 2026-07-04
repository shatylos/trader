package scalper

import (
	"github.com/shatylos/trader/internal/domain"
	"github.com/shatylos/trader/internal/trading/constant"
	"github.com/shatylos/trader/tools/apperrors"
	"github.com/shatylos/trader/tools/logger"
	_type "github.com/shatylos/trader/tools/type"
	"time"
)

type Config struct {
	Enabled              bool
	Verbose              bool
	TelegramNotifier     bool
	Code                 string
	Id                   string
	CoinPare             string
	MainCurrency         string
	TradeCurrency        string
	Leverage             int64
	TimeoutSeconds       time.Duration
	PricePrecision       int64
	QtyPrecision         int64
	MinQty               float64
	RiskPercent          float64
	WithdrawPercent      float64
	IncreaseRiskToMinQty bool

	// entry (lower) timeframe
	EntryResolution   string
	EntryCandleReview int64
	LtfEmaPeriod      int64
	RsiPeriod         int64
	RsiMomentumLevel  float64
	RsiOverbought     float64
	RsiOversold       float64
	AtrPeriod         int64
	MinAtrPercent     float64
	PullbackLookback  int64
	TpAtrMult         float64
	SlAtrMult         float64

	// bias (higher) timeframe
	BiasResolution           string
	BiasCandleReview         int64
	BiasCandlesCacheTTL      time.Duration
	HtfEmaFastPeriod         int64
	HtfEmaSlowPeriod         int64
	RequireTrendConfirmation bool
}

var allowedResolutions = map[string]bool{
	constant.Resol1m:  true,
	constant.Resol5m:  true,
	constant.Resol15m: true,
	constant.Resol30m: true,
	constant.Resol1h:  true,
	constant.Resol2h:  true,
	constant.Resol4h:  true,
	constant.Resol1d:  true,
	constant.Resol1w:  true,
}

func (s *Scalper) SetConfig(strategyConfig interface{}, domainConfig map[string]interface{}) (err error) {

	s.config = Config{}

	configMap, ok := strategyConfig.(map[interface{}]interface{})
	if !ok {
		logger.Error("Config for the strategy scalper is not valid")
		return apperrors.New("config for the strategy scalper is not valid")
	}

	s.config.Id, err = _type.ToString(configMap["id"])
	if err != nil {
		return apperrors.Wrap(err, "the field id is empty or contains not correct value type. Expects string value")
	}

	var enabled int64
	enabled, err = _type.ToInt64(configMap["enabled"])
	if err != nil {
		return apperrors.Wrap(err, "the field enabled is empty or contains not correct value type. Expects 1 or 0")
	}
	s.config.Enabled = enabled == 1

	var verbose int64
	verbose, err = _type.ToInt64(configMap["verbose"])
	if err != nil {
		return apperrors.Wrap(err, "the field verbose is empty or contains not correct value type. Expects 1 or 0")
	}
	s.config.Verbose = verbose == 1

	if configMap["telegram_notifier"] != nil {
		var tgNotifier int64
		tgNotifier, err = _type.ToInt64(configMap["telegram_notifier"])
		if err != nil {
			return apperrors.Wrap(err, "the field telegram_notifier is empty or contains not correct value type. Expects 1 or 0")
		}
		s.config.TelegramNotifier = tgNotifier == 1
	}

	s.config.CoinPare, err = _type.ToString(configMap["coin_pare"])
	if err != nil {
		return apperrors.Wrap(err, "the field coin_pare is empty or contains not correct value type")
	}

	s.config.MainCurrency, err = _type.ToString(configMap["main_currency"])
	if err != nil {
		return apperrors.Wrap(err, "the field main_currency is empty or contains not correct value type")
	}

	s.config.TradeCurrency, err = _type.ToString(configMap["trade_currency"])
	if err != nil {
		return apperrors.Wrap(err, "the field trade_currency is empty or contains not correct value type")
	}

	s.config.Leverage, err = _type.ToInt64(configMap["leverage"])
	if err != nil || s.config.Leverage == 0 {
		return apperrors.Wrap(err, "wrong value leverage")
	}

	s.config.TimeoutSeconds, err = _type.ToTimeDuration(configMap["timeout_seconds"])
	if err != nil {
		return apperrors.Wrap(err, "the field timeout_seconds is empty or contains not correct value")
	}

	s.config.PricePrecision, err = _type.ToInt64(configMap["price_precision"])
	if err != nil {
		return apperrors.Wrap(err, "the field price_precision is empty or contains not correct value type. Expects int64 value")
	}

	s.config.QtyPrecision, err = _type.ToInt64(configMap["qty_precision"])
	if err != nil {
		return apperrors.Wrap(err, "the field qty_precision is empty or contains not correct value type. Expects int64 value")
	}

	s.config.MinQty, err = _type.ToFloat64(configMap["min_qty"])
	if err != nil {
		return apperrors.Wrap(err, "the field min_qty is empty or contains not correct value type. Expects float64 value")
	}

	s.config.RiskPercent, err = _type.ToFloat64(configMap["risk_percent"])
	if err != nil || s.config.RiskPercent == 0 {
		return apperrors.Wrap(err, "empty value risk_percent")
	}

	s.config.WithdrawPercent, err = _type.ToFloat64(configMap["withdraw_percent"])
	if err != nil {
		return apperrors.Wrap(err, "the field withdraw_percent is empty")
	}

	if configMap["increase_risk_to_min_qty"] != nil {
		var increaseRiskToMinQty int64
		increaseRiskToMinQty, err = _type.ToInt64(configMap["increase_risk_to_min_qty"])
		if err != nil {
			return apperrors.Wrap(err, "the field increase_risk_to_min_qty is empty or contains not correct value type. Expects 1 or 0")
		}
		s.config.IncreaseRiskToMinQty = increaseRiskToMinQty == 1
	}

	err = s.setEntryTimeframeConfig(configMap)
	if err != nil {
		return apperrors.Wrap(err, "error set entry timeframe config")
	}

	err = s.setBiasTimeframeConfig(configMap)
	if err != nil {
		return apperrors.Wrap(err, "error set bias timeframe config")
	}

	err = s.applyDomainConfig(configMap, domainConfig)
	if err != nil {
		return apperrors.Wrap(err, "error apply domain config")
	}

	return
}

func (s *Scalper) setEntryTimeframeConfig(configMap map[interface{}]interface{}) (err error) {
	s.config.EntryResolution, err = _type.ToString(configMap["entry_resolution"])
	if err != nil {
		return apperrors.Wrap(err, "the field entry_resolution is empty or contains not correct value type")
	}
	if !allowedResolutions[s.config.EntryResolution] {
		return apperrors.New("the field entry_resolution contains not correct value: %s", s.config.EntryResolution)
	}

	s.config.LtfEmaPeriod, err = _type.ToInt64(configMap["ltf_ema_period"])
	if err != nil || s.config.LtfEmaPeriod < 2 {
		return apperrors.Wrap(err, "the field ltf_ema_period is empty or contains not correct value type. Expects int64 value more than 1")
	}

	s.config.RsiPeriod, err = _type.ToInt64(configMap["rsi_period"])
	if err != nil || s.config.RsiPeriod < 2 {
		return apperrors.Wrap(err, "the field rsi_period is empty or contains not correct value type. Expects int64 value more than 1")
	}

	s.config.RsiMomentumLevel, err = _type.ToFloat64(configMap["rsi_momentum_level"])
	if err != nil || s.config.RsiMomentumLevel <= 0 || s.config.RsiMomentumLevel >= 100 {
		return apperrors.Wrap(err, "the field rsi_momentum_level is empty or must be between 0 and 100")
	}

	s.config.RsiOverbought, err = _type.ToFloat64(configMap["rsi_overbought"])
	if err != nil || s.config.RsiOverbought <= s.config.RsiMomentumLevel {
		return apperrors.Wrap(err, "the field rsi_overbought is empty or must be greater than rsi_momentum_level")
	}

	s.config.RsiOversold, err = _type.ToFloat64(configMap["rsi_oversold"])
	if err != nil || s.config.RsiOversold >= s.config.RsiMomentumLevel {
		return apperrors.Wrap(err, "the field rsi_oversold is empty or must be less than rsi_momentum_level")
	}

	s.config.AtrPeriod, err = _type.ToInt64(configMap["atr_period"])
	if err != nil || s.config.AtrPeriod < 2 {
		return apperrors.Wrap(err, "the field atr_period is empty or contains not correct value type. Expects int64 value more than 1")
	}

	s.config.MinAtrPercent, err = _type.ToFloat64(configMap["min_atr_percent"])
	if err != nil || s.config.MinAtrPercent < 0 {
		return apperrors.Wrap(err, "the field min_atr_percent is empty or contains not correct value type. Expects float64 value 0 or more")
	}

	s.config.PullbackLookback, err = _type.ToInt64(configMap["pullback_lookback"])
	if err != nil || s.config.PullbackLookback < 1 {
		return apperrors.Wrap(err, "the field pullback_lookback is empty or contains not correct value type. Expects int64 value more than 0")
	}

	s.config.SlAtrMult, err = _type.ToFloat64(configMap["sl_atr_mult"])
	if err != nil || s.config.SlAtrMult <= 0 {
		return apperrors.Wrap(err, "the field sl_atr_mult is empty or contains not correct value type. Expects float64 value more than 0")
	}

	s.config.TpAtrMult, err = _type.ToFloat64(configMap["tp_atr_mult"])
	if err != nil || s.config.TpAtrMult <= 0 {
		return apperrors.Wrap(err, "the field tp_atr_mult is empty or contains not correct value type. Expects float64 value more than 0")
	}

	s.config.EntryCandleReview, err = _type.ToInt64(configMap["entry_candle_review"])
	if err != nil {
		return apperrors.Wrap(err, "the field entry_candle_review is empty or contains not correct value type. Expects int64 value")
	}
	// The indicators are calculated on closed candles only (the forming candle is
	// excluded) and the pullback scan needs valid indicator values over the whole
	// lookback window, so demand enough history for the slowest one plus a buffer
	// to let the Wilder smoothing settle.
	minReview := maxInt64(s.config.LtfEmaPeriod, maxInt64(s.config.RsiPeriod+1, s.config.AtrPeriod+1)) + s.config.PullbackLookback + 20
	if s.config.EntryCandleReview < minReview {
		return apperrors.New("the field entry_candle_review must be at least %d for the configured indicator periods", minReview)
	}

	return
}

func (s *Scalper) setBiasTimeframeConfig(configMap map[interface{}]interface{}) (err error) {
	s.config.BiasResolution, err = _type.ToString(configMap["bias_resolution"])
	if err != nil {
		return apperrors.Wrap(err, "the field bias_resolution is empty or contains not correct value type")
	}
	if !allowedResolutions[s.config.BiasResolution] {
		return apperrors.New("the field bias_resolution contains not correct value: %s", s.config.BiasResolution)
	}

	s.config.HtfEmaFastPeriod, err = _type.ToInt64(configMap["htf_ema_fast_period"])
	if err != nil || s.config.HtfEmaFastPeriod < 2 {
		return apperrors.Wrap(err, "the field htf_ema_fast_period is empty or contains not correct value type. Expects int64 value more than 1")
	}

	s.config.HtfEmaSlowPeriod, err = _type.ToInt64(configMap["htf_ema_slow_period"])
	if err != nil || s.config.HtfEmaSlowPeriod <= s.config.HtfEmaFastPeriod {
		return apperrors.Wrap(err, "the field htf_ema_slow_period is empty or must be greater than htf_ema_fast_period")
	}

	s.config.BiasCandleReview, err = _type.ToInt64(configMap["bias_candle_review"])
	if err != nil {
		return apperrors.Wrap(err, "the field bias_candle_review is empty or contains not correct value type. Expects int64 value")
	}
	minReview := s.config.HtfEmaSlowPeriod + 20
	if s.config.BiasCandleReview < minReview {
		return apperrors.New("the field bias_candle_review must be at least %d for the configured indicator periods", minReview)
	}

	s.config.BiasCandlesCacheTTL, err = _type.ToTimeDuration(configMap["bias_candles_cache_ttl"])
	if err != nil {
		return apperrors.Wrap(err, "the field bias_candles_cache_ttl is empty or contains not correct value")
	}
	s.config.BiasCandlesCacheTTL = s.config.BiasCandlesCacheTTL * time.Second

	var requireTrendConfirmation int64
	requireTrendConfirmation, err = _type.ToInt64(configMap["require_trend_confirmation"])
	if err != nil {
		return apperrors.Wrap(err, "the field require_trend_confirmation is empty or contains not correct value type. Expects 1 or 0")
	}
	s.config.RequireTrendConfirmation = requireTrendConfirmation == 1

	return
}

func (s *Scalper) applyDomainConfig(configMap map[interface{}]interface{}, domainConfig map[string]interface{}) error {
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
	s.provider = domainItem

	return nil
}

func maxInt64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
