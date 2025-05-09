package agent

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"unicode"
)

func CalculateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Expression string `json:"expression"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, fmt.Sprintf("ошибка в теле запроса: %v", err), http.StatusBadRequest)
		return
	}

	result, err := computeExpression(request.Expression)
	if err != nil {
		http.Error(w, fmt.Sprintf("ошибка вычисления: %v", err), http.StatusBadRequest)
		return
	}

	response := struct {
		Result float64 `json:"result"`
	}{
		Result: result,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("ошибка при отправке ответа: %v", err)
	}
}

func computeExpression(expression string) (float64, error) {
	expression = strings.ReplaceAll(expression, " ", "")
	tokens, err := tokenize(expression)
	if err != nil {
		return 0, err
	}
	postfix, err := infixToPostfix(tokens)
	if err != nil {
		return 0, err
	}
	result, err := evaluatePostfix(postfix)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func tokenize(expression string) ([]string, error) {
	var tokens []string
	var currentToken strings.Builder

	for i, char := range expression {
		if unicode.IsDigit(char) {
			currentToken.WriteRune(char)
		} else {
			if currentToken.Len() > 0 {
				tokens = append(tokens, currentToken.String())
				currentToken.Reset()
			}

			if char == '-' {
				if i == 0 || (len(tokens) > 0 && (tokens[len(tokens)-1] == "(" || isOperator(tokens[len(tokens)-1]))) {
					currentToken.WriteRune(char)
					continue
				}
			}

			if isOperator(string(char)) || char == '(' || char == ')' {
				tokens = append(tokens, string(char))
			} else if !unicode.IsSpace(char) {
				return nil, fmt.Errorf("неизвестный токен: %c", char)
			}
		}
	}

	if currentToken.Len() > 0 {
		tokens = append(tokens, currentToken.String())
	}

	return tokens, nil
}

func isOperator(s string) bool {
	return s == "+" || s == "-" || s == "*" || s == "/" || s == "^"
}

func isNumber(char string) bool {
	return strings.Contains("0123456789.", char)
}

func infixToPostfix(tokens []string) ([]string, error) {
	var output []string
	var operators []string
	precedence := map[string]int{
		"+": 1, "-": 1,
		"*": 2, "/": 2,
		"^": 3,
	}
	for _, token := range tokens {
		if isNumber(token) {
			output = append(output, token)
		} else if token == "(" {
			operators = append(operators, token)
		} else if token == ")" {
			for len(operators) > 0 && operators[len(operators)-1] != "(" {
				output = append(output, operators[len(operators)-1])
				operators = operators[:len(operators)-1]
			}
			if len(operators) == 0 {
				return nil, fmt.Errorf("ошибка: несогласованные скобки")
			}
			operators = operators[:len(operators)-1]
		} else if isOperator(token) {
			for len(operators) > 0 && precedence[operators[len(operators)-1]] >= precedence[token] &&
				(token != "^" || precedence[operators[len(operators)-1]] > precedence[token]) {
				output = append(output, operators[len(operators)-1])
				operators = operators[:len(operators)-1]
			}
			operators = append(operators, token)
		} else {
			return nil, fmt.Errorf("неизвестный токен: %s", token)
		}
	}
	for len(operators) > 0 {
		if operators[len(operators)-1] == "(" {
			return nil, fmt.Errorf("ошибка: несогласованные скобки")
		}
		output = append(output, operators[len(operators)-1])
		operators = operators[:len(operators)-1]
	}
	return output, nil
}

func evaluatePostfix(tokens []string) (float64, error) {
	var stack []float64
	for _, token := range tokens {
		if isNumber(token) {
			num, err := strconv.ParseFloat(token, 64)
			if err != nil {
				return 0, err
			}
			stack = append(stack, num)
		} else if isOperator(token) {
			if len(stack) < 2 {
				return 0, fmt.Errorf("ошибка: недостаточно аргументов для операции %s", token)
			}
			b := stack[len(stack)-1]
			a := stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			var res float64
			switch token {
			case "+":
				res = a + b
			case "-":
				res = a - b
			case "*":
				res = a * b
			case "/":
				if b == 0 {
					return 0, fmt.Errorf("ошибка: деление на ноль")
				}
				res = a / b
			case "^":
				res = math.Pow(a, b)
			}
			stack = append(stack, res)
		} else {
			return 0, fmt.Errorf("неизвестный токен: %s", token)
		}
	}
	if len(stack) != 1 {
		return 0, fmt.Errorf("ошибка: неверное количество значений в стеке")
	}
	return stack[0], nil
}
