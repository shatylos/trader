package mapping

import (
	"fmt"
	"github.com/shatylos/trader/internal/trading/constant"
	"github.com/shatylos/trader/tools"
	"github.com/shatylos/trader/tools/logger"
)

var bybitIntervals = map[string]string{
	constant.Resol1m:     "1",
	constant.Resol5m:     "5",
	constant.Resol15m:    "15",
	constant.Resol30m:    "30",
	constant.Resol1h:     "60",
	constant.Resol2h:     "120",
	constant.Resol4h:     "240",
	constant.Resol1d:     "D",
	constant.Resol1w:     "W",
	constant.Resol1month: "M",
}

func ToBybitInterval(domainInterval string) (bybitInterval string, err error) {
	var ok bool
	bybitInterval, ok = bybitIntervals[domainInterval]
	if ok == false {
		msg := fmt.Sprintf("Unexpected value (%s) to map Bybit interval", domainInterval)
		logger.Error(msg)
		err = tools.AppError{
			Message: msg,
		}
		return
	}
	return
}
