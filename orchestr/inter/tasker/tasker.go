package tasker

import (
	"context"
	"log"
	"sync"
)

// Достает рабочее пространство из контекста
func GetWs(ctx context.Context) (*WorkingSpace, bool) {
	mu := sync.Mutex{}
	mu.Lock()
	defer mu.Unlock()
	ws, ok := ctx.Value("ws").(*WorkingSpace)
	return ws, ok
}

// Создает рабочее пространство из сохраненной базы данных
func RunTasker() (*WorkingSpace, error) {
	// Восстанавливаем выражения и задачи из базы данных
	// Загрузка сохраненной БД
	ws, err := LoadDB()
	if err != nil {
		log.Println("ошибка загрузки БД", err)
		return ws, err
	}
	return ws, err
}
