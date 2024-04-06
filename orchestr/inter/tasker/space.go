package tasker

import (
	"arithmometer/pkg/expressions"
	"arithmometer/pkg/parser"
	"arithmometer/pkg/taskQueue"
	"arithmometer/pkg/timings"
	"arithmometer/pkg/treeExpression"
	"fmt"
	"log"
	"sync"
	"time"
)

type WorkingSpace struct {
	Tasks       *taskQueue.Tasks                 `json:"tasks"`
	Expressions *expressions.Expressions         `json:"expressions"`
	Timings     *timings.Timings                 `json:"timings"`
	AllNodes    *map[uint64]*treeExpression.Node `json:"allNodes"` // ключ id узла
	mu          sync.RWMutex
}

// Сохраняет рабочее пространство
func (ws *WorkingSpace) Save() error {
	return SaveWS(ws)
}

// При получении выполненого задания,
// проверяет на наличие ошибки деления на ноль,
// Записывает результат в узел и изменяет статус на - вычислено
// Обновляет очередь задач.
// Проверяет список выражений и если оно вычислено, обновляет его статус.
// Добавляет новую задачу в начало очереди задач.
func (ws *WorkingSpace) UpdateTasks(IdTask uint64, answer *Answer) error {
	ws.mu.RLock()
	// находим узел решенной задаче
	calculatedNode, ok := (*ws.AllNodes)[IdTask]
	ws.mu.RUnlock()
	if !ok {
		return fmt.Errorf("узел в мапе активных узлов не найден")
	}
	// Проверка деления на ноль и обновление выражения
	// с удалением не требующих решения задач,
	// а также изменение статуса выражения
	if answer.Err != nil {
		ws.updateWhileZeroDiv(calculatedNode, answer.Err)
		return answer.Err
	}
	result := answer.Result
	// Удаляем задачу из очереди
	ws.Tasks.RemoveTask(IdTask)
	// записываем результат вычисления в узел
	calculatedNode.Val = result
	calculatedNode.Calculated = true

	// Проверяем родительский узел
	parent := calculatedNode.Parent
	// Если это корень выражения
	if parent == nil {
		// Обновляем результат выражения и его статус
		ws.Expressions.UpdateStatus(calculatedNode, "done", result)
		return nil
	}
	// проверка готовности родительского узла и добавление его в очередь задач
	for checkAndUpdateNodeToTasks(ws, parent) {
		if parent.Parent == nil {
			break
		}
		parent = parent.Parent
	}
	return nil
}

// Добавляет новое выражение в структура,
// обновляет мапу узлов
// обновляет очередь вычислений
func (ws *WorkingSpace) AddExpression(expression *expressions.Expression) error {
	// построить дерево выражения и внести корень
	root, nodes, err := parser.GetTree(expression.Postfix)
	// создать id корневого узла
	root.CreateId()
	expression.RootId = root.NodeId

	// добавляем выражение в список выражений
	ws.Expressions.Add(expression)
	// Если выражение не может быть построено, возращаем ошибку
	if err != nil {
		return fmt.Errorf("обшибка построения выражения: %v", err)
	}

	// проходим по нему и добавляем узлы готовые к вычислению в очередь
	// сами узлы добавляем в AllNodes
	for _, node := range *nodes {
		// Создаем ID для узлов
		if node.NodeId == 0 {
			node.CreateId()
		}
		// заполняем словарь узлами
		ws.mu.Lock()
		(*ws.AllNodes)[node.NodeId] = node
		// Если узел не рассчитан и узла с таким ID нет в очереди задач
		if node.IsReadyToCalc() {
			// добавляем его в таски
			ws.Tasks.AddTask(&TaskContainer{
				IdTask:   node.NodeId,
				TaskN:    Task{X: node.X.Val, Y: node.Y.Val, Op: node.Op},
				Deadline: time.Now().Add(time.Hour * 1000),
				TimingsN: expression.Times,
			})
		}
		ws.mu.Unlock()
	}
	return nil
}

