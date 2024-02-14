package handler

import (
	"arithmometer/orchestr/inter/tasker"
	"arithmometer/orchestr/parsing"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func NewExpression(w http.ResponseWriter, r *http.Request) {
	// Получаем список выражений из контекста
	ws, ok := tasker.GetWs(r.Context())
	if !ok {
		log.Println("ошибка контекста")
		return
	}

	expressions := ws.Expressions

	// Проверить что это запрос POST
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("требуется метод POST"))
		return
	}
	// Читаем тело запроса, в котором записано выражение и тайминги операций
	var newExrp tasker.NewExpr
	err := json.NewDecoder(r.Body).Decode(&newExrp)
	if err != nil {
		log.Println("ошибка POST запроса")
		return
	}
	//Если тайминги не передаются, тогда они ставятся по умолчанию
	if newExrp.Timings == nil {
		newExrp.Timings = &tasker.Timings{
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

	// Создаем и сохраняем задачу
	expression := tasker.Expression{
		UserTask: newExrp.Expr,
		Postfix:  postfix,
		Times:    *newExrp.Timings,
	}
	expression.CreateId()

	expressions.Add(&expression)

	// err = SafeJSON("new_expression", expression)

	// Записываем тело ответа
	body := fmt.Sprintf("Expression:(id): %s\nTimings: %s", expression.Id, expression.Times.String())
	w.Write([]byte(body))

}
