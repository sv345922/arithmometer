package newExpression

import (
	"arithmometer/internal/parser"
	"arithmometer/internal/wSpace"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Длительность по умолчанию

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
		var newClientExpression NewExpr
		err := json.NewDecoder(r.Body).Decode(&newClientExpression)
		defer r.Body.Close()
		if err != nil {
			log.Println("ошибка POST запроса")
			return
		}
		//Если тайминги не передаются, тогда они ставятся по умолчанию
		setTimingsWhileZero(&newClientExpression)

		// Обновляем полученные тайминги
		ws.Timings = newClientExpression.Timings
		// Парсим выражение, и проверяем его
		// Предполагается, что если парсинг с ошибкой, значит невалидное выражение
		newExpression := parser.NewExpression()
		err = newExpression.Do(newClientExpression.Expression, *newClientExpression.Timings)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid expression"))
			return
		}

		log.Printf("Method: %s, Expression: %s, Timings: %s, id: %d",
			r.Method,
			newExpression.UserTask,
			newExpression.Times.String(),
			newExpression.Id,
		)
		// Добавляем новое выражение в WorkingSpace
		// Добавляем в Expressions
		err = ws.AddToExpressions(TransformParseExpression(newExpression))
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("expression is already exist"))
			return
		}
		// Добавляем в AllNodes
		_, err = ws.AddToAllNodes(newExpression.Nodes)
		if err != nil {
			log.Println(err)
		}
		// Добавляем в очередь
		_ = ws.Queue.AddExpressionNodes(newExpression.Nodes)

		// Записываем тело ответа в виде id выражения
		body := fmt.Sprintf("%d", newExpression.Id)
		w.Write([]byte(body))

		//Сохраняем ws
		err = ws.Save()
		if err != nil {
			log.Println(err)
		}
	}
}
