package main

import (
	"fmt"
	"math"
	"reflect"
	"testing"
)

type TestGroup struct {
	Name  string
	Cases []TestCase
}

type TestCase struct {
	Source string
	Want   interface{}
	Error  string
}

func TestInterpreter(t *testing.T) {
	// Success cases
	for _, tg := range []TestGroup{
		{
			"No data",
			[]TestCase{
				{"", nil, "expected expression, got end of data"},
			},
		},
		{
			"Integer Literals",
			[]TestCase{
				{"0", 0, ""},
				{"3", 3, ""},
			},
		},
		{
			"Floating point operator",
			[]TestCase{
				{"1.0", 1.0, ""},
				{"(1+2).(9+5)", 3.14, ""},
				{"1.2*3.4", 4.08, ""},
				{"1.(-2)", nil, "negative decimal part"},
				{"(1.2).3", nil, "unexpected values"},
				{"1.(2.3)", nil, "unexpected values"},
			},
		},
		{
			"Unary minus operator",
			[]TestCase{
				{"-2", -2, ""},
				{"-2.3", -2.3, ""},
				{"-2-2", -4, ""},
				{"-", nil, "expected expression, got end of data"},
			},
		},
		{
			"Binary plus operator",
			[]TestCase{
				{"2+3", 5, ""},
				{"2.3+4.5", 6.8, ""},
				{"1+2.0", nil, "unexpected values"},
				{"2+3+4", 9, ""},
				{"1+", nil, "expected expression, got end of data"},
			},
		},
		{
			"Binary minus operator",
			[]TestCase{
				{"3-2", 1, ""},
				{"5.4-3.2", 2.2, ""},
				{"3-2.0", nil, "unexpected values"},
				{"4-3-2", -1, ""},
				{"1-", nil, "expected expression, got end of data"},
			},
		},
		{
			"Binary times operator",
			[]TestCase{
				{"2*3", 6, ""},
				{"2.3*4.5", 10.35, ""},
				{"3*2.0", nil, "unexpected values"},
				{"2*3*4", 24, ""},
				{"1*", nil, "expected expression, got end of data"},
			},
		},
		{
			"Binary divide operator",
			[]TestCase{
				{"5/2", 2, ""},
				{"5.0/2.0", 2.5, ""},
				{"5/2.0", nil, "unexpected values"},
				{"6/3/2", 1, ""},
				{"5/0", nil, "dividing by zero"},
				{"5.0/0.0", math.Inf(1), ""},
				{"5/", nil, "expected expression, got end of data"},
			},
		},
		{
			"Binary modulo operator",
			[]TestCase{
				{"5%3", 2, ""},
				{"5.0%3.0", nil, "invalid floating point operator"},
				{"5%3.0", nil, "unexpected values"},
				{"5%3%3", 2, ""},
				{"3%0", nil, "dividing by zero"},
				{"5%", nil, "expected expression, got end of data"},
			},
		},
		{
			"Binary equal operator",
			[]TestCase{
				{"2=2", true, ""},
				{"2=3", false, ""},
				{"2.0=2.0", true, ""},
				{"2.0=3.0", false, ""},
				{"(2=2)=(2=2)", true, ""},
				{"(2=2)=(2=3)", false, ""},
				{"2=3.0", nil, "unexpected values"},
				{"1=2=(3=4)", true, ""},
				{"5=", nil, "expected expression, got end of data"},
			},
		},
		{
			"Unary not operator",
			[]TestCase{
				{"!(2=3)", true, ""},
				{"!(2=2)", false, ""},
				{"!!(2=2)", true, ""},
				{"!2", nil, "expected boolean value"},
				{"!", nil, "expected expression, got end of data"},
			},
		},
		{
			"Binary less than operator",
			[]TestCase{
				{"2<3", true, ""},
				{"2<2", false, ""},
				{"2.0<2.1", true, ""},
				{"2.0<2.0", false, ""},
				{"2<3.0", nil, "unexpected values"},
				{"1<2<3", nil, "unexpected values"},
				{"5<", nil, "expected expression, got end of data"},
			},
		},
		{
			"Binary greater than operator",
			[]TestCase{
				{"3>2", true, ""},
				{"2>2", false, ""},
				{"2.1>2.0", true, ""},
				{"2.0>2.0", false, ""},
				{"2>3.0", nil, "unexpected values"},
				{"1>2>3", nil, "unexpected values"},
				{"5>", nil, "expected expression, got end of data"},
			},
		},
		{
			"Binary and operator",
			[]TestCase{
				{"(1=1)&(1=1)", true, ""},
				{"(1=1)&(1=2)", false, ""},
				{"(1=2)&(1=1)", false, ""},
				{"(1=2)&(1=2)", false, ""},
				{"(1=2)&(1/0)", false, ""},
				{"2&3", nil, "expected boolean value"},
				{"(1=1)&(1=1)&(1=1)", true, ""},
				{"(1=1)&", nil, "expected expression, got end of data"},
			},
		},
		{
			"Binary or operator",
			[]TestCase{
				{"(1=1)|(1=1)", true, ""},
				{"(1=1)|(1=2)", true, ""},
				{"(1=2)|(1=1)", true, ""},
				{"(1=2)|(1=2)", false, ""},
				{"(1=1)|(1/0)", true, ""},
				{"2|3", nil, "expected boolean value"},
				{"(1=2)|(1=2)|(1=1)", true, ""},
				{"(1=2)|", nil, "expected expression, got end of data"},
			},
		},
		{
			"Multiple declarations",
			[]TestCase{
				{"1", 1, ""},
				{"1;2", 2, ""},
				{"1;2;3", 3, ""},
				{"1;", nil, "expected expression, got end of data"},
			},
		},
		{
			"Assignment",
			[]TestCase{
				{"$a=1", 1, ""},
				{"$a=1;a", 1, ""},
				{"$", nil, "expected identifier, got end of data"},
				{"$1=2", nil, "expected identifier, got 1"},
				{"$a=", nil, "expected expression, got end of data"},
			},
		},
	} {
		for i, tc := range tg.Cases {
			t.Run(fmt.Sprintf("%s [%d]", tg.Name, i), func(t *testing.T) {
				i := NewInterpreter()
				got, err := i.Eval(tc.Source)
				if err != nil {
					if tc.Error == "" {
						t.Fatalf("Unexpected error: %s", err.Error())
					}
					if tc.Error != err.Error() {
						t.Fatalf("Got error: %s. Want %s", err.Error(), tc.Error)
					}
					return
				}
				if !reflect.DeepEqual(tc.Want, got) {
					t.Fatalf("%s = %v. Want %v", tc.Source, got, tc.Want)
				}
			})
		}
	}
}
