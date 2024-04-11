package parser

import (
	"fmt"
)

// Node - узел выражения
type Node struct {
	NodeId     uint64  `json:"nodeId"`
	Op         string  `json:"op"` // оператор
	X          *Node   `json:"x"`
	Y          *Node   `json:"y"`          // потомки
	Val        float64 `json:"val"`        // значение узла
	Sheet      bool    `json:"sheet"`      // флаг листа
	Calculated bool    `json:"calculated"` // флаг вычисленного узла
	Parent     *Node   `json:"parent"`     // узел родитель
}

// Возвращает тип узла
func (n *Node) GetType() string {
	if n.Op != "" {
		return "Op"
	}
	return "num"
}

// Создает ID у узла
func (n *Node) CreateId() uint64 {
	id := NewId(n.String())
	n.NodeId = id
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

// Стрингер
func (n *Node) String() string {
	if n == nil {
		return ""
	}
	if n.Sheet {
		return fmt.Sprintf("%.4f", n.Val)
	}
	return fmt.Sprintf("%v%v%v", n.X.String(), n.Op, n.Y.String())
	//return fmt.Sprintf("id: %d,\tx_id: %v,\ty_id: %v,\tparent_id: %v,\tval: %s%.4f\n",
	//	n.Id,
	//	n.X,
	//	n.Y,
	//	n.Parent,
	//	n.Op,
	//	n.Val,
	//)
}
