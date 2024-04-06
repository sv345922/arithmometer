package parser

import (
	"arithmometer/pkg/stack"
	"arithmometer/pkg/treeExpression"
	"fmt"
	"strconv"
	"strings"
	"text/scanner"
)

// TODO обновить связи

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
func (s *Symbol) createNode() *treeExpression.Node {
	switch s.getType() {
	case "Op": // если символ оператор
		return &treeExpression.Node{Op: s.Val}
	default: // Если символ операнд
		val, _ := strconv.ParseFloat(s.Val, 64)
		return &treeExpression.Node{Val: val, Sheet: true, Calculated: true}
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
				return nil, fmt.Errorf("invalid expression: %s", text)
			}
		}
	}
	return getPostfix(SymbList)
}

// Создает постфиксную запись выражения
func getPostfix(input []*Symbol) ([]*Symbol, error) {
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
						return nil, fmt.Errorf("invalid paranthesis")
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

// GetTree Строит дерево выражения и возвращает корневой узел из постфиксного выражения
func GetTree(postfix []*Symbol) (*treeExpression.Node, *[]*treeExpression.Node, error) {
	if len(postfix) == 0 {
		return nil, nil, fmt.Errorf("expression is empty")
	}
	stack := stack.NewStack[treeExpression.Node](len(postfix))
	for _, symbol := range postfix {
		node := symbol.createNode()
		// Если узел оператор

		if node.GetType() != "num" {
			// если стек пустой, возвращаем ошибку выражения
			if stack.IsEmpty() {
				return nil, nil, fmt.Errorf("ошибка выражения, оператор без операнда")
			}
			y := stack.Pop() // взять
			x := stack.Pop() // взять

			// если в стеке нет x, создаем вместо него узел с val=0,
			// обработка унарных операторов
			if x == nil {
				node.X = &treeExpression.Node{Val: 0, Parent: node, Sheet: true, Calculated: true}
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
	nodes := make([]*treeExpression.Node, 0)
	GetNodes(root, &nodes)
	return root, &nodes, nil
}

// Проходит дерево выражения от корня и создает список узлов выражения
func GetNodes(node *treeExpression.Node, nodes *[]*treeExpression.Node) {
	node.CreateId()
	*nodes = append(*nodes, node)
	if node.Sheet {
		return
	}
	GetNodes(node.X, nodes)
	GetNodes(node.Y, nodes)
	return
}
