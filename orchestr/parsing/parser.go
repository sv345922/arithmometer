package parsing

import (
	"fmt"
	//"strconv"
	"strings"
	"text/scanner"
)

// Parse - парсит строку в дерево, возвращает список узлов, корень дерева и ошибку
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
	postFix, err := getPostfix(SymbList)
	if err != nil {
		return nil, err
	}
	return postFix, nil
}

// Создает постфиксную запись выражения
func getPostfix(input []*Symbol) ([]*Symbol, error) {
	var postFix []*Symbol                   // последовательсность постфиксного выражения
	opStack := newStack[Symbol](len(input)) // стек хранения операторов

	for _, currentSymbol := range input {
		switch currentSymbol.getType() {
		case "num":
			postFix = append(postFix, currentSymbol)
		case "Op":
			switch currentSymbol.Val {
			case "(":
				opStack.push(currentSymbol)
			case ")":
				for {
					headStack := opStack.pop()
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
				for !opStack.isEmpty() && opStack.top().getPriority() >= priorCur {
					postFix = append(postFix, opStack.pop())
				}
				opStack.push(currentSymbol)
			}
		}
	}
	for !opStack.isEmpty() {
		postFix = append(postFix, opStack.pop())
	}
	return postFix, nil
}

// Строит дерево выражения и возвращает корневой узел из постфиксного выражения
func GetTree(postfix []*Symbol) (*Node, error) {
	if len(postfix) == 0 {
		return nil, fmt.Errorf("expression is empty")
	}
	stack := newStack[Node](len(postfix))
	for _, symbol := range postfix {
		node := symbol.createNode()
		// Если узел оператор

		if node.getType() != "num" {
			// если стек пустой, возвращаем ошибку выражения
			if stack.isEmpty() {
				return nil, fmt.Errorf("ошибка выражения, оператор без операнда")
			}
			y := stack.pop() // взять
			x := stack.pop() // взять

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
			stack.push(node) // положить
		} else {
			// если узел не оператор, то он число
			stack.push(node) // положить
		}
	}

	return stack.top(), nil
}
