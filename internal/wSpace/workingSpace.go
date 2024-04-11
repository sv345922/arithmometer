package wSpace

import (
	"arithmometer/internal/dataBase"
	"arithmometer/internal/entities"
	"arithmometer/internal/taskQueue"
	"fmt"
	"log"
	"sync"
)

type WorkingSpace struct {
	Queue       *taskQueue.Queue                `json:"tasks"`
	Expressions map[uint64]*entities.Expression `json:"expressions"`
	AllNodes    map[uint64]*entities.Node       `json:"allnodes"`
	Timings     *entities.Timings               `json:"timings"`
	Mu          sync.RWMutex
}

func NewWorkingSpace() *WorkingSpace {
	ws := new(WorkingSpace)
	ws.Queue = taskQueue.NewQueue()
	expressions := make(map[uint64]*entities.Expression)
	ws.Expressions = expressions
	allNodes := make(map[uint64]*entities.Node)
	ws.AllNodes = allNodes
	return ws
}

// Сохраняет рабочее пространство
func (ws *WorkingSpace) Save() error {
	ws.Mu.RLock()
	// Создаем и заполняем структуру базы данных для сохранения
	db := dataBase.NewDB()
	// Заполняем структуру БД
	// заполняем тайминги
	db.Timings = ws.Timings
	// заполняем очередь задач
	db.Queue = ws.Queue
	// заполняем список выражений
	db.Expressions = ws.Expressions
	// заполняем AllNodes
	db.AllNodes = ws.AllNodes

	ws.Mu.RUnlock()
	err := db.Save()
	if err != nil {
		return err
	}
	log.Println("DB saved")
	return nil
}

// Добавляет новое выражение в структуры,
// обновляет мапу узлов
// обновляет очередь вычислений
func (ws *WorkingSpace) AddToExpressions(expression *entities.Expression) error {
	// добавляем выражение в список выражений
	id := expression.Id
	// Если выражение с таким id есть возвращаем ошибку
	if _, ok := ws.Expressions[id]; ok {
		return fmt.Errorf("expression already exists")
	}
	// Добавляем выражение в мапу выражений
	ws.Expressions[expression.Id] = expression

	return nil
}
func (ws *WorkingSpace) AddToAllNodes(nodes []*entities.Node) (int, error) {
	n := 0
	for _, node := range nodes {
		if _, ok := ws.AllNodes[node.Id]; !ok {
			ws.AllNodes[node.Id] = node
			n++
		}
	}
	if n != len(nodes) {
		return n, fmt.Errorf("some nodes are exist")
	}
	return n, nil
}

// Возвращает корень узла и выражение с этим корнем
func (ws *WorkingSpace) GetRoot(nodeId uint64) (*entities.Node, *entities.Expression) {
	root, ok := ws.AllNodes[nodeId]
	for ; !ok; root, ok = ws.AllNodes[root.Parent] {
	}
	var expression *entities.Expression
	for _, expression = range ws.Expressions {
		if expression.RootId == root.Id {
			break
		}
	}
	return root, expression
}

// Возвращает список id узлов выражения по id корня
func (ws *WorkingSpace) GetExpressionNodesID(nodeId uint64, nodesId []uint64) {
	nodesId = append(nodesId, nodeId)
	node, ok := ws.AllNodes[nodeId]
	if node.Sheet || !ok {
		return
	}

	ws.GetExpressionNodesID(node.X, nodesId)
	ws.GetExpressionNodesID(node.Y, nodesId)
}

// Загружает структуру из db и возвращает её
func LoadDB() (*WorkingSpace, error) {
	db := dataBase.NewDB()
	err := db.Load()
	if err != nil {
		log.Println(err)
	}

	// Создаем рабочее пространство
	ws := NewWorkingSpace()
	ws.Mu.Lock()
	// заполняем поля
	ws.Queue = db.Queue
	ws.Expressions = db.Expressions
	ws.AllNodes = db.AllNodes
	ws.Timings = db.Timings
	ws.Mu.Unlock()
	return ws, nil
}

// TODO Видимо надо удалить остальное

