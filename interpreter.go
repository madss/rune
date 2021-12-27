package main

import (
	"errors"
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
	case node.Type == '$':
		ident := rune(node.Args[0].Type)
		val, err := i.interpret(node.Args[1])
		if err != nil {
			return nil, err
		}
		i.env[ident] = val
		return val, nil
	case node.Type.IsDigit():
		return int(node.Type - '0'), nil
	case node.Type.IsIdent():
		val, ok := i.env[rune(node.Type)]
		if !ok {
			return nil, fmt.Errorf("undefined variable %c", node.Type)
		}
		return val, nil
	case node.Type == '+':
		left, err := i.interpret(node.Args[0])
		if err != nil {
			return nil, err
		}
		right, err := i.interpret(node.Args[1])
		if err != nil {
			return nil, err
		}
		return left.(int) + right.(int), nil
	case node.Type == '-':
		if len(node.Args) == 1 {
			// Negation
			val, err := i.interpret(node.Args[0])
			if err != nil {
				return nil, err
			}
			return -val.(int), nil
		} else {
			// Subtraction
			left, err := i.interpret(node.Args[0])
			if err != nil {
				return nil, err
			}
			right, err := i.interpret(node.Args[1])
			if err != nil {
				return nil, err
			}
			return left.(int) - right.(int), nil
		}
	case node.Type == '*':
		left, err := i.interpret(node.Args[0])
		if err != nil {
			return nil, err
		}
		right, err := i.interpret(node.Args[1])
		if err != nil {
			return nil, err
		}
		return left.(int) * right.(int), nil
	case node.Type == '/':
		left, err := i.interpret(node.Args[0])
		if err != nil {
			return nil, err
		}
		right, err := i.interpret(node.Args[1])
		if err != nil {
			return nil, err
		}
		if right.(int) == 0 {
			return nil, errors.New("dividing by zero")
		}
		return left.(int) / right.(int), nil
	}
	panic("invalid node")
}
