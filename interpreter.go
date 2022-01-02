package main

import (
	"fmt"
)

type Interpreter struct {
	env map[rune]interface{}
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		env: map[rune]interface{}{},
	}
}

func (i *Interpreter) Eval(s string) (interface{}, error) {
	ts := tokenize(s)

	node, err := parse(&ts)
	if err != nil {
		return nil, err
	}

	val, err := i.interpret(node)
	if err != nil {
		return nil, err
	}

	return val, nil
}

func (i *Interpreter) interpret(node Node) (interface{}, error) {
	switch {
	case node.Type.IsDigit() && len(node.Args) == 0:
		return int(node.Type - '0'), nil
	case node.Type.IsIdent() && len(node.Args) == 0:
		val, ok := i.env[rune(node.Type)]
		if !ok {
			return nil, fmt.Errorf("undefined variable %c", node.Type)
		}
		return val, nil
	case node.Type == '$' && len(node.Args) == 2:
		ident := rune(node.Args[0].Type)
		val, err := i.interpret(node.Args[1])
		if err != nil {
			return nil, err
		}
		i.env[ident] = val
		return val, nil
	case node.Type == '!' && len(node.Args) == 1:
		val, err := i.interpretBool(node.Args[0])
		return !val, err
	case node.Type == '&' && len(node.Args) == 2:
		left, err := i.interpretBool(node.Args[0])
		if err != nil || !left {
			return left, err
		}
		return i.interpretBool(node.Args[1])
	case node.Type == '|' && len(node.Args) == 2:
		left, err := i.interpretBool(node.Args[0])
		if err != nil || left {
			return left, err
		}
		return i.interpretBool(node.Args[1])
	case node.Type == '-' && len(node.Args) == 1:
		val, err := i.interpret(node.Args[0])
		if err != nil {
			return nil, err
		}
		switch val := val.(type) {
		case int:
			return -val, nil
		default:
			return nil, fmt.Errorf("expected numberic value")
		}
	case node.Type.IsOneOf('=', '<', '>', '+', '-', '*', '/', '%') && len(node.Args) == 2:
		left, err := i.interpret(node.Args[0])
		if err != nil {
			return nil, err
		}
		right, err := i.interpret(node.Args[1])
		if err != nil {
			return nil, err
		}
		switch {
		case isBoolPair(left, right):
			lval := left.(bool)
			rval := right.(bool)
			switch node.Type {
			case '=':
				return lval == rval, nil
			default:
				return nil, fmt.Errorf("unexpected boolean values")
			}
		case isIntPair(left, right):
			lval := left.(int)
			rval := right.(int)
			switch node.Type {
			case '=':
				return lval == rval, nil
			case '<':
				return lval < rval, nil
			case '>':
				return lval > rval, nil
			case '+':
				return lval + rval, nil
			case '-':
				return lval - rval, nil
			case '*':
				return lval * rval, nil
			case '/':
				if rval == 0 {
					return nil, fmt.Errorf("dividing by zero")
				}
				return lval / rval, nil
			case '%':
				if rval == 0 {
					return nil, fmt.Errorf("dividing by zero")
				}
				return lval % rval, nil
			default:
				panic("invalid operator")
			}
		default:
			return nil, fmt.Errorf("unexpected values")
		}
	default:
		panic("invalid node")
	}
}

func (i *Interpreter) interpretBool(node Node) (bool, error) {
	val, err := i.interpret(node)
	if err != nil {
		return false, err
	}
	switch val := val.(type) {
	case bool:
		return val, nil
	default:
		return false, fmt.Errorf("expected boolean value")
	}
}

func isBoolPair(left, right interface{}) bool {
	_, okleft := left.(bool)
	_, okright := right.(bool)
	return okleft && okright
}

func isIntPair(left, right interface{}) bool {
	_, okleft := left.(int)
	_, okright := right.(int)
	return okleft && okright
}
