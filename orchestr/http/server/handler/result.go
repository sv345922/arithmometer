package handler

import (
	"arithmometer/orchestr/inter/tasker"
	"fmt"
	"log"
	"net/http"
)

// Обрабатывает запросы клиента о проверке результата вычисллений
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
		log.Println("не найден id в запросе")
		return
	}
	log.Println("Id =", id)

	// Обновление списка выражений
	ws.Update()
	// Поиск выражения в списке выражений
	expression := tasker.FindExpression(id, expressions)

	// выражение не найдено
	if expression == nil {
		w.WriteHeader(http.StatusNoContent)
		w.Write([]byte("id не найден"))
		return
	}
	if expression.Calculated() {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(expression.Result))
		return
	}
	w.Write([]byte(fmt.Sprintf("выражение %s еще не посчитано", expression.UserTask)))

}
