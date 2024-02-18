package parsing

import (
	"fmt"
	"strconv"
	"sync"
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
	case "Op": // если символ оператор
		return &Node{Op: s.Val}
	default: // Если символ операнд
		val, _ := strconv.ParseFloat(s.Val, 64)
		return &Node{Val: val, Sheet: true, Calculated: true}
	}
}

// Node - узел выражения
type Node struct {
	NodeId     uint64       `json:"nodeId"`
	Op         string       `json:"op"` // оператор
	X          *Node        `json:"x"`
	Y          *Node        `json:"y"`     // потомки
	Val        float64      `json:"Val"`   // значение узла
	Sheet      bool         `json:"sheet"` // флаг листа
	Calculated bool         `json:"calculated"`
	ErrZeroDiv error        `json:"err"`
	Parent     *Node        `json:"parent"` // узел родитель
	Mu         sync.RWMutex `json:"-"`
}

// Создает ID у узла
func (n *Node) CreateId() uint64 {
	n.NodeId = GetId(n.String())
	//n.NodeId = int(time.Now().Unix())
	/*
		s := n.String()
		hasher := sha1.New()
		hasher.Write([]byte(s))
		n.NodeId = base64.URLEncoding.EncodeToString(hasher.Sum(nil))
		return n.NodeId
	*/
	return n.NodeId
}

// проверка на готовность к вычислению
func (n *Node) IsReadyToCalc() bool {
	if !n.Calculated {
		if n.X.Calculated && n.Y.Calculated {
			return true
		}
	}
	return false
}

// создает id
func GetId(s string) uint64 {
	res := uint64(0)
	for i, v := range []byte(s) {
		res += uint64(i)
		res += uint64(v)
	}
	return res
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
	if n.getType() != "Op" {
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
	ind int
}

func newStack[T additiveStack](size int) *Stack[T] {
	return &Stack[T]{
		val: make([]*T, size),
		ind: -1,
	}
}

// Извлечь верхний элемент из стека и удалить его
func (s *Stack[T]) pop() *T {
	if !s.isEmpty() {
		res := s.val[s.ind]
		s.val[s.ind] = nil
		s.ind--
		return res
	} else {
		return nil
	}
}

// Добавить элемент в стек
func (s *Stack[T]) push(l *T) {
	s.ind++
	s.val[s.ind] = l
}

// Проверить стек на пустоту
func (s *Stack[T]) isEmpty() bool {
	return s.ind < 0
}

// Вернуть значение верхнего элемента стека
func (s *Stack[T]) top() *T {
	if s.isEmpty() {
		return nil
	}
	return s.val[s.ind]
}
