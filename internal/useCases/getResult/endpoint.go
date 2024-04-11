package getResult

import (
	"arithmometer/internal/wSpace"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// Обрабатывает запросы клиента о проверке результата вычислений
func GetResult(ws *wSpace.WorkingSpace) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

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

		// Обновление списка задач и выражений ?

		// преобразуем id в число
		idInt, _ := strconv.ParseUint(id, 10, 64)
		// Поиск выражения в списке выражений
		expression, ok := ws.Expressions[idInt]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("id не найден"))
			return
		}

		w.WriteHeader(http.StatusOK)
		switch expression.Status {
		case "done":
			w.Write([]byte(fmt.Sprintf("результат выражения %s = %f", expression.UserTask, expression.Result)))
			return
		case "zero division":
			w.Write([]byte(fmt.Sprintf("Выражение %s содержит ошибку деления на ноль", expression.UserTask)))
			return
		default:
			w.Write([]byte(fmt.Sprintf("выражение %s еще не посчитано", expression.UserTask)))
		}
	}
}
