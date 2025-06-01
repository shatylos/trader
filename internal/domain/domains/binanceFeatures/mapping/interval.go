package mapping

import (
	"fmt"
	"github.com/shatylos/trader/internal/trading/constant"
	"github.com/shatylos/trader/tools"
	"github.com/shatylos/trader/tools/logger"
)

var binanceIntervals = map[string]string{
	constant.Resol1m:     "1m",
	constant.Resol5m:     "5m",
	constant.Resol15m:    "15m",
	constant.Resol30m:    "30m",
	constant.Resol1h:     "1h",
	constant.Resol2h:     "2h",
	constant.Resol4h:     "4h",
	constant.Resol1d:     "1d",
	constant.Resol1w:     "1w",
	constant.Resol1month: "1M",
}

func ToBinanceInterval(domainInterval string) (binanceInterval string, err error) {
	var ok bool
	binanceInterval, ok = binanceIntervals[domainInterval]
	if ok == false {
		msg := fmt.Sprintf("Unexpected value (%s) to map Binance interval", domainInterval)
		logger.Error(msg)
		err = tools.AppError{
			Message: msg,
		}
		return
	}
	return
}
