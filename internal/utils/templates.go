package utils

import (
	"html/template"
)

func OrderTemplate() *template.Template {
	t, err := template.ParseFiles("internal/utils/orderTemplate.html")
	if err != nil {
		Logger.Error(err.Error())
	}
	return t
}
