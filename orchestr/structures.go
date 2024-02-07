package orchestr

import "fmt"

// Узел
type Node struct {
	op     string  // оператор
	x, y   *Node   // потомки
	val    float64 // значение узла
	sheet  bool    // флаг листа
	parent *Node   // узел родитель
	level  int	   // уровень вложенности
}

func (n *Node) String() string {
	if n.sheet {
		return fmt.Sprintf("%f", n.val)
	}
	return fmt.Sprintf("(%s%s%s)", n.x, n.op, n.y)
}

type Tree struct {
	nodes     []*Node
	relations map[int][]int
}

type relNode struct {
	x, y *Node
}
