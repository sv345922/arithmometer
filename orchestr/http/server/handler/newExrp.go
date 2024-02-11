package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"arithmometer/orchestr/parsing"
)

func NewExpression(w http.ResponseWriter, r *http.Request) {
	// Проверить что это запрос POST
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("требуется метод POST"))
		return
	}
	// Читаем тело запроса, в котором записано выражение и тайминги операций
	var newExrp NewExpr
	err := json.NewDecoder(r.Body).Decode(&newExrp)
	if err != nil {
		log.Println("ошибка POST запроса")
		return
	}
	//Если тайминги не передаются, тогда они ставятся по умолчанию
	if newExrp.Timings == nil {
		newExrp.Timings = &Timings{
			Plus:  1,
			Minus: 1,
			Mult:  1,
			Div:   1,
		}
	}
	log.Printf("Method: %s, Expression: %s, Timings: %s", r.Method, newExrp.Expr, newExrp.Timings.String())
	// Парсим выражение, и проверяем его
	// предполагается, что если парсинг с ошибкой, значит невалидное выражение
	postfix, err := parsing.Parse(newExrp.Expr)
	// если невалидное выражение
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid expression"))
		return
	}

	// Сохраняем задачу
	expression := Expression{
		Postfix: postfix,
		Times:   *newExrp.Timings,
	}
	expression.doId()
	err = SafeJSON("new_expression", expression)
	body := fmt.Sprintf("Expression:(id): %s\nTimings: %s", expression.Id, expression.Times.String())
	w.Write([]byte(body))
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
func LoadJSON[T additiveJSON](name string) (*T, error) {
	var result T
	wd, err := os.Getwd()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	path := wd + "/orchestr/db/" + name + ".json"
	data, err := os.ReadFile(path)
	err = json.Unmarshal(data, result)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &result, nil
}
