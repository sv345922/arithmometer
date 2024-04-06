package treeExpression

import (
	"fmt"
	"sync"
)

// Node - узел выражения
type Node struct {
	NodeId     uint64
	Op         string // оператор
	X          *Node
	Y          *Node   // потомки
	Val        float64 // значение узла
	Sheet      bool    // флаг листа
	Calculated bool    // флаг вычисленного узла
	Parent     *Node   // узел родитель
	Mu         sync.RWMutex
}

// Создает ID у узла
func (n *Node) CreateId() uint64 {
	n.Mu.Lock()
	defer n.Mu.Unlock()
	n.NodeId = GetId(n.String())
	return n.NodeId
}

// создает id из строки
func GetId(s string) uint64 {
	res := uint64(0)
	for i, v := range []byte(s) {
		res += uint64(i)
		res += uint64(v)
	}
	return res
}

// проверка на готовность к вычислению
func (n *Node) IsReadyToCalc() bool {
	if !n.Calculated && n.X.Calculated && n.Y.Calculated {
		return true
	}
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
	if n.GetType() != "Op" {
		return fmt.Sprintf("%f", n.Val)
	}
	return fmt.Sprintf("(%s%s%s)", n.X, n.Op, n.Y)
}

// Возвращает значение узла
func (n *Node) getVal() string {
	n.Mu.RLock()
	defer n.Mu.RUnlock()
	if n.Op == "" {
		return fmt.Sprint(n.Val)
	}
	return n.Op
}

func (n *Node) IsCalculated() bool {
	n.Mu.RLock()
	defer n.Mu.RUnlock()
	return n.Calculated
}
