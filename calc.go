package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func Calc(expression string) (float64, error) {
	expression = strings.ReplaceAll(expression, " ", "")
	tokens := tokenize(expression)
	var values []float64
	var ops []rune

	for _, token := range tokens {
		if num, err := strconv.ParseFloat(token, 64); err == nil {
			values = append(values, num)
		} else if token == "(" {
			ops = append(ops, '(')
		} else if token == ")" {
			for len(ops) > 0 && ops[len(ops)-1] != '(' {
				if err := applyOp(&values, ops[len(ops)-1]); err != nil {
					return 0, err
				}
				ops = ops[:len(ops)-1]
			}
			ops = ops[:len(ops)-1]
		} else if isOperator(rune(token[0])) {
			for len(ops) > 0 && precedence(ops[len(ops)-1]) >= precedence(rune(token[0])) {
				if err := applyOp(&values, ops[len(ops)-1]); err != nil {
					return 0, err
				}
				ops = ops[:len(ops)-1]
			}
			ops = append(ops, rune(token[0]))
		}
	}

	for len(ops) > 0 {
		if err := applyOp(&values, ops[len(ops)-1]); err != nil {
			return 0, err
		}
		ops = ops[:len(ops)-1]
	}

	if len(values) != 1 {
		return 0, errors.New("invalid expression")
	}
	return values[0], nil
}

func applyOp(values *[]float64, op rune) error {
	if len(*values) < 2 {
		return errors.New("invalid expression")
	}

	a := (*values)[len(*values)-2]
	b := (*values)[len(*values)-1]
	*values = (*values)[:len(*values)-2]

	switch op {
	case '+':
		*values = append(*values, a+b)
	case '-':
		*values = append(*values, a-b)
	case '*':
		*values = append(*values, a*b)
	case '/':
		if b == 0 {
			return errors.New("division by zero")
		}
		*values = append(*values, a/b)
	default:
		return errors.New("unknown operator")
	}

	return nil
}

func precedence(op rune) int {
	if op == '+' || op == '-' {
		return 1
	}

	if op == '*' || op == '/' {
		return 2
	}
	return 0
}

func tokenize(expression string) []string {
	var tokens []string
	var numBuffer strings.Builder

	for _, char := range expression {
		if isDigit(char) || char == '.' {
			numBuffer.WriteRune(char)
		} else {
			if numBuffer.Len() > 0 {
				tokens = append(tokens, numBuffer.String())
				numBuffer.Reset()
			}
			tokens = append(tokens, string(char))
		}
	}

	if numBuffer.Len() > 0 {
		tokens = append(tokens, numBuffer.String())
	}

	return tokens
}

func isDigit(char rune) bool {
	return char >= '0' && char <= '9'
}

func isOperator(char rune) bool {
	return char == '+' || char == '-' || char == '*' || char == '/'
}

func main() {
	result, err := Calc("(300 - 200) + (2 * 100 - 100) / 0.88")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(result)
}
