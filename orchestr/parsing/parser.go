package parsing

import (
	"fmt"
	//"strconv"
	"strings"
	"text/scanner"
)

// Parse - парсит строку в дерево, возвращает список узлов, корень дерева и ошибку
func Parse(input string) ([]*Symbol, []*Node, *Node, error) {
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
				return nil, nil, nil, fmt.Errorf("invalid expression: %s", text)
			}
		}
	}
	postFix, err := getPostfix(SymbList)
	if err != nil {
		return nil, nil, nil, err
	}
	nodeList, root, err := getTree(postFix)
	if err != nil {
		return nil, nil, nil, err
	}
	return postFix, nodeList, root, nil
}

// Создает постфиксную запись выражения
func getPostfix(input []*Symbol) ([]*Symbol, error) {
	var opStack Stack[Symbol]
	var postFix []*Symbol

	for _, val := range input {
		switch val.getType() {
		case "num":
			postFix = append(postFix, val)
		case "Op":
			switch val.Val {
			case "(":
				opStack.push(val)
			case ")":
				for n := opStack.pop(); n.Val != "("; n = opStack.pop() {
					if n == nil {
						return nil, fmt.Errorf("invalid paranthesis")
					}
					postFix = append(postFix, val)
				}
			default: // Val оператор
				priorCur := val.getPriority()
				for !opStack.isEmpty() && opStack.get().getPriority() >= priorCur {
					postFix = append(postFix, opStack.pop())
				}
				opStack.push(val)
			}
		}
	}
	for !opStack.isEmpty() {
		postFix = append(postFix, opStack.pop())
	}
	return postFix, nil
}

// Строит дерево выражения и возвращает список узлов и корневой узел
func getTree(postfix []*Symbol) ([]*Node, *Node, error) {
	if len(postfix) == 0 {
		return nil, nil, fmt.Errorf("expression is empty")
	}
	stack := new(Stack[Node])
	for _, symb := range postfix {
		node := symb.createNode()
		if node.getType() == "Op" {
			y := stack.pop()
			x := stack.get()
			// при первом отрицательном числе
			if x == nil {
				node = y
				node.Val = -node.Val
			} else {
				node.X = stack.pop()
				node.Y = y
			}
			stack.push(node)
		} else {
			stack.push(node)
		}
	}
	return stack.val, stack.get(), nil
}
