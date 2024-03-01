package handler

import (
	"arithmometer/orchestr/inter/tasker"
	"arithmometer/orchestr/parsing"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Обработчик создания нового выражения
func NewExpression(ws *tasker.WorkingSpace) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Проверить что это запрос POST
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("требуется метод POST"))
			return
		}
		// Читаем тело запроса, в котором записано выражение и тайминги операций
		var newExrp tasker.NewExpr
		err := json.NewDecoder(r.Body).Decode(&newExrp)
		defer r.Body.Close()
		if err != nil {
			log.Println("ошибка POST запроса")
			return
		}
		//Если тайминги не передаются, тогда они ставятся по умолчанию
		if newExrp.Timings == nil {
			newExrp.Timings = &tasker.Timings{Plus: 1, Minus: 1, Mult: 1, Div: 1}
		}
		// Парсим выражение, и проверяем его
		// Предполагается, что если парсинг с ошибкой, значит невалидное выражение
		postfix, err := parsing.Parse(newExrp.Expr)
		// если невалидное выражение
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid expression"))
			return
		}

		// Создаем и сохраняем задачу в задачах
		expression := tasker.Expression{
			UserTask: newExrp.Expr,
			Postfix:  postfix,
			Times:    *newExrp.Timings,
		}
		expression.CreateId()
		log.Printf("Method: %s, Expression: %s, Timings: %s, id: %d",
			r.Method,
			newExrp.Expr,
			newExrp.Timings.String(),
			expression.IdExpression,
		)

		ws.Expressions.Add(&expression)
		// Обновляем рабочее пространство
		ws.Update()
		// Сохраняем базу данных
		err = ws.Save()
		if err != nil {
			log.Println("ошибка сохранение после нового задания", err)
		}

		// Записываем тело ответа
		body := fmt.Sprintf("%d", expression.IdExpression)
		w.Write([]byte(body))
	}
}
