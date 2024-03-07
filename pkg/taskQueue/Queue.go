package taskQueue

import (
	"arithmometer/orchestr/inter/tasker"
	"sync"
)

// Tasks - Очередь задач.
// Waiting - задачи, готовые для выдачи вычислителям,
// Working - задачи, взятые вычислителем
// WaitingIds - id готовых для вычисления задач
// L - количество элементов в очереди (всего)
type Tasks struct {
	Waiting    []*tasker.TaskContainer          `json:"waiting"`
	Working    map[uint64]*tasker.TaskContainer `json:"working"`
	WaitingIds map[uint64]struct{}              `json:"waitingIds"`
	L          uint
	mu         sync.Mutex
}

// NewTasks Возвращает указатель на новую очередь задач
func NewTasks() *Tasks {
	return &Tasks{
		Waiting:    make([]*tasker.TaskContainer, 0, 100),
		WaitingIds: make(map[uint64]struct{}),
		Working:    make(map[uint64]*tasker.TaskContainer),
		L:          0,
		mu:         sync.Mutex{},
	}
}

// AddTask Добавляет задачу в список задач
func (ts *Tasks) AddTask(task *tasker.TaskContainer) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	ts.Waiting = append(ts.Waiting, task)
	ts.WaitingIds[task.GetID()] = struct{}{}
	ts.L++
}

// RemoveTask Удаляет задачу из очереди задач
func (ts *Tasks) RemoveTask(idTask uint64) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	if _, ok := ts.Working[idTask]; ok {
		delete(ts.Working, idTask)
		ts.L--
		return
	}
	if _, ok := ts.WaitingIds[idTask]; ok {
		for i, task := range ts.Waiting {
			if task.GetID() == idTask {
				ts.Waiting = append(ts.Waiting[:i], ts.Waiting[i+1:]...)
				delete(ts.WaitingIds, idTask)
				ts.L--
			}
		}
	}
}

// GetTask
// Возвращает свободную задачу для вычислителя,
// переносит эту задачу в мапу работающих задач. При пустой очереди возвращает nil
func (ts *Tasks) GetTask(calcId int) *tasker.TaskContainer {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	l := len(ts.Waiting)
	// если очередь пустая
	if l == 0 {
		return nil
	}
	// Получаем задачу - первый элемент очереди
	result := ts.Waiting[0]
	switch l {
	case 1:
		ts.Waiting = ts.Waiting[:0] // при длине 1 опустошаем очередь
	default:
		ts.Waiting = ts.Waiting[1:] // иначе оставляем очередь без первого элемента
	}
	// удаляем id из мапы ожидающих
	delete(ts.WaitingIds, result.GetID())
	// добавляем выданную задачу в обрабатываемые
	ts.Working[result.GetID()] = result
	// Устанавливаем id калькулятора выданной задаче
	result.SetCalc(calcId)
	return result
}
