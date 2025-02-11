package constant

const Resol1m = "1"
const Resol5m = "5"
const Resol15m = "15"
const Resol30m = "30"
const Resol45m = "45"
const Resol1h = "60"
const Resol2h = "120"
const Resol3h = "180"
const Resol4h = "240"
const Resol1d = "D"
const Resol1w = "W"
const Resol1month = "M"

var ResolToSec = map[string]int64{
	"1":   60,
	"5":   300,
	"15":  900,
	"30":  1800,
	"45":  2700,
	"60":  3600,
	"120": 7200,
	"180": 10800,
	"240": 14400,
	"D":   86400,
	"W":   604800,
	"M":   2592000,
}
