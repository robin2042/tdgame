package telegram

import (
	"bytes"
	"fmt"
	"text/template"
)

// main.go
type UserInfo struct {
	Name   string
	Gender string
	Age    int
}

func TemplateNiuniu_Bet() string {
	var b bytes.Buffer
	tmpl, err := template.ParseFiles("./templates/niuniu/bet.tmpl")
	if err != nil {
		fmt.Println("create template failed,err:", err)
		return "无效"
	}

	tmpl.Execute(&b, nil)
	return b.String()
}
