package tasker

import (
	"arithmometer/orchestr/parsing"
	"fmt"
	"log"
	"sync"
	"time"
)

// Контейнер задач для отправки вычислителю
type TaskContainer struct {
	// идентификатор задачи, передается вычислителю
	IdTask   uint64     `json:"id"`
	TaskN    Task       `json:"taskN"`    // задача
	Err      error      `json:"err"`      // ошибка
	TimingsN Timings    `json:"timingsN"` // тайминги
	CalcId   int        `json:"calcId"`   // id вычислителя задачи
	Deadline time.Time  `json:"deadline"`
	mu       sync.Mutex `json:"-"`
}

// Зачада для вычислителя
type Task struct {
	X  float64 `json:"x"`  // операнд X
	Y  float64 `json:"y"`  // операнд Y
	Op string  `json:"op"` // операция
}

// Содержит задачи для вычислителей
// .Dict - словарь задач, ключ IdExpression
// .mu - мьютекс для блокировки словаря
type Tasks struct {
	Queue Dequeue                   `json:"queue"`
	Dict  map[uint64]*TaskContainer `json:"dict"` // ключ IdTask
	mu    sync.RWMutex              `json:"-"`
}

// Добавляет задачу в список задач
func (t *Tasks) AddTask(task TaskContainer) {
	t.mu.Lock()
	defer t.mu.Unlock()
	// Добавляем в словарь
	t.Dict[task.IdTask] = &task
	// Добавляем в начало очередь
	t.Queue.AddFront(&task)
	t.Queue.Update()
}

// Проверка на наличие задачи в очереди
func (t *Tasks) isContent(node *parsing.Node) bool {
	id := node.NodeId
	if _, ok := t.Dict[id]; ok {
		return true
	}
	return false
}

// Удаляет задачу
func (t *Tasks) RemoveTask(idTask uint64) {
	t.mu.Lock()
	defer t.mu.Unlock()
	// удаление из очереди
	t.Queue.removeTask(idTask) // Пока пропускаем ошибку
	// удаление из словаря
	delete(t.Dict, idTask)
}

// Возвращает задачу (без удаления) для передачи ее вычислителю,
// также записывает id вычислителя в поле TaskContainer.CalcId
// и переставляет взятую задачу в начало очереди
func (t *Tasks) GetTask(calcId int) *TaskContainer {
	t.mu.Lock()
	defer t.mu.Unlock()
	// берем последний элемент очереди
	task, err := t.Queue.PopBack()
	// если очередь пуста, возвращаем nil
	if err != nil {
		log.Println(err)
		return nil
	}
	// если последний элемент очереди не взят вычислителем в обработку, возвращаем его
	task.mu.Lock()
	if task.CalcId != 0 {
		task.CalcId = calcId
		task.mu.Unlock()
		t.Queue.AddFront(task)
		return task
	}
	task.mu.Unlock()
	// иначе возвращаем task на прежнее место и
	// возвращаем nil - очередь пуста, все элементы в обработке
	t.Queue.AddBack(task)
	return nil
}

// Содержит выражения пользователя
// .Dict - словарь ссылок на выражения, ключ IdExpression
// .ListExpr - список с ссылками на выражения
// .mu - мьютекс для блокировки
type Expressions struct {
	Dict     map[uint64]*Expression `json:"dict"` // ключ id запроса выражения/запроса клиента
	ListExpr []*Expression          `json:"listExpr"`
	mu       sync.RWMutex           `json:"-"`
}

// Обновляет список выражений, при вычисленном корне выражения, либо делении на ноль,
// ставит статус вычислено/деление на ноль
// и результат вычислений
func (e *Expressions) UpdateStatus(root *parsing.Node, status string, result float64) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	for _, expression := range e.ListExpr {
		expression.mu.Lock()
		if expression.RootId == root.NodeId {
			expression.Status = status
			expression.Result = result
			expression.mu.Unlock()
			return
		}
		expression.mu.Unlock()
	}
}

/*
	func (e *Expressions) Remove(id string) error {
		e.mu.Lock()
		defer e.mu.Unlock()
		delete(e.dict, id)
		e.listExpr = append(e.listExpr[:index], e.listExpr[index+1:]...)
	}
*/
// Выражение
type Expression struct {
	IdExpression uint64            `json:"id"`        // id запроса клиента
	UserTask     string            `json:"userTask"`  // задание клиента
	Postfix      []*parsing.Symbol `json:"postfix"`   // постфиксная запись выражения
	Times        Timings           `json:"times"`     // тайминги
	Result       float64           `json:"result"`    // результат,
	Status       string            `json:"status"`    // ""/"done"/"деление на ноль"
	RootId       uint64            `json:"rootId"`    // id кореневого узла
	ParsError    error             `json:"parsError"` // Ошибка парсинга
	mu           sync.Mutex
}

// Создает id выражения
func (e *Expression) CreateId() {
	s := ""
	for _, symbol := range e.Postfix {
		s = s + symbol.Val
	}
	e.IdExpression = parsing.GetId(s)

	/*

		hasher := sha1.New()
		hasher.Write([]byte(s))
		e.IdExpression = base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	*/
}

// Определяет, вычислено ли выражение
func (e *Expression) Calculated() bool {
	if e.Status == "done" {
		return true
	}
	return false
}

// Тайминги для операторов
type Timings struct {
	Plus  int `json:"plus"`
	Minus int `json:"minus"`
	Mult  int `json:"mult"`
	Div   int `json:"div"`
}

// Стрингер
func (t *Timings) String() string {
	return fmt.Sprintf("+: %d s, -: %d s, *: %d s, /: %d s", t.Plus, t.Minus, t.Mult, t.Div)
}

type NewExpr struct {
	Expr    string   `json:"expr"`
	Timings *Timings `json:"timings"`
}
type Answer struct {
	Result float64 `json:"result"`
	Err    error   `json:"err"`
}
type AnswerContainer struct {
	Id      string `json:"id"`
	AnswerN Answer `json:"answerN"`
}
