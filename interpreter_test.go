package main

import (
	"reflect"
	"testing"
)

func TestInterpreter(t *testing.T) {
	// Success cases
	for _, tc := range []struct {
		name string
		expr string
		want interface{}
	}{
		// Declarations
		{"assignment", "$a=1", 1},
		// Expressions
		{"Single integer", "1", 1},
		{"Arithmic egation", "-1", -1},
		{"Logic negation", "!(1=0)", true},
		{"Equality", "1=0", false},
		{"Addition", "1+2", 3},
		{"Addition chain", "1+2+3", 6},
		{"Precedence weakest first", "1+2*3", 7},
		{"Precedence weakest last", "1*2+3", 5},
		{"Paranthesises", "(1+2)*3", 9},
	} {
		t.Run(tc.name, func(t *testing.T) {
			i := NewInterpreter()
			got, err := i.Eval(tc.expr)
			if err != nil {
				t.Fatalf("Unexpected error: %s", err.Error())
			}
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("%s = %d. Want %d", tc.expr, got, tc.want)
			}
		})
	}

	// Error cases
	for _, tc := range []struct {
		name string
		expr string
		want string
	}{
		{"No expression", "", "expected expression, got end of data"},
		{"Unfinished binary expression", "1+", "expected expression, got end of data"},
		{"Unfinished negation", "-", "expected expression, got end of data"},
		{"Unexpected binary operator", "+", "expected expression, got +"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			i := NewInterpreter()
			got, err := i.Eval(tc.expr)
			if err == nil {
				t.Fatalf("%s = %d. Expected error %s", tc.expr, got, err.Error())
			}
			if err.Error() != tc.want {
				t.Fatalf("%s -> %s. Want %s", tc.expr, err.Error(), tc.want)
			}
		})
	}
}
