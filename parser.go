package main

import (
	"fmt"
)

type Node struct {
	Type token
	Args []Node
}

func precedence(t token) int {
	switch t {
	case '.':
		return 8
	case '!': // Also unary minus that cannot be represented because of binary minus
		return 7
	case '*', '/', '%':
		return 6
	case '+', '-':
		return 5
	case '<', '>':
		return 4
	case '=':
		return 3
	case '&':
		return 2
	case '|':
		return 1
	default:
		return 0
	}
}

func parse(ts *tokens) (Node, error) {
	node, err := parseDecls(ts)
	if err != nil {
		return Node{}, err
	}
	if t := ts.Current(); t != 0 {
		return Node{}, fmt.Errorf("expected end of data, got %s", t)
	}
	return node, nil
}

func parseDecls(ts *tokens) (Node, error) {
	decls := Node{Type: ';'}
	for {
		decl, err := parseDecl(ts)
		if err != nil {
			return Node{}, err
		}
		decls.Args = append(decls.Args, decl)

		if ts.Current() != ';' {
			return decls, nil
		}
		ts.Advance()
	}
}

func parseDecl(ts *tokens) (Node, error) {
	t := ts.Current()
	switch t {
	case '$':
		ts.Advance()
		name := ts.Current()
		if !name.IsIdent() {
			return Node{}, fmt.Errorf("expected identifier, got %s", name)
		}
		ts.Advance()
		if ts.Current() != '=' {
			return Node{}, fmt.Errorf("expected =, got %s", t)
		}
		ts.Advance()
		n, err := parseExpr(ts, 0)
		if err != nil {
			return Node{}, fmt.Errorf("expected expression, got %s", ts.Current())
		}
		return Node{t, []Node{{Type: name}, n}}, nil
	default:
		return parseExpr(ts, 0)
	}
}

func parseExpr(ts *tokens, prec int) (Node, error) {
	var node Node

	t := ts.Current()
	switch {
	case t.IsDigit() || t.IsIdent():
		ts.Advance()
		node = Node{t, nil}
	case t.IsOneOf('-', '!'):
		ts.Advance()
		other, err := parseExpr(ts, precedence('!'))
		if err != nil {
			return Node{}, err
		}
		node = Node{t, []Node{other}}
	case t == '(':
		ts.Advance()
		n, err := parseExpr(ts, 0)
		if err != nil {
			return Node{}, err
		}
		if t := ts.Current(); t != ')' {
			return Node{}, fmt.Errorf("expected ), got %s", t)
		}
		ts.Advance()
		node = n
	default:
		return node, fmt.Errorf("expected expression, got %s", t)
	}

	for {
		t := ts.Current()

		if precedence(t) <= prec {
			return node, nil
		}
		switch {
		case t.IsOneOf('+', '-', '*', '/', '%', '=', '<', '>', '&', '|', '.'):
			ts.Advance()
			other, err := parseExpr(ts, precedence(t))
			if err != nil {
				return Node{}, err
			}
			node = Node{t, []Node{node, other}}
		default:
			return node, nil
		}
	}
}
