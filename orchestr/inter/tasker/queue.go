package tasker

import (
	"fmt"
	"sync"
	"time"
)

type Dequeue struct {
	Q  []*TaskContainer `json:"q"` // Очередь
	L  int              `json:"l"` // Длина очереди
	mu sync.Mutex
}

func NewDequeue() *Dequeue {
	result := Dequeue{}
	result.Q = make([]*TaskContainer, 0)
	result.L = 0
	return &result
}
func (d *Dequeue) AddBack(newVal *TaskContainer) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.Q = append(d.Q, newVal)
	d.L++
}
func (d *Dequeue) AddFront(newVal *TaskContainer) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.Q = append([]*TaskContainer{newVal}, d.Q...)
	d.L++
}
func (d *Dequeue) PopBack() (*TaskContainer, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.L == 0 {
		return nil, fmt.Errorf("пустая очередь")
	}
	result := d.Q[d.L-1]
	d.Q = d.Q[:d.L-1]
	d.L--
	return result, nil
}
func (d *Dequeue) PopFront() (*TaskContainer, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.L == 0 {
		return nil, fmt.Errorf("пустая очередь")
	}
	result := d.Q[0]
	d.Q = d.Q[1:]
	d.L--
	return result, nil
}
func (d *Dequeue) removeTask(idTask uint64) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	for i, element := range d.Q {
		if element.IdTask == idTask {
			d.Q = append(d.Q[:i], d.Q[i+1:]...)
			d.L--
			return nil
		}
	}
	return fmt.Errorf("задача не найдена в очереди")
}

// Обновляет структуру очереди, перемещая задачи с прошедшим таймаутом в начало очереди
// и сбрасывая id их вычислителя
func (d *Dequeue) UpdateWithTimeouts() {
	d.mu.Lock()
	defer d.mu.Unlock()
	var waitingTasks []*TaskContainer

	// Перебираем задачи из конца очереди (ожидающие результата от вычислителя)
	// и проверяем их дедлайны
	for currentTask, err := d.PopFront(); err == nil && currentTask.CalcId != 0; currentTask, err = d.PopFront() {
		// если дедлайн прошел обновляем задачу и ставим в начало очереди
		if currentTask.Deadline.Before(time.Now()) {
			// установка дедлайна на будущее, можно и больше
			currentTask.Deadline = time.Now().Add(time.Hour * 1000)
			// и сброс id вычислителя
			currentTask.CalcId = 0
			d.AddBack(currentTask)
			continue
		} else {
			// иначе заносим во временный список ожидающих результата вычисления
			waitingTasks = append(waitingTasks, currentTask)
		}
	}
	// Переносим элементы из временного списка в конец очереди
	for _, task := range waitingTasks {
		d.AddFront(task)
	}
}
