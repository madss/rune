package main

import (
	"fmt"
	"strings"
)

type Node struct {
	Type token
	Args []Node
}

func (n *Node) String() string {
	b := strings.Builder{}
	b.WriteRune('(')
	b.WriteRune(rune(n.Type))
	for _, arg := range n.Args {
		b.WriteRune(' ')
		b.WriteString(arg.String())
	}
	b.WriteRune(')')
	return b.String()
}

func precedence(t token) int {
	switch t {
	case '*', '/':
		return 2
	case '+', '-':
		return 1
	default:
		return 0
	}
}

func parse(ts *tokens) (Node, error) {
	node, err := parseDecl(ts)
	if err != nil {
		return Node{}, err
	}
	if t := ts.Current(); t != 0 {
		return Node{}, fmt.Errorf("expected end of data, got %s", t)
	}
	return node, nil
}

func parseDecl(ts *tokens) (Node, error) {
	t := ts.Current()
	switch t {
	case '$':
		ts.Advance()
		ident := ts.Current()
		if !ident.IsIdent() {
			return Node{}, fmt.Errorf("expected identifier, got %s", t)
		}
		ts.Advance()
		if ts.Current() != '=' {
			return Node{}, fmt.Errorf("expected =, got %s", t)
		}
		ts.Advance()
		n, err := parseExpr(ts, 0) // FIXME: Is this a good precedence?!
		if err != nil {
			return Node{}, fmt.Errorf("expected expression, got %s", ts.Current())
		}
		return Node{t, []Node{{Type: ident}, n}}, nil
	default:
		return parseExpr(ts, 0)
	}
}

func parseExpr(ts *tokens, prec int) (Node, error) {
	logf(prec, "parseExpr(%d)", prec)
	var node Node

	t := ts.Current()
	logf(prec, "Found prefix %c", t)
	switch {
	case t.IsDigit() || t.IsIdent():
		ts.Advance()
		node = Node{t, nil}
	case t == '-':
		ts.Advance()
		other, err := parseExpr(ts, 0)
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
			logf(prec, "Found ending %c", t)
			return Node{}, fmt.Errorf("expected ), got %s", t)
		}
		ts.Advance()
		node = n
	default:
		return node, fmt.Errorf("expected expression, got %s", t)
	}

	for {
		t := ts.Current()
		logf(prec, "Found suffix %c (%d)", t, precedence(t))

		if precedence(t) < prec {
			logf(prec, "too low precedence")
			return node, nil

		}
		switch t {
		case '+', '-', '*', '/':
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

func logf(indent int, format string, args ...interface{}) {
	// fmt.Printf(strings.Repeat(" ", indent)+format+"\n", args...)
}
