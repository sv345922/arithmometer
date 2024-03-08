package tasker

import (
	"arithmometer/orchestr/parsing"
	"arithmometer/pkg/timings"
	"log"
	"sync"
	"time"
)

// TaskContainer Контейнер задачи для очереди задач
type TaskContainer struct {
	IdTask   uint64          `json:"id"`       // идентификатор задачи, передается вычислителю
	TaskN    Task            `json:"taskN"`    // задача
	Err      error           `json:"err"`      // ошибка
	TimingsN timings.Timings `json:"timingsN"` // тайминги
	CalcId   int             `json:"calcId"`   // id вычислителя задачи
	Deadline time.Time       `json:"deadline"`
	mu       sync.RWMutex    `json:"-"`
}

// GetID Возвращает id
func (tc *TaskContainer) GetID() uint64 {
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	return tc.IdTask
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

// SetDeadline устанавливает дедлайн задаче от текущего момента
func (tc *TaskContainer) SetDeadline(add time.Duration) {
	tc.mu.Lock()
	tc.Deadline = time.Now().Add(add)
	tc.mu.Unlock()
}

func (tc *TaskContainer) GetTiming() time.Duration {
	return tc.TaskN.Duration
}

// Task Зачада для вычислителя
type Task struct {
	X        float64       `json:"x"`        // операнд X
	Y        float64       `json:"y"`        // операнд Y
	Op       string        `json:"op"`       // операция
	Duration time.Duration `json:"duration"` // длительность операции
}

// Tasks Содержит задачи для вычислителей
// .Dict - словарь задач, ключ IdExpression
// .mu - мьютекс для блокировки словаря
type Tasks struct {
	Queue *Dequeue     `json:"queue"`
	mu    sync.RWMutex `json:"-"`
}

// NewTasks Создает новый список задач
// .Dict - словарь с задачами
// .Queue - очередь задач
func NewTasks() *Tasks {
	res := Tasks{}
	//res.Dict = make(map[uint64]*TaskContainer)
	res.Queue = NewDequeue()
	return &res
}

// AddTask Добавляет задачу в список задач
func (t *Tasks) AddTask(task *TaskContainer) {
	t.mu.Lock()
	defer t.mu.Unlock()
	// Добавляем в словарь
	//t.Dict[task.IdTask] = task
	// Добавляем в начало очередь
	t.Queue.AddFront(task)
}

// Проверка на наличие задачи в очереди
func (t *Tasks) isContent(node *parsing.Node) bool {
	id := node.NodeId
	for _, val := range t.Queue.Q {
		if val.IdTask == id {
			return true
		}
	}
	return false
}

// RemoveTask Удаляет задачу
func (t *Tasks) RemoveTask(idTask uint64) {
	// удаление из очереди
	t.Queue.removeTask(idTask) // Пока пропускаем ошибку
}

// Возвращает задачу (без удаления) для передачи ее вычислителю,
// также записывает id вычислителя в поле TaskContainer.CalcId
// и переставляет взятую задачу в начало очереди
func (t *Tasks) GetTask(calcId int) *TaskContainer {
	t.mu.Lock()
	defer t.mu.Unlock()
	// если очередь пуста, возвращаем nil
	if t.Queue.L == 0 {
		log.Println("очередь задач пустая")
		return nil
	}
	// берем элемент из начала очереди (последний в списке Q)
	task, _ := t.Queue.PopBack()
	// Если задача уже взята вычислителем возвращаем task на прежнее место и
	// возвращаем nil - все элементы в обработке
	if task.CalcId != 0 {
		t.Queue.AddBack(task)
		return nil
	}

	// если последний элемент очереди не взят вычислителем в обработку, возвращаем его
	// сам элемент кладем в начало очереди, ставим id вычислителя
	task.mu.Lock()
	task.CalcId = calcId
	task.mu.Unlock()
	t.Queue.AddFront(task)
	return task

}

//// Тайминги для операторов
//type Timings struct {
//	Plus  int `json:"plus"`
//	Minus int `json:"minus"`
//	Mult  int `json:"mult"`
//	Div   int `json:"div"`
//}
//
//// Стрингер
//func (t *Timings) String() string {
//	return fmt.Sprintf("+: %ds, -: %ds, *: %ds, /: %ds", t.Plus, t.Minus, t.Mult, t.Div)
//}

// Для получения выражения от клиента
type NewExpr struct {
	Expr    string           `json:"expr"`
	Timings *timings.Timings `json:"timings"`
}

// Для получения ответа на задачу
type Answer struct {
	Result float64 `json:"result"`
	Err    error   `json:"err"`
}
type AnswerContainer struct {
	Id      string `json:"id"`
	AnswerN Answer `json:"answerN"`
}
