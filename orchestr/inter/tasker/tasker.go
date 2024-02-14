package tasker

import (
	"arithmometer/orchestr/parsing"
	"context"
	"fmt"
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

// Создает новый список выражений
// .Dict - словарь с ссылками на выражения
// .ListExpr - список с ссылками на выраженя, повторяет .Dict
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
// .Dict - словарь с задачами
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

/*
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
		// TODO доделать проверку на ошибки вычисления / можно не делать?
	}
}
*/

// Создает список задач и выражений и сохраняет его
func RunTasker() (*WorkingSpace, error) {
	tasks := NewTasks()
	expressions := NewExpressions()
	timings := &Timings{}
	// Восстатнавливаем выражения и задачи из базы данных
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
	workingSpace := WorkingSpace{
		Tasks:       tasks,
		Expressions: expressions,
		Timings:     timings,
		AllNodes:    &allNodes,
	}
	err = workingSpace.save()
	if err != nil {
		err = fmt.Errorf("ошибка сохранения БД: %v", err)
	}
	return &workingSpace, err
}

// Восстанавливает рабочее пространство из сохраненной базы данных
func restoreTaskExpr(tasks *Tasks, expressions *Expressions, timings *Timings) error {
	db, err := LoadJSON[DataBase]("db")
	if err != nil {
		log.Print(err)
	}
	if db.Timings == nil {
		timings = db.Timings
	}
	if db.Tasks != nil {
		tasks = db.Tasks
	}
	if db.Expressions != nil {
		expressions = db.Expressions
	}
	return nil
}

// Отдать ожидающие задачи на выполнение
