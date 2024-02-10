package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"arithmometer/orchestr/parsing"
)

func NewExpression(w http.ResponseWriter, r *http.Request) {
	// Проверить что это запрос POST
	if r.Method != http.MethodGet { //http.MethodPost {
		log.Println(r.Method)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid method"))
		return
	}
	err := r.ParseForm()
	// Парсим выражение, и проверяем его
	// предполагается, что если парсинг с ошибкой, значит невалидное выражение
	expr := r.Form.Get("expr")
	log.Println(expr)
	postfix, _, _, err := parsing.Parse(expr)
	// если невалидное выражение
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid expression"))
		return
	}
	// Читаем тело запроса, в котором записаны тайминги операций
	var timings Timings
	err = json.NewDecoder(r.Body).Decode(&timings)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Set default timings\n"))
		// тайминги по умолчанию
		timings = Timings{1, 1, 1, 1}
	}

	// Сохраняем задачу
	// TODO

	expression := Expression{
		Postfix: postfix,
		Times:   timings,
	}
	expression.doId()
	err = SafeJSON("new_expression", expression)
	w.Write([]byte(expression.Id))
}

// getNodes Возвращает список узлов

// SafeJSON Сохраняет структуру в базе данных, в папке db
func SafeJSON[T additiveJSON](name string, expr T) error {
	jsonBytes, err := json.Marshal(expr)
	if err != nil {
		log.Println(err)
		return err
	}
	wd, err := os.Getwd()
	if err != nil {
		log.Println(err)
		return err
	}
	path := wd + "/orchestr/db/" + name + ".json"
	err = os.WriteFile(path, jsonBytes, 0666)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
func LoadJSON[T additiveJSON](name string) {
	// TODO
}
