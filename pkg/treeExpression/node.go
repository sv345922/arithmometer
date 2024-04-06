package treeExpression

import (
	"fmt"
	"sync"
)

type Nodes struct {
	AllNodes map[uint64]*Node
	mu       sync.RWMutex
}

// Создать новой список узлов
func NewNodes() *Nodes {
	nodes := make(map[uint64]*Node)
	return &Nodes{AllNodes: nodes}
}

// Если узел есть в списке, возвращает true, иниче false
func (ns *Nodes) Contains(key uint64) bool {
	ns.mu.RLock()
	defer ns.mu.RUnlock()
	if _, ok := ns.AllNodes[key]; ok {
		return true
	}
	return false
}

// Добавить узел в список узлов, если узел с таким id существует в списке, возвращает false
func (ns *Nodes) Add(n *Node) bool {
	key := n.NodeId
	if ns.Contains(key) {
		return false
	}
	ns.mu.Lock()
	ns.AllNodes[key] = n
	ns.mu.Unlock()
	return true
}

// возвращает узел из исписка по его id, если id нет в списке узлов, возвращает nil
func (ns *Nodes) Get(key uint64) *Node {
	if ns.Contains(key) {
		ns.mu.RLock()
		defer ns.mu.RUnlock()
		return ns.AllNodes[key]
	}
	return nil
}

// Удалить узел из списка узлов по его id, если узла с таким id нет, возвращает false
func (ns *Nodes) Remove(key uint64) bool {
	if ns.Contains(key) {
		ns.mu.Lock()
		delete(ns.AllNodes, key)
		return true
	}
	return false
}

// Стрингер
func (ns *Nodes) String() string {
	result := ""
	for key, val := range ns.AllNodes {
		result += fmt.Sprintf("key=%d, val=%s\n",
			key,
			val.getVal())
	}
	return result
}

// Node - узел выражения
type Node struct {
	NodeId     uint64  `json:"nodeId"`
	Op         string  `json:"op"` // оператор
	X          uint64  `json:"x"`
	Y          uint64  `json:"y"`          // потомки
	Val        float64 `json:"val"`        // значение узла
	Sheet      bool    `json:"sheet"`      // флаг листа
	Calculated bool    `json:"calculated"` // флаг вычисленного узла
	Parent     uint64  `json:"parent"`     // узел родитель
	Mu         sync.RWMutex
}

// Создает новый пустой узел
func NewNode() *Node {
	return new(Node)
}

// Возвращает узел X
func (n *Node) GetX(nodes *Nodes) *Node {
	n.Mu.RLock()
	defer n.Mu.RUnlock()
	result := nodes.Get(n.X)
	return result
}

// Возвращает узел Y
func (n *Node) GetY(nodes *Nodes) *Node {
	n.Mu.RLock()
	defer n.Mu.RUnlock()
	result := nodes.Get(n.Y)
	return result
}

// Возвращает узел Parent
func (n *Node) GetParent(nodes *Nodes) *Node {
	n.Mu.RLock()
	defer n.Mu.RUnlock()
	result := nodes.Get(n.Parent)
	return result
}

// Создает ID у узла
func (n *Node) CreateId() uint64 {
	id := NewId(n.String())
	n.Mu.Lock()
	n.NodeId = id
	n.Mu.Unlock()
	return n.NodeId
}

// создает id из строки
func NewId(s string) uint64 {
	res := uint64(0)
	for i, v := range []byte(s) {
		res += uint64(i)
		res += uint64(v)
	}
	return res
}

// проверка на готовность к вычислению
func (n *Node) IsReadyToCalc(nodes *Nodes) bool {
	n.Mu.RLock()
	if !n.Calculated {
		n.Mu.RUnlock()
		if n.GetX(nodes).IsCalculated() && n.GetY(nodes).IsCalculated() {
			return true
		}
	}
	n.Mu.RUnlock()
	return false
}

// Возвращает тип узла
func (n *Node) GetType() string {
	n.Mu.RLock()
	defer n.Mu.RUnlock()
	if n.Op != "" {
		return "Op"
	}
	return "num"
}

// Стрингер
func (n *Node) String() string {
	return fmt.Sprintf("id: %d,\tx_id: %d,\ty_id: %d,\tparent_id: %d,\tval: %.4v\n",
		n.NodeId,
		n.X,
		n.Y,
		n.Parent,
		n.getVal(),
	)
}

// Возвращает значение узла
func (n *Node) getVal() string {
	n.Mu.RLock()
	defer n.Mu.RUnlock()
	if n.Op == "" {
		return fmt.Sprintf("%f", n.Val)
	}
	return n.Op
}

func (n *Node) IsCalculated() bool {
	n.Mu.RLock()
	defer n.Mu.RUnlock()
	return n.Calculated
}
