package fibonacci

import (
	"fmt"
	domainStructs "github.com/shatylos/trader/internal/domain/structs"
	"github.com/shatylos/trader/tools"
	"time"
)

func (f *Fibonacci) getCurrentPrice() (currentPrice float64, err error) {
	var candles []domainStructs.DomainCandle
	candles, err = f.provider.LoadCandleHistory(f.config.CoinPare, f.config.Resolution, time.Now().Unix(), 1)
	if err != nil {
		return
	}
	if len(candles) != 1 {
		err = tools.AppError{Message: fmt.Sprintf("Unexpected length of candles (%d items). Expected 1 item", len(candles))}
	}
	currentPrice = candles[0].Close
	return
}
