package parsing

import (
	"crypto/sha1"
	"encoding/base64"
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
	Err        error   `json:"err"` // узел родитель
}

// Создает ID у узла
func (n *Node) CreateId() string {
	s := n.String()
	hasher := sha1.New()
	hasher.Write([]byte(s))
	n.NodeId = base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return n.NodeId
}

// Возвращает тип узла
func (n *Node) getType() string {
	if n.Op != "" {
		return "Op"
	}
	return "num"
}

// Стрингер
func (n *Node) String() string {
	if n.Op == "" {
		return fmt.Sprintf("%f", n.Val)
	}
	return fmt.Sprintf("(%s%s%s)", n.X, n.Op, n.Y)
}

// Возвращает значение узла
func (n *Node) getVal() string {
	if n.Op == "" {
		return fmt.Sprint(n.Val)
	}
	return n.Op
}

type additiveStack interface {
	Node | Symbol
}

// Стэк для реализации алгоритма Дийксты
type Stack[T additiveStack] struct {
	val []*T
}

// Извлечь верхний элемент из стека и удалить его
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

// Добавить элемент в стек
func (s *Stack[T]) push(l *T) {
	s.val = append(s.val, l)
}

// Проверить стек на пустоту
func (s *Stack[T]) isEmpty() bool {
	if len(s.val) == 0 {
		return true
	}
	return false
}

// Вернуть значение верхнего элемента стека
func (s *Stack[T]) get() *T {
	if s.isEmpty() {
		return nil
	}
	return s.val[len(s.val)-1]
}
