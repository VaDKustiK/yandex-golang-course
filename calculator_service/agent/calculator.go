package agent

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
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
	// Возвращаем статус 200 OK
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
	var currentToken string

	for i := 0; i < len(expression); i++ {
		char := string(expression[i])
		if isNumber(char) {
			currentToken += char
		} else if isOperator(char) {
			if currentToken != "" {
				tokens = append(tokens, currentToken)
				currentToken = ""
			}
			tokens = append(tokens, char)
		} else {
			return nil, fmt.Errorf("некорректный символ: %s", char)
		}
	}
	if currentToken != "" {
		tokens = append(tokens, currentToken)
	}
	return tokens, nil
}

func isNumber(char string) bool {
	return strings.Contains("0123456789.", char)
}

func isOperator(char string) bool {
	return strings.Contains("+-*/^", char)
}

func infixToPostfix(tokens []string) ([]string, error) {
	var stack []string
	var output []string
	for _, token := range tokens {
		if isNumber(token) {
			output = append(output, token)
		} else if isOperator(token) {
			for len(stack) > 0 && precedence(stack[len(stack)-1]) >= precedence(token) {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, token)
		}
	}
	for len(stack) > 0 {
		output = append(output, stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}
	return output, nil
}

func precedence(operator string) int {
	switch operator {
	case "+", "-":
		return 1
	case "*", "/":
		return 2
	case "^":
		return 3
	default:
		return 0
	}
}

func evaluatePostfix(postfix []string) (float64, error) {
	var stack []float64
	for _, token := range postfix {
		if isNumber(token) {
			num, err := strconv.ParseFloat(token, 64)
			if err != nil {
				return 0, err
			}
			stack = append(stack, num)
		} else if isOperator(token) {
			if len(stack) < 2 {
				return 0, fmt.Errorf("недостаточно операндов для оператора %s", token)
			}
			b := stack[len(stack)-1]
			a := stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			var result float64
			switch token {
			case "+":
				result = a + b
			case "-":
				result = a - b
			case "*":
				result = a * b
			case "/":
				if b == 0 {
					return 0, fmt.Errorf("деление на ноль")
				}
				result = a / b
			case "^":
				result = power(a, b)
			default:
				return 0, fmt.Errorf("неизвестный оператор: %s", token)
			}
			stack = append(stack, result)
		}
	}
	if len(stack) != 1 {
		return 0, fmt.Errorf("неверное количество результатов в стеке")
	}
	return stack[0], nil
}

func power(base, exponent float64) float64 {
	return math.Pow(base, exponent)
}
