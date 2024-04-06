package taskQueue

import (
	"log"
	"sync"
	"time"
)

type Element interface {
	IsTimeout() bool
	SetDeadline(time.Duration)
	GetID() uint64
	SetCalc(uint64)
	IsReadyToCalc() bool
	String() string
}

// Tasks - Очередь задач.
// Waiting - задачи, готовые для выдачи вычислителям,
// Working - задачи, взятые вычислителем
// NotReady - задачи, ожидающие решения других задач
// WaitingIds - id готовых для вычисления задач
// L - количество элементов в очереди (всего)
type Tasks struct {
	AllTasks map[uint64]struct{} `json:"allTasks"` // мапа всех задач
	Waiting  []*Element          `json:"waiting"`  // готовые для выдачи вычислителям,
	Working  map[uint64]*Element `json:"working"`  // взятые вычислителем
	NotReady map[uint64]*Element `json:"norReady"` // ожидающие решения других задач
	L        uint                `json:"l"`        // количество элементов в очереди (всего)
	mu       sync.RWMutex
}

// NewTasks Возвращает указатель на новую очередь задач
func NewTasks() *Tasks {
	return &Tasks{
		AllTasks: make(map[uint64]struct{}),
		Waiting:  make([]*Element, 0),
		Working:  make(map[uint64]*Element),
		NotReady: make(map[uint64]*Element),
		L:        0,
		mu:       sync.RWMutex{},
	}
}

// AddTask Добавляет задачу в список задач NotReady и увеличивает счетчик L
func (ts *Tasks) AddTask(task *Element) bool {
	ind := (*task).GetID()
	if _, ok := ts.AllTasks[ind]; ok {
		return false
	}
	ts.NotReady[ind] = task
	ts.mu.Lock()
	ts.AllTasks[ind] = struct{}{}
	ts.L++
	ts.mu.Unlock()
	return true
}

// RemoveTask Удаляет задачу из очереди задач
func (ts *Tasks) RemoveTask(idTask uint64) bool {
	ts.mu.RLock()
	// проверяем наличие задачи в очереди
	if _, ok := ts.AllTasks[idTask]; !ok {
		ts.mu.RUnlock()
		return false
	}
	ts.mu.RUnlock()

	ts.mu.Lock()
	defer ts.mu.Unlock()

	// удаляем из Working
	if _, ok := ts.Working[idTask]; ok {
		delete(ts.Working, idTask)
		delete(ts.AllTasks, idTask)
		ts.L--
		return true
	}
	// удаляем из NotReady
	if _, ok := ts.NotReady[idTask]; ok {
		delete(ts.NotReady, idTask)
		delete(ts.AllTasks, idTask)
		ts.L--
		return true
	}
	// удалем из Waiting
	for i, task := range ts.Waiting {
		if (*task).GetID() == idTask {
			ts.Waiting = append(ts.Waiting[:i], ts.Waiting[i+1:]...) // TODO возможная ошибка
			delete(ts.AllTasks, idTask)
			ts.L--
			return true
		}
	}
	log.Printf("ошибка в очереди (элемент есть в AllTask но нет в других местах) при id=%d", idTask)
	return false
}


// GetTask
// Возвращает свободную задачу для вычислителя,
// переносит эту задачу в мапу работающих задач. При пустой очереди возвращает nil.
// Работа с таймингами и id вычислителя снаружи функции.
// Сначала обновляет очередь: проверяем в NotReady и если задача готова для вычисления
// переносим её в waiting.
func (ts *Tasks) GetTask() *Element {
	// обновляем очереди
	_ = ts.CheckDeadlines()
	ts.mu.Lock()
	for key, val := range ts.NotReady {
		if (*val).IsReadyToCalc() {
			ts.Waiting = append(ts.Waiting, val)
			delete(ts.NotReady, key) // TODO maybe error
		}
	}
	ts.mu.Unlock()

	// берем первый элемент из очереди
	ts.mu.Lock()
	defer ts.mu.Unlock()
	l := len(ts.Waiting)
	// если очередь пустая возвращаем nil
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
	// и переносим в мапу ожидающих решения
	ts.Working[(*result).GetID()] = result

	return result
}

// CheckDeadlines Функция обновления очереди по состоянию таймингов
// Если среди работающих задач есть с простроченным дедлайном,
// то задача переносится в список ожидающих
func (ts *Tasks) CheckDeadlines() int {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	// получаем список ключей
	keys := make([]uint64, len(ts.Working))
	i := 0
	for k := range ts.Working {
		keys[i] = k
		i++
	}
	n := 0
	for _, key := range keys {
		task := ts.Working[key]
		// если задача с прошедшим дедлайном
		if (*task).IsTimeout() {
			// увеличиваем счетчик просроченных
			n++
			// устанавливаем дедлайн в далекое будущее
			(*task).SetDeadline(time.Hour * 1000)
			// и перемещаем задачу в начало очереди ожидающих
			ts.Waiting = append([]*Element{task}, ts.Waiting...)
			delete(ts.Working, (*task).GetID())
		}
	}
	return n
}

