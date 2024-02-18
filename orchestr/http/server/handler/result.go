package handler

import (
	"arithmometer/orchestr/inter/tasker"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// Обрабатывает запросы клиента о проверке результата вычислений
func GetResult(w http.ResponseWriter, r *http.Request) {
	// Получаем рабочее пространство из контекста
	ws, ok := tasker.GetWs(r.Context())
	if !ok {
		log.Println("ошибка контекста")
		return
	}
	expressions := ws.Expressions

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
	//idInt, _ := strconv.Atoi(id)
	// Поиск выражения в списке выражений
	expression := tasker.FindExpression(idInt, expressions)

	// выражение не найдено
	if expression == nil {
		w.Write([]byte("id не найден"))
		//w.WriteHeader(http.StatusNoContent)
		return
	}
	if expression.Calculated() {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("результат выражения %s = %f", expression.UserTask, expression.Result)))
		return
	}
	w.Write([]byte(fmt.Sprintf("выражение %s еще не посчитано", expression.UserTask)))

}
