package tasker

import (
	"arithmometer/orchestr/parsing"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
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
	res.Dict = make(map[int]*Expression)
	res.ListExpr = make([]*Expression, 0)
	return &res
}

// Добавляет выражение в список выражений
func (e *Expressions) Add(expression *Expression) {
	e.Dict[expression.IdExpression] = expression
	e.ListExpr = append(e.ListExpr, expression)
}

// возвращает выражение из списка задач
func FindExpression(id int, e *Expressions) *Expression {
	e.mu.RLock()
	defer e.mu.RUnlock()
	if task, ok := e.Dict[id]; ok {
		return task
	}
	return nil
}

// Создает новый список задач
// .Dict - словарь с задачами
// .Queue - очередь задач
func NewTasks() *Tasks {
	res := Tasks{}
	res.Dict = make(map[int]*TaskContainer)
	res.Queue = NewDequeue()
	return &res
}

// возвращает задачу из списка (мапы) задач
func FindNodes(id int, n map[int]*parsing.Node) *parsing.Node {
	if node, ok := n[id]; ok {
		return node
	}
	return nil
}

// Создает рабочее пространство и сохраняет базу данных
func RunTasker() (*WorkingSpace, error) {
	tasks := NewTasks()
	expressions := NewExpressions()
	timings := &Timings{}
	allNodes := make(map[int]*parsing.Node) // ключом является id узла
	/*
		// Проверка на существование сохраненной базы данных с созданием пустого файла
		if ok, err := checkDb(); !ok {
			if err != nil {
				return nil, err
			}
		}
	*/

	// Восстанавливаем выражения, задачи и тайминги из базы данных
	err := restoreTaskExpr(tasks, expressions, timings)
	if err != nil {
		log.Println("ошибка восстановления из БД", err)
	}
	// Если база данных не содержит задач, но содержит выражения
	// создаем очередь задач из выражений
	if tasks == nil && expressions != nil {
		for _, expression := range expressions.ListExpr {
			// строим дерево выражения
			nodes, root, err := parsing.GetTree(expression.Postfix)
			// Записываем в выражение ошибку, если она возникла при построении дерева
			// выражения
			if err != nil {
				expression.ParsError = err
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
					tasks.AddTask(TaskContainer{
						IdTask: node.NodeId,
						TaskN:  Task{X: node.X.Val, Y: node.X.Val, Op: node.Op},
					})
				}
			}
		}
	}
	workingSpace := WorkingSpace{
		Tasks:       tasks,
		Expressions: expressions,
		//Timings:     timings,
		AllNodes: allNodes,
	}
	err = workingSpace.Save()
	if err != nil {
		err = fmt.Errorf("ошибка сохранения БД: %v", err)
	}
	return &workingSpace, err
}

// Восстанавливает рабочее пространство из сохраненной базы данных
func restoreTaskExpr(tasks *Tasks, expressions *Expressions, timings *Timings) error {
	var result = &DataBase{
		Expressions: NewExpressions(),
		Tasks:       NewTasks(),
		//Timings:     &Timings{0, 0, 0, 0},
		mu: sync.Mutex{},
	}
	wd, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	path := wd + "\\orchestr\\db\\" + "db.json"
	data, err := os.ReadFile(path)
	if err != nil {
		log.Println("ошибка открытия json", err)
	}
	err = json.Unmarshal(data, result)
	if err != nil {
		log.Println(err)
	}

	db := result
	if err != nil {
		log.Print(err)
		return err
	}
	/*
		if db.Timings == nil {
			timings = db.Timings
		}
	
	*/
	if db.Tasks != nil {
		tasks = db.Tasks
	}
	if db.Expressions != nil {
		expressions = db.Expressions
	}
	return nil
}

// Проверяет существование базы данных и создает пустой файл при необходимости
func checkDb() (bool, error) {
	wd, err := os.Getwd()
	if err != nil {
		log.Println(err)
		return false, err
	}
	path := wd + "\\orchestr\\db\\db.json"

	fmt.Println("path=", path)

	_, fileError := os.Stat(path)
	if os.IsExist(fileError) {
		return true, nil
	}
	file, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Println("ошибка создания файла")
		return false, err
	}
	log.Println("создан файл базы данных")
	return true, file.Close()
}
