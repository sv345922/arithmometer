package handler

import (
	"arithmometer/internal/wSpace"
	"arithmometer/pkg/expressions"
	"arithmometer/pkg/parser"
	"arithmometer/pkg/timings"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Обработчик создания нового выражения
func NewExpression(ws *wSpace.WorkingSpace) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Проверить что это запрос POST
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("требуется метод POST"))
			return
		}
		// Читаем тело запроса, в котором записано выражение и тайминги операций
		var newExrp wSpace.NewExpr
		err := json.NewDecoder(r.Body).Decode(&newExrp)
		defer r.Body.Close()
		if err != nil {
			log.Println("ошибка POST запроса")
			return
		}
		//Если тайминги не передаются, тогда они ставятся по умолчанию
		if newExrp.Timings == nil {
			newExrp.Timings = &timings.Timings{Plus: 1, Minus: 1, Mult: 1, Div: 1}
		}
		// Парсим выражение, и проверяем его
		// Предполагается, что если парсинг с ошибкой, значит невалидное выражение
		postfix, err := parser.Parse(newExrp.Expr)
		// если невалидное выражение
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid expression"))
			return
		}
		// Создаем выражение
		expression := expressions.Expression{
			UserTask: newExrp.Expr,
			Postfix:  postfix,
			Times:    *newExrp.Timings,
			//RootId:   root.NodeId,
		}
		// создаем id выражения
		expression.CreateId()

		log.Printf("Method: %s, Expression: %s, Timings: %s, id: %d",
			r.Method,
			expression.UserTask,
			expression.Times.String(),
			expression.IdExpression,
		)
		// добавляем выражение в список выражений
		// также выполняем необходимые действия с рабочей структурой
		err = ws.AddExpression(&expression)
		if err != nil {
			log.Println("ошибка добавления нового задания", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Выражение не может быть вычислено"))
		}
		err = ws.Save()
		if err != nil {
			log.Println("ошибка сохранение после нового задания", err)
		}
		// Записываем тело ответа
		body := fmt.Sprintf("%d", expression.IdExpression)
		w.Write([]byte(body))
	}
}
