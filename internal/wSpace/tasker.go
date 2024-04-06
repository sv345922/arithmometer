package wSpace

import (
	"arithmometer/pkg/treeExpression"
	"sync"
	"time"
)

type Task struct {
	X        float64       `json:"x"`        // операнд X
	Y        float64       `json:"y"`        // операнд Y
	Op       string        `json:"op"`       // операция
	Duration time.Duration `json:"duration"` // длительность операции
}

// TaskContainer Контейнер задачи для очереди задач
type TaskContainer struct {
	treeExpression.Node
	TaskAn   Task         `json:"taskN"`  // задача
	Err      error        `json:"err"`    // ошибка
	CalcId   int          `json:"calcId"` // id вычислителя задачи
	Deadline time.Time    `json:"deadline"`
	mu       sync.RWMutex `json:"-"`
}

// GetID Возвращает id
func (tc *TaskContainer) GetID() uint64 {
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	return tc.NodeId
}

// SetCalc Присваивает id вычислителя
func (tc *TaskContainer) SetCalc(calcId int) {
	tc.CalcId = calcId
}

// Проверка на завершение дедлайна задачи, если время вышло, возвращает true
func (tc *TaskContainer) IsTimeout() bool {
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	if tc.Deadline.Before(time.Now()) {
		return true
	}
	return false
}

// проверяет готовность задачи для расчетов
func (tc *TaskContainer) IsReadyToCalc() bool {
	return tc.IsReadyToCalc()
}

// SetDeadline устанавливает дедлайн задаче от текущего момента
func (tc *TaskContainer) SetDeadline(add time.Duration) {
	tc.mu.Lock()
	tc.Deadline = time.Now().Add(add)
	tc.mu.Unlock()
}

func (tc *TaskContainer) GetTiming() time.Duration {
	return tc.TaskAn.Duration
}

// Task Зачада для вычислителя
