package helper

import (
	"fmt"
	"github.com/shatylos/trader/tools/apperrors"
	"github.com/shatylos/trader/tools/math"
	"html/template"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func GetTemplate(fileName string) (*template.Template, error) {

	baseFileName := filepath.Base(fileName)

	tmpl, err := template.New(baseFileName).Funcs(template.FuncMap{
		"dict":               dict,
		"dateFormat":         dateFormat,
		"dateFormatFromUnix": dateFormatFromUnix,
		"longFloatShort":     longFloatShort,
		"moneyFormat":        moneyFormat,
		"add":                add,
		"mul":                mul,
		"isZeroTime":         isZeroTime,
	}).ParseFiles(
		fileName,
		"web/template/widget/baseMenu.html",
		"web/template/widget/menu.html",
	)

	if err != nil {
		err = apperrors.Wrap(err, "error parse template for fileName %s", fileName)
		return nil, err
	}

	return tmpl, nil
}

func dict(values ...any) (dict map[string]any) {
	dict = make(map[string]any)
	for i := 0; i < len(values); i += 2 {
		key := values[i].(string)
		dict[key] = values[i+1]
	}
	return dict
}

func dateFormat(layout string, t time.Time) string {
	return t.Format(layout)
}

func dateFormatFromUnix(layout string, timeUnix int64) string {
	t := time.Unix(timeUnix, 0)
	return t.Format(layout)
}

func longFloatShort(numsAfterDot int, value float64) string {
	valueStr := strconv.FormatFloat(value, 'f', numsAfterDot, 64)
	return valueStr
}

func add(x, y int) int {
	return x + y
}

func mul(x, y float64) float64 {
	return math.Mul(x, y)
}

func isZeroTime(t time.Time) bool {
	return t.IsZero()
}

func moneyFormat(numsAfterDot int, num float64) string {
	s := fmt.Sprintf("%.*f", numsAfterDot, num) // adjust precision as needed
	s = strings.TrimRight(s, "0")
	s = strings.TrimRight(s, ".")

	parts := strings.SplitN(s, ".", 2)
	intPart := parts[0]

	intPart = formatWithSeparator(intPart, ',', 3, false)

	if len(parts) == 1 {
		return intPart
	}

	fracPart := parts[1]
	fracPart = formatWithSeparator(fracPart, ',', 3, true)

	return intPart + "." + fracPart
}

func formatWithSeparator(s string, sep rune, group int, fromLeft bool) string {
	var out []rune
	runes := []rune(s)

	if fromLeft {
		for i, r := range runes {
			if i > 0 && i%group == 0 {
				out = append(out, sep)
			}
			out = append(out, r)
		}
	} else {
		count := 0
		for i := len(runes) - 1; i >= 0; i-- {
			if count > 0 && count%group == 0 && runes[i] != '-' {
				out = append([]rune{sep}, out...)
			}
			out = append([]rune{runes[i]}, out...)
			count++
		}
	}
	return string(out)
}
