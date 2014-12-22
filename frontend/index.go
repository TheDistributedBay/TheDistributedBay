package frontend

import (
	"fmt"
	"html/template"
	"net/http"
)

func IndexHandler(resp http.ResponseWriter, req *http.Request) {
	t, err := template.ParseFiles("frontend/templates/index.html")
	if err != nil {
		fmt.Println(err)
	}
	t.Execute(resp, 0)
}