// TODO - не используется
// При поступлении нового выражения
// проходит по списку выражений, создает дерево узлов выражения,
// включает в рабочее пространство список узлов - ws.AllNodes
// созадет очередь задач для вычислителей - ws.tasks
func (ws *WorkingSpace) Update() {
	//ws.mu.Lock()
	//defer ws.mu.Unlock()
	//// Взять выражения
	//// проверить на существование списка выражений
	//if ws.Expressions == nil {
	//	return
	//}
	////проходим по задачам
	//for _, expression := range ws.Expressions.ListExpr {
	//	// строим дерево выражения
	//	root, nodes, err := parsing.GetTree(expression.Postfix)
	//
	//	// Записываем в выражение ошибку, если она возникла при построении дерева
	//	// выражения
	//	if err != nil {
	//		expression.ParsError = err
	//		continue
	//	}
	//	// Создаем дерево задач
	//	for _, node := range *nodes {
	//		// Создаем ID для узлов
	//		node.CreateId()
	//		// проверить наличие задачи в tasks
	//		// заполняем словарь узлами
	//		(*ws.AllNodes)[node.NodeId] = node
	//		// Если узел не рассчитан и узла с таким ID нет в очереди задач
	//		if node.IsReadyToCalc() && !ws.Tasks.isContent(node) {
	//			// добавляем его в таски
	//			ws.Tasks.AddTask(&TaskContainer{
	//				IdTask:   node.NodeId,
	//				TaskAn:    Task{X: node.X.Val, Y: node.Y.Val, Op: node.Op},
	//				Deadline: time.Now().Add(time.Hour * 1000),
	//				TimingsN: expression.Times,
	//			})
	//		}
	//	}
	//	expression.RootId = root.NodeId
	//}
}

// Проходит дерево выражения от корня и создает список узлов выражения - удалить
//func GetNodes(root *parsing.NodeDB, nodes *[]*parsing.NodeDB) []*parsing.NodeDB {
//	nodes = append(nodes, root)
//	if root.Sheet {
//		return nodes
//	}
//	nodes = GetNodes(root.X, nodes)
//	nodes = GetNodes(root.Y, nodes)
//	return nodes
//}

// Проверяет на готовность узел, при готовности добавляет его в очередь задач
// TODO проверить, возможно отсюда идет ошибка очереди
func checkAndUpdateNodeToTasks(ws *WorkingSpace, node *treeExpression.Node) bool {
	// Если x и y вычислены
	if node.X.IsCalculated() && node.Y.IsCalculated() {
		// создаем задачу и кладем её в очередь
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
// проверяет узлы в дереве выражения и обновляет их
func (ws *WorkingSpace) updateWhileZeroDiv(node *treeExpression.Node, err error) {
	log.Println("в выражении присутствует деление на ноль")
	err = fmt.Errorf(err.Error() + "in Expression")
	// находим кореневой узел выражения
	root := node.Parent
	for ; root != nil; root = node.Parent {
	}
	// Изменяем статус выражения с ошибкой
	ws.Expressions.UpdateStatus(root, err.Error(), 0)

	//Удаляем узлы выражения из очереди и мапы узлов
	ws.removeCalculatedNodes(root)
}

// Удаляем узлы выражения из очереди и мапы узлов по корневому узлу
func (ws *WorkingSpace) removeCalculatedNodes(node *treeExpression.Node) {
	ws.mu.RLock()
	defer ws.mu.RUnlock()
	for node.X != nil {
		ws.removeCalculatedNodes(node.X)
	}
	for node.Y != nil {
		ws.removeCalculatedNodes(node.Y)
	}
	ws.mu.Lock()
	delete((*ws.AllNodes), node.NodeId)
	ws.mu.Unlock()
	ws.Tasks.RemoveTask(node.NodeId)
}
