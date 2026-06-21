package vwapReversion

import (
	"fmt"
	domainStructs "github.com/shatylos/trader/internal/domain/structs"
	"github.com/shatylos/trader/tools/apperrors"
	"time"
)

func (v *VwapReversion) getAvailableBalance() (balance float64, err error) {
	var wallet domainStructs.DomainWallet
	wallet, err = v.provider.GetWallet()
	if err != nil {
		err = apperrors.Wrap(err, "error get wallet")
		return
	}
	for _, coin := range wallet.Available {
		if coin.Coin == v.config.MainCurrency {
			balance = coin.Amount
		}
	}
	return
}

type candleCacheEntry struct {
	candles  []domainStructs.DomainCandle
	loadedAt time.Time
}

func (v *VwapReversion) LoadLTCandleHistory() ([]domainStructs.DomainCandle, error) {
	return v.loadCachedCandleHistory(v.config.LongTrendResolution, v.config.LongTrendCandleReview, v.config.LTCandlesCacheTTL)
}

func (v *VwapReversion) LoadSTCandleHistory() ([]domainStructs.DomainCandle, error) {
	return v.loadCachedCandleHistory(v.config.Resolution, v.config.VwapPeriod, time.Duration(0))
}

func (v *VwapReversion) loadCachedCandleHistory(resolution string, limit int64, cacheTTL time.Duration) ([]domainStructs.DomainCandle, error) {
	cacheKey := fmt.Sprintf("%s:%d", resolution, limit)

	if v.candleCache == nil {
		v.candleCache = make(map[string]candleCacheEntry)
	}

	entry, ok := v.candleCache[cacheKey]
	if ok && time.Since(entry.loadedAt) < cacheTTL {
		return entry.candles, nil
	}

	candles, err := v.provider.LoadCandleHistory(v.config.CoinPare, resolution, limit)
	if err != nil {
		return nil, apperrors.Wrap(err, "error load candle history")
	}

	v.candleCache[cacheKey] = candleCacheEntry{
		candles:  candles,
		loadedAt: time.Now(),
	}

	return candles, nil
}
