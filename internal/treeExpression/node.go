package treeExpression

import (
	"arithmometer/internal/entities"
	"fmt"
	"sync"
)

// Node - узел выражения
type Node struct {
	entities.Node
	mu sync.RWMutex
}

// Создает новый пустой узел
func NewNode(node *entities.Node) *Node {
	result := new(Node)
	result.Node = *node
	return result
}

// Возвращает узел X
func (n *Node) GetX(nodes *Nodes) *Node {
	n.mu.RLock()
	defer n.mu.RUnlock()
	result := nodes.Get(n.X)
	return result
}

// Возвращает узел Y
func (n *Node) GetY(nodes *Nodes) *Node {
	n.mu.RLock()
	defer n.mu.RUnlock()
	result := nodes.Get(n.Y)
	return result
}

// Возвращает узел Parent
func (n *Node) GetParent(nodes *Nodes) *Node {
	n.mu.RLock()
	defer n.mu.RUnlock()
	result := nodes.Get(n.Parent)
	return result
}

// Создает ID у узла
func (n *Node) CreateId() uint64 {
	id := NewId(n.String())
	n.mu.Lock()
	n.Id = id
	n.mu.Unlock()
	return n.Id
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
	n.mu.RLock()
	if !n.Calculated {
		n.mu.RUnlock()
		if n.GetX(nodes).IsCalculated() && n.GetY(nodes).IsCalculated() {
			return true
		}
	} else {
		n.mu.RUnlock()
	}
	return false
}

// Возвращает тип узла
func (n *Node) GetType() string {
	n.mu.RLock()
	defer n.mu.RUnlock()
	if n.Op != "" {
		return "Op"
	}
	return "num"
}

// Стрингер
func (n *Node) String() string {
	return fmt.Sprintf("id: %d,\tx_id: %d,\ty_id: %d,\tparent_id: %d,\tval: %.4v\n",
		n.Id,
		n.X,
		n.Y,
		n.Parent,
		n.getVal(),
	)
}

// Возвращает значение узла
func (n *Node) getVal() string {
	n.mu.RLock()
	defer n.mu.RUnlock()
	if n.Op == "" {
		return fmt.Sprintf("%f", n.Val)
	}
	return n.Op
}

func (n *Node) IsCalculated() bool {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.Calculated
}