//// TODO
//// При получении выполненого задания,
//// проверяет на наличие ошибки деления на ноль,
//// Записывает результат в узел и изменяет статус на - вычислено
//// Обновляет очередь задач.
//// Проверяет список выражений и если оно вычислено, обновляет его статус.
//// Добавляет новую задачу в начало очереди задач.
//func (ws *WorkingSpace) UpdateTasks(IdTask uint64, answer *Answer) error {
//	ws.Mu.RLock()
//	// находим узел решенной задаче
//	calculatedNode, ok := (*ws.AllNodes)[IdTask]
//	ws.Mu.RUnlock()
//	if !ok {
//		return fmt.Errorf("узел в мапе активных узлов не найден")
//	}
//	// Проверка деления на ноль и обновление выражения
//	// с удалением не требующих решения задач,
//	// а также изменение статуса выражения
//	if answer.Err != nil {
//		ws.updateWhileZeroDiv(calculatedNode, answer.Err)
//		return answer.Err
//	}
//	result := answer.Result
//	// Удаляем задачу из очереди
//	ws.Queue.RemoveTask(IdTask)
//	// записываем результат вычисления в узел
//	calculatedNode.Val = result
//	calculatedNode.Calculated = true
//
//	// Проверяем родительский узел
//	parent := calculatedNode.Parent
//	// Если это корень выражения
//	if parent == nil {
//		// Обновляем результат выражения и его статус
//		ws.Expressions.UpdateStatus(calculatedNode, "done", result)
//		return nil
//	}
//	// проверка готовности родительского узла и добавление его в очередь задач
//	for checkAndUpdateNodeToTasks(ws, parent) {
//		if parent.Parent == nil {
//			break
//		}
//		parent = parent.Parent
//	}
//	return nil
//}
//
//// TODO - не используется
//// При поступлении нового выражения
//// проходит по списку выражений, создает дерево узлов выражения,
//// включает в рабочее пространство список узлов - ws.AllTask
//// созадет очередь задач для вычислителей - ws.tasks
//func (ws *WorkingSpace) Update() {
//	//ws.mu.Lock()
//	//defer ws.mu.Unlock()
//	//// Взять выражения
//	//// проверить на существование списка выражений
//	//if ws.Expressions == nil {
//	//	return
//	//}
//	////проходим по задачам
//	//for _, expression := range ws.Expressions.ListExpr {
//	//	// строим дерево выражения
//	//	root, nodes, err := parsing.GetTree(expression.Postfix)
//	//
//	//	// Записываем в выражение ошибку, если она возникла при построении дерева
//	//	// выражения
//	//	if err != nil {
//	//		expression.ParsError = err
//	//		continue
//	//	}
//	//	// Создаем дерево задач
//	//	for _, node := range *nodes {
//	//		// Создаем ID для узлов
//	//		node.CreateId()
//	//		// проверить наличие задачи в tasks
//	//		// заполняем словарь узлами
//	//		(*ws.AllTask)[node.Id] = node
//	//		// Если узел не рассчитан и узла с таким ID нет в очереди задач
//	//		if node.IsReadyToCalc() && !ws.Queue.isContent(node) {
//	//			// добавляем его в таски
//	//			ws.Queue.AddTask(&TaskContainer{
//	//				IdTask:   node.Id,
//	//				TaskAn:    Task{X: node.X.Val, Y: node.Y.Val, Op: node.Op},
//	//				Deadline: time.Now().Add(time.Hour * 1000),
//	//				TimingsN: expression.Times,
//	//			})
//	//		}
//	//	}
//	//	expression.RootId = root.Id
//	//}
//}
//
//// Проходит дерево выражения от корня и создает список узлов выражения - удалить
////func GetNodes(root *parsing.NodeDB, nodes *[]*parsing.NodeDB) []*parsing.NodeDB {
////	nodes = append(nodes, root)
////	if root.Sheet {
////		return nodes
////	}
////	nodes = GetNodes(root.X, nodes)
////	nodes = GetNodes(root.Y, nodes)
////	return nodes
////}
//
//// Проверяет на готовность узел, при готовности добавляет его в очередь задач
//// TODO проверить, возможно отсюда идет ошибка очереди
//func checkAndUpdateNodeToTasks(ws *WorkingSpace, node *treeExpression.Node) bool {
//	// Если x и y вычислены
//	if node.X.IsCalculated() && node.Y.IsCalculated() {
//		// создаем задачу и кладем её в очередь
//		task := &wSpace.TaskContainer{
//			IdTask: node.Id,
//			TaskAn: wSpace.Task{
//				X:  node.X.Val,
//				Y:  node.Y.Val,
//				Op: node.Op,
//			},
//			Err:      nil,
//			TimingsN: *ws.Timings,
//		}
//		ws.Queue.AddTask(task)
//		return true
//	}
//	return false
//}
//
//// Обновляет рабочее пространство при обнаружении деления на ноль,
//// проверяет узлы в дереве выражения и обновляет их
//func (ws *WorkingSpace) updateWhileZeroDiv(node *treeExpression.Node, err error) {
//	log.Println("в выражении присутствует деление на ноль")
//	err = fmt.Errorf(err.Error() + "in Expression")
//	// находим кореневой узел выражения
//	root := node.Parent
//	for ; root != nil; root = node.Parent {
//	}
//	// Изменяем статус выражения с ошибкой
//	ws.Expressions.UpdateStatus(root, err.Error(), 0)
//
//	//Удаляем узлы выражения из очереди и мапы узлов
//	ws.removeCalculatedNodes(root)
//}
//
//// Удаляем узлы выражения из очереди и мапы узлов по корневому узлу
//func (ws *WorkingSpace) removeCalculatedNodes(node *treeExpression.Node) {
//	ws.Mu.RLock()
//	defer ws.Mu.RUnlock()
//	for node.X != nil {
//		ws.removeCalculatedNodes(node.X)
//	}
//	for node.Y != nil {
//		ws.removeCalculatedNodes(node.Y)
//	}
//	ws.Mu.Lock()
//	delete(*ws.AllNodes, node.Id)
//	ws.Mu.Unlock()
//	ws.Queue.RemoveTask(node.Id)
//}
//
