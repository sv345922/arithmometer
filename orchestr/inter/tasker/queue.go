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

func NewDequeue() Dequeue {
	result := Dequeue{}
	result.Q = make([]*TaskContainer, 0)
	result.L = 0
	return result
}
func (d *Dequeue) AddBack(newVal *TaskContainer) {
	d.Q = append(d.Q, newVal)
	d.L++
}
func (d *Dequeue) AddFront(newVal *TaskContainer) {
	d.Q = append([]*TaskContainer{newVal}, d.Q...)
	d.L++
}
func (d *Dequeue) PopBack() (*TaskContainer, error) {
	if d.L == 0 {
		return nil, fmt.Errorf("пустая очередь")
	}
	result := d.Q[d.L-1]
	d.Q = d.Q[:d.L-1]
	d.L--
	return result, nil
}
func (d *Dequeue) PopFront() (*TaskContainer, error) {
	if d.L == 0 {
		return nil, fmt.Errorf("пустая очередь")
	}
	result := d.Q[0]
	d.Q = d.Q[1:]
	d.L--
	return result, nil
}
func (d *Dequeue) removeTask(idTask uint64) error {
	for i, element := range d.Q {
		if element.IdTask == idTask {
			d.Q = append(d.Q[:i], d.Q[i+1:]...)
			d.L--
			return nil
		}

	}
	return fmt.Errorf("пустая очередь, либо задача не найдена в очереди")
}

// Обновляет структуру очереди, чтобы в конце были только невзятые задачи, а в начале
// только взятые в обработку
func (d *Dequeue) Update() {
	var inWork []*TaskContainer
	var notInWork []*TaskContainer
	d.mu.Lock()
	defer d.mu.Unlock()

	for i := 0; i < d.L; i++ {
		// Если текущий элемент не вычисляется
		task := d.Q[i]
		if task.CalcId == 0 {
			// заносим его в список notInWork
			notInWork = append(notInWork, d.Q[i])
		} else {
			// иначе заносим его в список inWork
			// если вычислитель не вернул результат до дедлайна
			if task.Deadline.Before(time.Now()) {
				// установка дедлайна на будущее, можно и больше
				task.Deadline = time.Now().Add(time.Hour * 1000)
				// и сброс id вычислителя
				task.CalcId = 0
				notInWork = append(notInWork, task)
			}
			inWork = append(inWork, d.Q[i])
		}
	}
	d.Q = append(inWork, notInWork...)
}
