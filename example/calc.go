package main

import (
	. "github.com/hydra13142/webui"
	"net/http"
	"strconv"
	"strings"
)

var (
	last float64
	op   string
	in   bool
)

func init() {
	last, op, in = 0, "=", false
}

func main() {
	w := &Window{
		Width:  180,
		Height: 220,
		Sub: []Object{
			&Container{
				Common: Common{"ct", "", 10, 10, 160, 200, nil},
				Sub: []Object{
					&Text{Common: Common{"tx", "0", 0, 0, 160, 40, nil}, Readonly: true},
					&Button{Common: Common{"n0", "0", 00, 160, 40, 40, func(p Param) (Param, error) {
						if in {
							s := p["tx"]
							if len(s) >= 16 {
								return Param{}, nil
							}
							if s != "0" {
								s = s + "0"
							}
							return Param{"tx": s}, nil
						} else {
							in = true
							return Param{"tx": "0"}, nil
						}
					}}},
					&Button{Common: Common{"n1", "1", 00, 120, 40, 40, add("1")}},
					&Button{Common: Common{"n2", "2", 40, 120, 40, 40, add("2")}},
					&Button{Common: Common{"n3", "3", 80, 120, 40, 40, add("3")}},
					&Button{Common: Common{"n4", "4", 00, 80, 40, 40, add("4")}},
					&Button{Common: Common{"n5", "5", 40, 80, 40, 40, add("5")}},
					&Button{Common: Common{"n6", "6", 80, 80, 40, 40, add("6")}},
					&Button{Common: Common{"n7", "7", 00, 40, 40, 40, add("7")}},
					&Button{Common: Common{"n8", "8", 40, 40, 40, 40, add("8")}},
					&Button{Common: Common{"n9", "9", 80, 40, 40, 40, add("9")}},
					&Button{Common: Common{"pt", ".", 40, 160, 40, 40, func(p Param) (Param, error) {
						if in {
							s := p["tx"]
							if len(s) >= 16 {
								return Param{}, nil
							}
							if !strings.ContainsRune(s, '.') {
								s = s + "."
							}
							return Param{"tx": s}, nil
						} else {
							in = true
							return Param{"tx": "0."}, nil
						}
					}}},
					&Button{Common: Common{"o1", "+", 120, 40, 40, 40, run("+")}},
					&Button{Common: Common{"o2", "-", 120, 80, 40, 40, run("-")}},
					&Button{Common: Common{"o3", "*", 120, 120, 40, 40, run("*")}},
					&Button{Common: Common{"o4", "/", 120, 160, 40, 40, run("/")}},
					&Button{Common: Common{"eq", "=", 80, 160, 40, 40, run("=")}},
				},
			},
		},
	}
	h := NewHandler(w, "calc.htm")
	http.ListenAndServe(":9999", h)
}

func opr(a, b float64, o string) float64 {
	switch o {
	case "+":
		return a + b
	case "-":
		return a - b
	case "*":
		return a * b
	case "/":
		return a / b
	case "=":
		return b
	}
	return a
}
func add(c string) func(p Param) (Param, error) {
	return func(p Param) (Param, error) {
		if in {
			s := p["tx"]
			if len(s) >= 16 {
				return Param{}, nil
			}
			if s[0] == '0' && !strings.ContainsRune(s, '.') {
				s = s[1:]
			}
			return Param{"tx": s + c}, nil
		} else {
			in = true
			return Param{"tx": c}, nil
		}
	}
}
func run(o string) func(p Param) (Param, error) {
	return func(p Param) (Param, error) {
		x, _ := strconv.ParseFloat(p["tx"], 64)
		x = opr(last, x, op)
		last, op, in = x, o, false
		return Param{"tx": strconv.FormatFloat(x, 'g', -1, 64)}, nil
	}
}
