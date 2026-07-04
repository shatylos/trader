package scalper

import (
	"fmt"
	domainStructs "github.com/shatylos/trader/internal/domain/structs"
	"github.com/shatylos/trader/tools/apperrors"
	"time"
)

func (s *Scalper) getAvailableBalance() (balance float64, err error) {
	var wallet domainStructs.DomainWallet
	wallet, err = s.provider.GetWallet()
	if err != nil {
		err = apperrors.Wrap(err, "error get wallet")
		return
	}
	for _, coin := range wallet.Available {
		if coin.Coin == s.config.MainCurrency {
			balance = coin.Amount
		}
	}
	return
}

type candleCacheEntry struct {
	candles  []domainStructs.DomainCandle
	loadedAt time.Time
}

func (s *Scalper) LoadBiasCandleHistory() ([]domainStructs.DomainCandle, error) {
	return s.loadCachedCandleHistory(s.config.BiasResolution, s.config.BiasCandleReview, s.config.BiasCandlesCacheTTL)
}

func (s *Scalper) LoadEntryCandleHistory() ([]domainStructs.DomainCandle, error) {
	return s.loadCachedCandleHistory(s.config.EntryResolution, s.config.EntryCandleReview, time.Duration(0))
}

func (s *Scalper) loadCachedCandleHistory(resolution string, limit int64, cacheTTL time.Duration) ([]domainStructs.DomainCandle, error) {
	cacheKey := fmt.Sprintf("%s:%d", resolution, limit)

	if s.candleCache == nil {
		s.candleCache = make(map[string]candleCacheEntry)
	}

	entry, ok := s.candleCache[cacheKey]
	if ok && time.Since(entry.loadedAt) < cacheTTL {
		return entry.candles, nil
	}

	candles, err := s.provider.LoadCandleHistory(s.config.CoinPare, resolution, limit)
	if err != nil {
		return nil, apperrors.Wrap(err, "error load candle history")
	}

	s.candleCache[cacheKey] = candleCacheEntry{
		candles:  candles,
		loadedAt: time.Now(),
	}

	return candles, nil
}
