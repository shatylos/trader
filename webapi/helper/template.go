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
		"dateFormat":     dateFormat,
		"longFloatShort": longFloatShort,
	}).ParseFiles(fileName)

	if err != nil {
		return nil, err
	}

	return tmpl, nil
}

func dateFormat(layout string, t time.Time) string {
	return t.Format(layout)
}

func longFloatShort(numsAfterDot int, value float64) string {
	valueStr := strconv.FormatFloat(value, 'f', numsAfterDot, 64)
	return valueStr
}
