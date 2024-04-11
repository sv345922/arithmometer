package parser

import (
	"arithmometer/internal/entities"
	"arithmometer/internal/stack"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"text/scanner"
)

// var errExpression = errors.New("invalid expression")

var priority = map[string]int{
	"+": 1,
	"-": 1,
	"*": 2,
	"/": 2,
}

// Symbol - содержит символ выражения
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

// Возвращает узел вычисления полученный из символа при построении дерева вычисления
func (s *Symbol) createNode() *Node {
	switch s.getType() {
	case "Op": // если символ оператор
		return &Node{Op: s.Val}
	default: // Если символ операнд
		val, _ := strconv.ParseFloat(s.Val, 64)
		return &Node{Val: val, Sheet: true, Calculated: true}
	}
}

// Parse - парсит выражение в символы
func Parse(input string) ([]*Symbol, error) {
	var s scanner.Scanner
	s.Init(strings.NewReader(input))
	s.Mode = scanner.ScanFloats | scanner.ScanInts // | scanner.ScanIdents
	var SymbList []*Symbol                         // список символов выражения, без пробелов и некорректных символов
	for token := s.Scan(); token != scanner.EOF; token = s.Scan() {
		text := s.TokenText()
		switch token {
		case scanner.Int, scanner.Float:
			SymbList = append(SymbList, &Symbol{text})
		default:
			switch text {
			case "+", "-", "*", "/", "(", ")":
				SymbList = append(SymbList, &Symbol{text})
			default:
				return nil, errors.Join(entities.ErrExpression, fmt.Errorf("invalid expression: %s", text))
			}
		}
	}
	return SymbList, nil
}

// Создает постфиксную запись выражения
func GetPostfix(input []*Symbol) ([]*Symbol, error) {
	var postFix []*Symbol                         // последовательность постфиксного выражения
	opStack := stack.NewStack[Symbol](len(input)) // стек хранения операторов

	for _, currentSymbol := range input {
		switch currentSymbol.getType() {
		case "num":
			postFix = append(postFix, currentSymbol)
		case "Op":
			switch currentSymbol.Val {
			case "(":
				opStack.Push(currentSymbol)
			case ")":
				for {
					headStack := opStack.Pop()
					if headStack == nil {
						return nil, errors.Join(entities.ErrExpression, fmt.Errorf("invalid paranthesis"))
					}
					if headStack.Val != "(" {
						postFix = append(postFix, headStack)
					} else {
						break
					}
				}
			default: // Val оператор
				priorCur := currentSymbol.getPriority()
				for !opStack.IsEmpty() && opStack.Top().getPriority() >= priorCur {
					postFix = append(postFix, opStack.Pop())
				}
				opStack.Push(currentSymbol)
			}
		}
	}
	for !opStack.IsEmpty() {
		postFix = append(postFix, opStack.Pop())
	}
	return postFix, nil
}

// GetTree Строит дерево выражения и возвращает корневой узел и список узлов
// из постфиксного выражения
func GetTree(postfix []*Symbol) (*entities.Node, []*entities.Node, error) {
	if len(postfix) == 0 {
		return nil, nil, errors.Join(entities.ErrExpression, fmt.Errorf("expression is empty"))
	}
	stack := stack.NewStack[Node](len(postfix))
	for _, symbol := range postfix {
		node := symbol.createNode()
		// Если узел оператор

		if node.GetType() != "num" {
			// если стек пустой, возвращаем ошибку выражения
			if stack.IsEmpty() {
				return nil, nil, errors.Join(entities.ErrExpression, fmt.Errorf("оператор без операнда"))
			}
			y := stack.Pop() // взять
			x := stack.Pop() // взять

			// если в стеке нет x, создаем вместо него узел с val=0,
			// обработка унарных операторов
			if x == nil {
				node.X = &Node{Val: 0, Parent: node, Sheet: true, Calculated: true}
				node.Y = y
				// устанавливаем родителя
				y.Parent = node
			} else {
				node.X = x
				node.Y = y
				// устанавливаем родителя
				x.Parent = node
				y.Parent = node
			}
			stack.Push(node) // положить
		} else {
			// если узел не оператор, то он число
			stack.Push(node) // положить
		}
	}
	// получаем список узлов выражения
	root := stack.Top()
	nodes := make([]*Node, 0)
	GetNodesOfExpression(root, &nodes)
	rootOut := TransformNode(root)
	nodesOut := make([]*entities.Node, len(nodes))
	for i, value := range nodes {
		nodesOut[i] = TransformNode(value)
	}
	return rootOut, nodesOut, nil
}

// Проходит дерево выражения от корня и создает список узлов выражения
func GetNodesOfExpression(node *Node, nodes *[]*Node) {
	node.CreateId()
	//fmt.Printf("created ID=%d for node: %v\n", node.Id, node)
	*nodes = append(*nodes, node)
	if node.Sheet {
		return
	}
	GetNodesOfExpression(node.X, nodes)
	GetNodesOfExpression(node.Y, nodes)
	return
}

// преобразует ноды парсера в ноды общие
func TransformNode(node *Node) *entities.Node {
	var x, y, parent uint64
	if node.X != nil {
		x = node.X.NodeId
	}
	if node.Y != nil {
		y = node.Y.NodeId
	}
	if node.Parent != nil {
		parent = node.Parent.NodeId
	}
	result := &entities.Node{
		Id:         node.NodeId,
		Op:         node.Op,
		X:          x,
		Y:          y,
		Val:        node.Val,
		Sheet:      node.Sheet,
		Calculated: node.Calculated,
		Parent:     parent,
	}
	return result
}
