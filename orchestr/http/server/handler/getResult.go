package handler

import (
	"arithmometer/internal/wSpace"
	"arithmometer/pkg/expressions"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// Обрабатывает запросы клиента о проверке результата вычислений
func GetResult(ws *wSpace.WorkingSpace) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		exprs := ws.Expressions

		// Проверить метод
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("требуется метод Get"))
			return
		}
		// Читаем id из параметров запроса
		id := r.URL.Query().Get("id")
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("не найден id в запросе"))
			log.Println("не найден id в запросе")
			return
		}
		log.Println("Id запрошенного выражения =", id)

		// Обновление списка задач и выражений

		// преобразуем id в число
		idInt, _ := strconv.ParseUint(id, 10, 64)
		// Поиск выражения в списке выражений
		expression := expressions.FindExpression(idInt, exprs)

		// выражение не найдено
		if expression == nil {
			w.Write([]byte("id не найден"))
			//w.WriteHeader(http.StatusNoContent)
			return
		}
		w.WriteHeader(http.StatusOK)
		if expression.Calculated() {
			w.Write([]byte(fmt.Sprintf("результат выражения %s = %f", expression.UserTask, expression.Result)))
			return
		}
		w.Write([]byte(fmt.Sprintf("выражение %s еще не посчитано", expression.UserTask)))
	}
}
