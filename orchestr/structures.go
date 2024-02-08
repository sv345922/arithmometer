package orchestr

import "fmt"

var priority map[string]int = map[string]int{
	"+": 1,
	"-": 1,
	"*": 2,
	"/": 2,
}

// Узел
type Node struct {
	op     string  // оператор
	x, y   *Node   // потомки
	val    float64 // значение узла
	sheet  bool    // флаг листа
	parent *Node   // узел родитель
}

func (n *Node) getType() string {
	if n.op != "" {
		return "op"
	}
	return "num"
}
func (n *Node) getPriority() int {

	if n.op != "" {
		return priority[n.op]
	}
	return 10
}
func (n *Node) String() string {
	if n.sheet {
		return fmt.Sprintf("%f", n.val)
	}
	return fmt.Sprintf("(%s%s%s)", n.x, n.op, n.y)
	/*
		if n.op != "" {
			return fmt.Sprintf("%s", n.op)
		}
		return fmt.Sprintf("%f", n.val)
	*/
}

type stack struct {
	val []*Node
}

func (s *stack) pop() *Node {
	length := len(s.val)
	if !s.isEmpty() {
		res := s.val[length-1]
		s.val = s.val[:length-1]
		return res
	} else {
		return nil
	}
}
func (s *stack) push(l *Node) {
	s.val = append(s.val, l)
}
func (s *stack) isEmpty() bool {
	if len(s.val) == 0 {
		return true
	}
	return false
}
func (s *stack) get() *Node {
	if s.isEmpty() {
		return nil
	}
	return s.val[len(s.val)-1]
}

func (s *stack) getPrev() *Node {
	if len(s.val) > 1 {
		return s.val[len(s.val)-2]
	}
	return nil
}
