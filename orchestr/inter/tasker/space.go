package tasker

import (
	"arithmometer/orchestr/parsing"
	"fmt"
	"log"
	"time"
)

type WorkingSpace struct {
	Tasks       *Tasks                   `json:"tasks"`
	Expressions *Expressions             `json:"expressions"`
	Timings     *Timings                 `json:"timings"`
	AllNodes    map[uint64]*parsing.Node // ключ id узла
}

// Сохраняет рабочее пространство
func (ws *WorkingSpace) Save() error {
	db := DataBase{
		Expressions: ws.Expressions,
		Tasks:       ws.Tasks,
		//Timings:     ws.Timings,
	}

	err := SafeJSON[DataBase]("db", db)
	if err != nil {
		return err
	}
	log.Println("DB saved")
	return nil
}

// При получении выполненого задания,
// проверяем на наличие ошибки деления на ноль.
// Записывает результат в узел. и изменяет статус на вычислено
// Обновляет очередь задач.
// Проверяет список выражений и если оно вычислено, обновляет его статус.
// Добавляет новую задачу в начало очереди задач.
func (ws *WorkingSpace) UpdateTasks(IdTask uint64, answer *Answer) error {
	defer ws.Tasks.Queue.Update()
	currentNode, ok := ws.AllNodes[IdTask]
	if !ok {
		log.Println("не найден узел")
	}
	// Проверка деления на ноль и обновление выражения
	// с удалением не требующих решения задач,
	// а также изменение статуса выражения
	if answer.Err != nil {
		log.Println("в выражении присутсвует деление на ноль")
		currentNode.ErrZeroDiv = answer.Err
		ws.updateWhileZero(currentNode)
		return nil
	}
	result := answer.Result
	// Удаляем задачу из очереди
	ws.Tasks.RemoveTask(IdTask)
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
	for checkAndUpdateNodeToTasks(ws, parent) {
		if parent.Parent == nil {
			// ws.Expressions.UpdateStatus(parent, "done", 0)
			break
		}
		parent = parent.Parent
	}
	return nil
}

// проходит по списку выражений, создает дерево узлов выражения,
// включает в рабочее пространство список узлов - ws.AllNodes
// созадет очередь задач для вычислителей - ws.tasks
func (ws *WorkingSpace) Update() {
	// Взять выражения
	// проверить на существование списка выражений
	if ws.Expressions == nil {
		return
	}
	//проходим по задачам
	for _, expression := range ws.Expressions.ListExpr {
		// строим дерево выражения
		root, err := parsing.GetTree(expression.Postfix)
		nodes := make([]*parsing.Node, 0)
		nodes = GetNodes(root, nodes)
		// Записываем в выражение ошибку, если она возникла при построении дерева
		// выражения
		if err != nil {
			expression.ParsError = err
			continue
		}
		// Создаем дерево задач
		for _, node := range nodes {
			// Создаем ID для узлов
			node.CreateId()
			// проверить наличие задачи в tasks
			// заполняем словарь узлами
			ws.AllNodes[node.NodeId] = node
			// Если узел не рассчитан и узла с таким ID не в очереди задач
			if node.IsReadyToCalc() && !ws.Tasks.isContent(node) {
				// добавляем его в таски
				ws.Tasks.AddTask(&TaskContainer{
					IdTask:   node.NodeId,
					TaskN:    Task{X: node.X.Val, Y: node.Y.Val, Op: node.Op},
					Deadline: time.Now().Add(time.Hour * 1000),
					TimingsN: expression.Times,
				})
			}
		}
		expression.RootId = root.NodeId
	}
}

// Проходит дерево выражения от корня и создает список узлов выражения
func GetNodes(root *parsing.Node, nodes []*parsing.Node) []*parsing.Node {
	nodes = append(nodes, root)
	if root.Sheet {
		return nodes
	}
	nodes = GetNodes(root.X, nodes)
	nodes = GetNodes(root.Y, nodes)
	return nodes
}

// Проверяет на готовность родительский узел, при готовности добавляет его в очередь задач
func checkAndUpdateNodeToTasks(ws *WorkingSpace, node *parsing.Node) bool {
	node.Mu.RLock()
	defer node.Mu.RUnlock()
	if node.X.Calculated && node.Y.Calculated {
		task := &TaskContainer{
			IdTask: node.NodeId,
			TaskN: Task{
				X:  node.X.Val,
				Y:  node.Y.Val,
				Op: node.Op,
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
