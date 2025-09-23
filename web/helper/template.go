package helper

import (
	"html/template"
	"path/filepath"
	"strconv"
	"time"
)

func GetTemplate(fileName string) (*template.Template, error) {

	baseFileName := filepath.Base(fileName)

	tmpl, err := template.New(baseFileName).Funcs(template.FuncMap{
		"dateFormat":         dateFormat,
		"dateFormatFromUnix": dateFormatFromUnix,
		"longFloatShort":     longFloatShort,
		"add":                add,
	}).ParseFiles(
		fileName,
		"web/template/widget/baseMenu.html",
		"web/template/widget/menu.html",
	)

	if err != nil {
		return nil, err
	}

	return tmpl, nil
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
