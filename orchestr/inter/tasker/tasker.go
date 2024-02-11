package tasker

import (
	"arithmometer/orchestr/parsing"
	"context"
	"fmt"
	"log"
	"sync"
)

// Достает список узлов из контекста
func GetNodes(ctx context.Context) (*map[string]*parsing.Node, error) {
	mu := sync.Mutex{}
	mu.Lock()
	defer mu.Unlock()
	value := ctx.Value("nodes")
	switch res := value.(type) {
	case *map[string]*parsing.Node:
		return res, nil
	default:
		return nil, fmt.Errorf("ошибка контекста (nodes)")
	}
}

// Достает тайминги из контекста
func GetTimings(ctx context.Context) (*Timings, error) {
	mu := sync.Mutex{}
	mu.Lock()
	defer mu.Unlock()
	value := ctx.Value("timings")
	switch res := value.(type) {
	case *Timings:
		return res, nil
	default:
		return nil, fmt.Errorf("ошибка контекста (timings)")
	}
}

// Достает список задач из контекста
func GetTasks(ctx context.Context) (*Tasks, error) {
	mu := sync.Mutex{}
	mu.Lock()
	defer mu.Unlock()
	value := ctx.Value("tasks")
	switch res := value.(type) {
	case *Tasks:
		return res, nil
	default:
		return nil, fmt.Errorf("ошибка контекста (tasks)")
	}
}

// Дстает список выражений из контекста
func GetExpressions(ctx context.Context) (*Expressions, error) {
	mu := sync.Mutex{}
	mu.Lock()
	defer mu.Unlock()
	value := ctx.Value("expressions")
	switch res := value.(type) {
	case *Expressions:
		return res, nil
	default:
		return nil, fmt.Errorf("ошибка контекста (expressions)")
	}
}

// Создает новый список выражений
func NewExpressions() *Expressions {
	res := Expressions{}
	res.Dict = make(map[string]*Expression)
	res.ListExpr = make([]*Expression, 0)
	return &res
}

// Добавляет выражение в список выражений
func (e *Expressions) Add(expression *Expression) {
	e.Dict[expression.Id] = expression
	e.ListExpr = append(e.ListExpr, expression)
}

// Проверяет задачи, при обнаружении завершенной обновляет её статус
// TODO
func (e *Expressions) Update() {
	e.mu.RLock()
	defer e.mu.RUnlock()
	for index, expr := range e.ListExpr {
		if expr.Calculated() {

		}

	}
}

// возвращает выражение из списка задач
func FindExpression(id string, e *Expressions) *Expression {
	e.mu.RLock()
	defer e.mu.RUnlock()
	if task, ok := e.Dict[id]; ok {
		return task
	}
	return nil
}

// Создает новый список задач
func NewTasks() *Tasks {
	res := Tasks{}
	res.Dict = make(map[string]*TaskContainer)
	return &res
}

// возвращает задачу из списка (мапы) задач
func FindNodes(id string, n map[string]*parsing.Node) *parsing.Node {
	if node, ok := n[id]; ok {
		return node
	}
	return nil
}

// Проверяет задачу, при её завершении удаляет
func (t *Tasks) Update(*map[string]*parsing.Node) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	for id, task := range t.Dict {
		if task.Calculated() {
			t.mu.Lock()
			delete(t.Dict, id)
			t.mu.Unlock()
		}
		// TODO доделать проверку на ошибки вычисления
	}
}

// Создает список задач и выражений
func RunTasker() (*Tasks, *Expressions, Timings, map[string]*parsing.Node) {
	tasks := NewTasks()
	expressions := NewExpressions()
	// Восстатнавливаем выражения и задачи из базы данных
	timings := Timings{}
	err := restoreTaskExpr(tasks, expressions, timings)
	if err != nil {
		log.Println("ошибка восстановления из БД", err)
	}
	allNodes := make(map[string]*parsing.Node) // ключом является id узла
	if tasks == nil {
		// строим дерево выражения
		for _, expression := range expressions.ListExpr {
			nodes, root, err := parsing.GetTree(expression.Postfix)
			// Записываем в выражение ошибку, если она возникла при построении дерева выражения
			if err != nil {
				expression.Result = err.Error()
				continue
			}
			expression.Root = root
			// создаем дерево задач
			for _, node := range nodes {
				// Создаем ID для узлов
				node.CreateId()
				// заполняем словарь узлами
				allNodes[node.NodeId] = node
				// Если узел не рассчитан
				if !node.Calculated {
					// добавляем его в таски
					tasks.Add(TaskContainer{
						Id:    node.NodeId,
						TaskN: Task{X: node.X.Val, Y: node.X.Val, Op: node.Op},
					})
				}
			}
		}
	}
	return tasks, expressions, timings, allNodes
}
func restoreTaskExpr(tasks *Tasks, expressions *Expressions, timings Timings) error {
	db, err := LoadJSON[DataBase]("db")
	if err != nil {
		log.Print(err)
	}
	if db.Timings != nil {
		timings = *db.Timings
	}
	if db.Tasks != nil {
		tasks = db.Tasks
	}
	if db.Expressions != nil {
		expressions = db.Expressions
	}
	return nil
}

// Создать список узлов выражения и построить списки задач на выполнение

// Отдать ожидающие задачи на выполнение

// При получении ответа от вычислителей обновить очереди

// При завершении вычисления выражения сохранить результат
