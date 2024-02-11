package parsing

import (
	"fmt"
	"strconv"
)

var priority = map[string]int{
	"+": 1,
	"-": 1,
	"*": 2,
	"/": 2,
}

// PostfixExpr - узел постфиксной записи
type Symbol struct {
	Val string
}

func (s *Symbol) getPriority() int {
	switch s.Val {
	case "+", "-", "*", "/": // если это оператор
		return priority[s.Val]
	case "(", ")":
		return 0
	default:
		return 10
	}
}
func (s *Symbol) getType() string {
	switch s.Val {
	case "+", "-", "*", "/", "(", ")": // если это оператор
		return "Op"
	default:
		return "num"
	}
}
func (s *Symbol) String() string {
	return s.Val
}

// Возвращает узел вычисления полученный из символа
func (s *Symbol) createNode() *Node {
	switch s.getType() {
	case "+", "-", "*", "/": // если символ оператор
		return &Node{Op: s.Val}
	default: // Если символ операнд
		val, _ := strconv.ParseFloat(s.Val, 64)
		return &Node{Val: val, Sheet: true, Calculated: true}
	}
}

// Node - узел выражения
type Node struct {
	NodeId     string  `json:"nodeId"`
	Op         string  `json:"op"` // оператор
	X          *Node   `json:"x"`
	Y          *Node   `json:"y"`     // потомки
	Val        float64 `json:"Val"`   // значение узла
	Sheet      bool    `json:"sheet"` // флаг листа
	Calculated bool    `json:"calculated"`
	// Parent     *Node   `json:"parent"` // узел родитель
}

func (n *Node) getType() string {
	if n.Op != "" {
		return "Op"
	}
	return "num"
}
func (n *Node) String() string {
	if n.Op == "" {
		return fmt.Sprintf("%f", n.Val)
	}
	return fmt.Sprintf("(%s%s%s)", n.X, n.Op, n.Y)
}
func (n *Node) doId() {
	n.NodeId = n.String()
}
func (n *Node) getVal() string {
	if n.Op == "" {
		return fmt.Sprint(n.Val)
	}
	return n.Op
}

type additiveStack interface {
	Node | Symbol
}
type Stack[T additiveStack] struct {
	val []*T
}

func (s *Stack[T]) pop() *T {
	length := len(s.val)
	if !s.isEmpty() {
		res := s.val[length-1]
		s.val = s.val[:length-1]
		return res
	} else {
		return nil
	}
}
func (s *Stack[T]) push(l *T) {
	s.val = append(s.val, l)
}
func (s *Stack[T]) isEmpty() bool {
	if len(s.val) == 0 {
		return true
	}
	return false
}
func (s *Stack[T]) get() *T {
	if s.isEmpty() {
		return nil
	}
	return s.val[len(s.val)-1]
}

func (s *Stack[T]) getPrev() *T {
	if len(s.val) > 1 {
		return s.val[len(s.val)-2]
	}
	return nil
}
