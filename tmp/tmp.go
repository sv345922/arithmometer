for token := s.Scan(); token != scanner.EOF; token = s.Scan() {
	text := s.TokenText()
	switch token {
	case scanner.Int, scanner.Float:
		// игнорируем ошибку, так как в данном кейсе в строке число
		num, _ := strconv.ParseFloat(text, 64)
		// Если первый элемент узла nil, записываем туда текущее число
		if node.x == nil {
			node.x = &Node{val: num, sheet: true, level: node.level+1}
			// Иначе число записываем во второй элемент, и отмечаем его посчитанным (листом дерева)
		} else if node.y == nil {
			node.y = &Node{val: num, sheet: true, level: node.level+1}
		// если x и y узла не nil значит ошибка в выражении - два числа подряд
		} else {
			return node, fmt.Errorf("expression error, doble numbers whithout operator")
		}
	default:
		// проверка на допустимые операторы
		switch text {
		case "+", "-":
			// Если узел заполнен, оператор заполнен всегда, если node.y != nil
			if node.x == nil {

			}
			if node.x != nil && node.y != nil {
				switch node.op {
				case "+", "-":
					node = getParentX(node)
				case "*", "/":
					return node, nil
				}
			}
			node.op = text
		case "*", "/":
			// y-элемент вычисляем рекурсивно
			// создаем y-элемент
			yEl := &Node{level: (node.level+1)}
			// ставим его на место y-элемента
			node.y = yEl.x
			// устанавливаем оператор
			yEl.op = text
			// вычисляем y-элемент рекурсивно
			_, err := GetTree(yEl, s)
			if err != nil {
				return nil, err
			}
			node.y = yEl
		case "(":
			xEl := &Node{level: (node.level+1)}
			node, err := GetTree(node, s)
			if err != nil {
				return nil, err
			}
			node.x = xEl
		case ")":
			return node, nil
		default:
			return nil, fmt.Errorf("invalid char %s", text)
		}
	}
}

func GetTree(node *Node, s *scanner.Scanner) (*Node, error) {
	//node := &Node{op: "", x: nil, y: nil, val: 0, sheet: false}
	token := s.Scan()
	if token == scanner.EOF {
		return node, nil
	}
	text := s.TokenText()
	switch token {
	case scanner.Int, scanner.Float:
		num, _ := strconv.ParseFloat(text, 64)
		switch {
		case node.x == nil:
			node.x = &Node{val: num, sheet: true, level: node.level+1, parent: node}
		case node.y == nil:
			node.y = &Node{val: num, sheet: true, level: node.level+1, parent: node}
			node.x 
		
		default:
			return nil, fmt.Errorf("invalid expression (или ошибка компоновки узлов) num = %f", num)
		}
	default:
		switch text {
		case "+", "-":
			// если это первый элемент выражения
			if node.x == nil{
				op := text
				token = s.Scan()
				text := s.TokenText()
				switch token {
				case scanner.Int, scanner.Float:
					// Не учитываем ошибку так как в кейсе число
					num, _ := strconv.ParseFloat(text, 64)
					if op == "-"{
						num = -num
					}
					node.x = &Node{val: num, sheet: true, level: node.level+1, parent: node}
				default:
					return nil, fmt.Errorf("invalid expression: %s", text)
				}
			}
			// если y-элемент определен и соответственно определен оператор
			if node.y != nil{
				node.x = createChild(node)
			} 
			node.op = text	
		case "*", "/":
			switch {
			// когда оператор первый в выражении
			case node.y == nil && node.x != nil:
				node.op = text
			// когда оператор в глубине выражения
			case node.y != nil && node.x != nil:
				childNode := &Node{x: node.y, op: text, level: node.level+1, parent: node}
				var err error
				node.y, err = GetTree(childNode, s)
				if err != nil {
					return nil, fmt.Errorf("invalid expression: %s", text)
				}

				//TODO
			// первый элемент выражения не может быть * или /
			case node.x == nil:
				return nil, fmt.Errorf("invalid expression: %s", text)
			default:
				return nil, fmt.Errorf("invalid expression: %s", text)
			}
			
			//TODO
		case "(":
			//TODO
		case ")":
			//TODO
		default:
			return nil, fmt.Errorf("invalid expression: %s", text)
		}
	}
	node, err := GetTree(node, s)
	if err != nil {
		return nil, err
	}
	return node, nil
}
