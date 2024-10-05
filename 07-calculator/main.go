package main

import (
	"fmt"
	"os"
	"strconv"
)

type operator struct {
	symbol     string
	precedence int
}

func (o operator) isLowerPrecedenceThan(o2 operator) bool {
	return o.precedence < o2.precedence
}

var operatorMap = map[string]operator{
	"+": {"+", 1},
	"-": {"-", 1},
	"*": {"*", 2},
	"/": {"*", 2},
	"^": {"^", 3},
}

func main() {
	// expression := "( 1 + 2 ) * ( 8 / 4 ) ^ ( 1 + 1 )"
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go \"( 1 + 2 ) * ( 8 / 4 ) ^ ( 1 + 1 )\"")
		return
	}
	expression := os.Args[1]

	tokens := parseExpression(expression)
	rpn := shuntingYard(tokens)
	result := evaluateRPN(rpn)
	fmt.Println(result)
}

func parseExpression(expression string) []string {
	tokens := []string{}
	for i := 0; i < len(expression); i++ {
		if expression[i] == ' ' {
			continue
		}
		if expression[i] >= '0' && expression[i] <= '9' {
			j := i
			for j < len(expression) && expression[j] >= '0' && expression[j] <= '9' {
				j++
			}
			tokens = append(tokens, expression[i:j])
			i = j - 1
		} else {
			tokens = append(tokens, string(expression[i]))
		}
	}
	return tokens
}

func shuntingYard(tokens []string) []string {
	output := []string{}
	operators := []string{}
	for _, token := range tokens {
		if isOperator(token) {
			output, operators = handleOperator(token, output, operators)
		} else if token == "(" {
			operators = append(operators, token)
		} else if token == ")" {
			output, operators = handleClosingParenthesis(output, operators)
		} else {
			output = append(output, token)
		}
	}
	return appendRemainingOperators(output, operators)
}

func isOperator(token string) bool {
	_, ok := operatorMap[token]
	return ok
}

func handleOperator(token string, output, operators []string) ([]string, []string) {
	currentOperator := operatorMap[token]
	for len(operators) > 0 {
		lastToken := operators[len(operators)-1]
		if !isOperator(lastToken) {
			break
		}
		lastOperator := operatorMap[lastToken]
		if lastOperator.isLowerPrecedenceThan(currentOperator) {
			break
		}
		output = append(output, operators[len(operators)-1])
		operators = operators[:len(operators)-1]
	}
	operators = append(operators, token)
	return output, operators
}

func handleClosingParenthesis(output, operators []string) ([]string, []string) {
	for {
		if len(operators) == 0 {
			panic("Mismatched parentheses")
		}
		if operators[len(operators)-1] == "(" {
			break
		}
		output = append(output, operators[len(operators)-1])
		operators = operators[:len(operators)-1]
	}
	operators = operators[:len(operators)-1]
	return output, operators
}

func appendRemainingOperators(output, operators []string) []string {
	for len(operators) > 0 {
		output = append(output, operators[len(operators)-1])
		operators = operators[:len(operators)-1]
	}
	return output
}

func evaluateRPN(rpn []string) int {
	stack := []int{}
	for _, token := range rpn {
		if isOperator(token) {
			if len(stack) < 2 {
				panic("Invalid RPN expression")
			}
			operandRight := stack[len(stack)-1]
			operandLeft := stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			stack = append(stack, applyOperator(token, operandLeft, operandRight))
		} else {
			tokenInt, err := strconv.Atoi(token)
			if err != nil {
				panic(err)
			}
			stack = append(stack, tokenInt)
		}
	}
	if len(stack) != 1 {
		panic("Invalid RPN expression")
	}
	return stack[0]
}

func applyOperator(operator string, left, right int) int {
	switch operator {
	case "+":
		return left + right
	case "-":
		return left - right
	case "*":
		return left * right
	case "/":
		if right == 0 {
			panic("Division by zero")
		}
		return left / right
	case "^":
		result := 1
		for i := 0; i < right; i++ {
			result *= left
		}
		return result
	default:
		panic("Unknown operator")
	}
}
