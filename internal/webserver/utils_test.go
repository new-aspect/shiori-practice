package webserver

import (
	"html/template"
	"os"
	"testing"
)

func TestTemplate(t *testing.T) {
	var srcHTML = `
		<!DOCTYPE html>
		<html>
		<head>
			<title>{{ . }}</title>
		</head>
		<body>
			<h1>{{ .  }}</h1>
		</body>
		</html>
	`

	temp, err := template.New("test").Parse(srcHTML)
	if err != nil {
		panic(err)
	}

	// 这个验证说明是空的啊，他不会给你一个模板的默认值
	err = temp.Execute(os.Stdout, "/")
	if err != nil {
		panic(err)
	}
}
