package tasker

import (
	"arithmometer/orchestr/parsing"
	"fmt"
	"log"
)

type WorkingSpace struct {
	Tasks       *Tasks                `json:"tasks"`
	Expressions *Expressions          `json:"expressions"`
	Timings     *Timings              `json:"timings"`
	AllNodes    map[int]*parsing.Node // ключ id узла
}

// Сохраняет рабочее пространство
func (ws *WorkingSpace) Save() error {
	db := DataBase{
		Expressions: ws.Expressions,
		Tasks:       ws.Tasks,
		//Timings:     ws.Timings,
	}
	db.mu.Lock()
	defer db.mu.Unlock()
	err := SafeJSON[DataBase]("db", db)
	if err != nil {
		return err
	}
	log.Println("база данных сохранена")
	return nil
}

// При получении выполненого задания
// Проверяем на наличие ошибки деления на ноль
// Обновляет очередь задач
// Проверяет список выражений и если оно вычислено, обновляет его статус
// Добавляет новую задачу в начало очереди задач, если выражение не корень выражения
func (ws *WorkingSpace) UpdateTasks(IdTask int, answer *Answer) error {
	currentNode := ws.AllNodes[IdTask]
	// Проверка деления на ноль и обновление выражения
	// с удалением не требующих решения задач,
	// а также изменение статуса выражения
	if answer.Err != nil {
		currentNode.ErrZeroDiv = answer.Err
		ws.updateWhileZero(currentNode)
		return nil
	}
	result := answer.Result
	// Удаляем задачу из очереди
	ws.Tasks.RemoveTask(IdTask)
	// удаляем задачу из словаря узлов TODO (надо ли)
	delete(ws.AllNodes, IdTask)
	// записываем результат вычисления в узел
	currentNode.Val = result
	currentNode.Calculated = true

	// Проверяем готовность родительского узла и добавляем его в очередь задач при готовности
	// TODO вроде сделано
	parent := currentNode.Parent
	// Если это корень выражения
	if parent == nil {
		// Обновляем результат выражения и его статус
		ws.Expressions.UpdateStatus(currentNode, "done", result)
		return nil
	}
	// проверка готовности узла и добавление в очередь задач
	for checkAndUpdateParent(ws, parent) {
		parent = parent.Parent
		if parent == nil {
			ws.Expressions.UpdateStatus(parent, "done", 0)
			break
		}
	}
	// ws.Tasks.Queue.Update()
	return nil
}
func (ws *WorkingSpace) Update() {
	for _, expression := range ws.Expressions.ListExpr {
		// строим дерево выражения
		nodes, root, err := parsing.GetTree(expression.Postfix)
		// Записываем в выражение ошибку, если она возникла при построении дерева
		// выражения
		if err != nil {
			expression.ParsError = err
			continue
		}
		expression.Root = root
		// Создаем дерево задач
		for _, node := range nodes {
			// Создаем ID для узлов
			node.CreateId()
			// заполняем словарь узлами
			ws.AllNodes[node.NodeId] = node
			// Если узел не рассчитан
			if !node.Calculated {
				// добавляем его в таски
				ws.Tasks.AddTask(TaskContainer{
					IdTask: node.NodeId,
					TaskN:  Task{X: node.X.Val, Y: node.X.Val, Op: node.Op},
				})
			}
		}
	}
}

// Проверяет на готовность родительский узел, при готовности добавляет его в очередь задач
func checkAndUpdateParent(ws *WorkingSpace, parent *parsing.Node) bool {
	parent.Mu.RLock()
	defer parent.Mu.RUnlock()
	if parent.X.Calculated && parent.Y.Calculated {
		task := TaskContainer{
			IdTask: parent.NodeId,
			TaskN: Task{
				X:  parent.X.Val,
				Y:  parent.Y.Val,
				Op: parent.Op,
			},
			Err:      nil,
			TimingsN: *ws.Timings,
		}
		ws.Tasks.AddTask(task)
		return true
	}
	return false
}

// Обновляет рабочее пространство при обнаружении деления на ноль,
// Проверяет узлы в дереве выражения и обновляет их,
// Удаляет задачи из очереди задач,
// Сохраняет базу данных
func (ws *WorkingSpace) updateWhileZero(node *parsing.Node) {
	// изменяем поле ошибка родительского узла и его дочерних узлов
	rootNode := errorUpdate(node)
	// Изменяем статус выражения с ошибкой
	ws.Expressions.UpdateStatus(rootNode, "деление на ноль в выражении", 0)
	// Удаляем узлы с ошибкой
	for key, val := range ws.AllNodes {
		if val.ErrZeroDiv != nil {
			delete(ws.AllNodes, key)
		}
	}
	// Удаляем задачи с ошибкой
	for key, val := range ws.Tasks.Dict {
		if val.Err != nil {
			ws.Tasks.RemoveTask(key)
		}
	}
	// Сохранение рабочего пространства
	ws.Save()
}

// Обновляет статус ошибки узла вниз (до листа) рекурсивно
func errorUpdateToSheet(n *parsing.Node) {
	// если узел это лист, или ветка узла уже обработана
	if n.Sheet || n.ErrZeroDiv != nil {
		return
	}
	n.ErrZeroDiv = fmt.Errorf("деление на ноль в выражении")
	errorUpdateToSheet(n.X)
	errorUpdateToSheet(n.Y)
}

// Обновляет статус ошибки узла вверх (до корня), и соседних ветвей
// и возвращает корень выражения. Работает рекурсивно
func errorUpdate(n *parsing.Node) *parsing.Node {
	if n.Parent == nil {
		n.ErrZeroDiv = fmt.Errorf("деление на ноль в выражении")
		return n
	}
	n.ErrZeroDiv = fmt.Errorf("деление на ноль в выражении")
	errorUpdateToSheet(n.X)
	errorUpdateToSheet(n.Y)
	return errorUpdate(n.Parent)
}
