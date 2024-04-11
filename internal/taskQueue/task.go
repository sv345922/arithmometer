package taskQueue

import (
	"arithmometer/internal/entities"
	"sync"
	"time"
)

type Task struct {
	Node     *entities.Node `json:"node"`
	X        float64        `json:"x"`
	XReady   bool           `json:"xready"`
	Y        float64        `json:"y"`
	YReady   bool           `json:"yready"`
	Error    error          `json:"error"`    // ошибка
	CalcId   uint64         `json:"calc_id"`  // id вычислителя задачи
	Deadline time.Time      `json:"deadline"` // дедлайн задачи
	Duration time.Duration  `json:"duration"` // длительность операции
	mu       sync.RWMutex
}

func NewTask(node *entities.Node) *Task {
	return &Task{Node: node, Deadline: time.Now().Add(time.Hour * 1000)}
}

// GetID Возвращает id
func (t *Task) GetID() uint64 {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.Node.Id
}

// SetCalc Присваивает id вычислителя
func (t *Task) SetCalc(calcId uint64) {
	t.CalcId = calcId
}

// Проверка на завершение дедлайна задачи, если время вышло, возвращает true
func (t *Task) IsTimeout() bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	if t.Deadline.Before(time.Now()) {
		return true
	}
	return false
}

// проверяет готовность задачи для расчетов
func (t *Task) IsReadyToCalc() bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	if !t.Node.Calculated {
		if t.XReady && t.YReady {
			return true
		}
	}
	return false
}

// SetDeadline устанавливает дедлайн задаче от текущего момента
func (t *Task) SetDeadline(add time.Duration) {
	t.mu.Lock()
	t.Deadline = time.Now().Add(add)
	t.mu.Unlock()
}

func (t *Task) SetDuration(duration time.Duration) {
	t.mu.Lock()
	t.Duration = duration
	t.mu.Unlock()
}

// Заносит результат выражения в узел
func (t *Task) SetResult(result float64) {
	t.mu.Lock()
	t.Node.Calculated = true
	t.Node.Val = result
	t.mu.Unlock()
}
