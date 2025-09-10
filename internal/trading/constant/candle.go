package constant

import "github.com/shatylos/trader/tools"

const Resol1m = "1m"
const Resol5m = "5m"
const Resol15m = "15m"
const Resol30m = "30m"
const Resol1h = "1h"
const Resol2h = "2h"
const Resol4h = "4h"
const Resol1d = "D"
const Resol1w = "W"
const Resol1month = "M"

var resolToSec = map[string]int64{
	"1m":  60,
	"5m":  300,
	"15m": 900,
	"30m": 1800,
	"1h":  3600,
	"2h":  7200,
	"4h":  14400,
	"D":   86400,
	"W":   604800,
	"M":   2592000,
}

func ResolutionToSeconds(resolution string) (seconds int64, err error) {
	var exists bool
	seconds, exists = resolToSec[resolution]
	if !exists {
		err = tools.AppError{
			Message: "Unexpected resolution value",
		}
	}
	return
}
