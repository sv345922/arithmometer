package handler

import (
	"arithmometer/orchestr/parsing"
	"log"
	"net/http"
)

func GetResult(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	_, err := parsing.Parse(r.Form.Get("id"))
	if err != nil {
		log.Println(err)
		return
	}

	// прочитать json выражения
	// если посчитано, вернуть ответ
	// иначе не посчитано ответить "ожидайте"
}
