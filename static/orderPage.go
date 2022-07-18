package static

import (
	"L0/internal"
	"fmt"
	"reflect"
)

const header = "<!DOCTYPE html>\n" +
	"<html lang=\"en\">\n" +
	"<head>\n" +
	"<meta charset=\"UTF-8\">\n" +
	"<title>WB L0</title>\n" +
	"<style>\n" +
	".zebra {\n" +
	"list-style: none;\n" +
	"border-left: 10px solid #FC7574;\n" +
	"padding: 0;\n" +
	"font-family: \"Lucida Sans\";\n" +
	"}\n" +
	".zebra li {padding: 10px;}\n" +
	".zebra li:nth-child(odd) {background: #E1F1FF;}\n" +
	".zebra li:nth-child(even) {background: white;}\n" +
	"</style>\n" +
	"</head>\n" +
	"<body>\n" + "\t<form action=\"\" method=\"get\">\n" +
	"ID: <input type=\"text\" name=\"uid\">\n" +
	"<input type=\"submit\" value=\"uid\">\n" +
	"</form>"

func generateList(data interface{}) string {
	val := reflect.ValueOf(data).Elem()
	htmlString := "<ul class=\"zebra\">\n"

	for i := 0; i < val.NumField(); i++ {
		htmlString = fmt.Sprint(htmlString+"<li><span>", val.Type().Field(i).Name,
			": </span><em>", val.Field(i).Interface(), "</em></li>\n")
	}
	htmlString += "</ul>\n"
	return htmlString
}

func GeneratePage(order *internal.Order) string {
	return header + generateList(order) + "</body>\n</html>"
}

func GenerateNotFound() string {
	return header + "<h1> Order not found </h1>\n" + "</body>\n</html>"
}
