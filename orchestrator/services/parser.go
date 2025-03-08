package services

import (
	"calculator_advanced_v2/orchestrator/models"
	"calculator_advanced_v2/utils"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func brokeString(str string) []string {
	var result []string
	var partOfString strings.Builder
	str = strings.ReplaceAll(str, " ", "")

	for _, elem := range str {
		if elem == '+' || elem == '-' || elem == '/' || elem == '*' || elem == '(' || elem == ')' {
			if partOfString.Len() > 0 {
				result = append(result, partOfString.String())
			}
			partOfString.Reset()
			result = append(result, string(elem))
		} else {
			partOfString.WriteRune(elem)
		}
	}

	if partOfString.Len() > 0 {
		result = append(result, partOfString.String())
	}

	return result
}

func isNumber(element string) bool {
	_, err := strconv.ParseFloat(element, 64)
	return err == nil
}

func isOperator(element string) bool {
	if element == "+" || element == "-" || element == "/" || element == "*" {
		return true
	}
	return false
}

func getPriority(operator string) int {
	switch operator {
	case "+", "-":
		return 1
	case "*", "/":
		return 2
	default:
		return 0
	}
}

func to_prefix_form(main_expression []string) ([]string, error) {
	var result []string
	var stackOperators []string
	for _, element := range main_expression {
		if isNumber(element) {
			result = append(result, element)
		} else if element == "(" {
			stackOperators = append(stackOperators, element)
		} else if element == ")" {
			for len(stackOperators) > 0 && stackOperators[len(stackOperators)-1] != "(" {
				result = append(result, stackOperators[len(stackOperators)-1])
				stackOperators = stackOperators[:len(stackOperators)-1]
			}
			if len(stackOperators) == 0 {
				return nil, errors.New("Error")
			}
			stackOperators = stackOperators[:len(stackOperators)-1]
		} else if isOperator(element) {
			for len(stackOperators) > 0 && getPriority(stackOperators[len(stackOperators)-1]) >= getPriority(element) {
				result = append(result, stackOperators[len(stackOperators)-1])
				stackOperators = stackOperators[:len(stackOperators)-1]
			}
			stackOperators = append(stackOperators, element)
		} else {
			return nil, errors.New("Error")
		}
	}
	for len(stackOperators) > 0 {
		if stackOperators[len(stackOperators)-1] == "(" {
			return nil, errors.New("Error")
		}
		result = append(result, stackOperators[len(stackOperators)-1])
		stackOperators = stackOperators[:len(stackOperators)-1]
	}
	return result, nil
}

func ParseExpressionAndMakeTasks(expressionID string, expression string) {
	var task models.Task
	var tasks []models.Task
	var stack []string
	var arg1, arg2 string
	var lastTaskID string
	brokeStr := brokeString(expression)
	workExpression, _ := to_prefix_form(brokeStr)
	fmt.Println("workExpression:", workExpression)
	for _, elem := range workExpression {
		if isNumber(elem) {
			stack = append(stack, elem)
		} else if isOperator(elem) {
			if len(stack) < 2 {
				fmt.Errorf("Calculation error. Wrong input")
			}
			if len(stack) > 0 {
				arg2 = stack[len(stack)-1]
				stack = stack[:len(stack)-1]
			} else {
				arg2 = lastTaskID
			}
			if len(stack) > 0 {
				arg1 = stack[len(stack)-1]
				stack = stack[:len(stack)-1]
			} else {
				arg1 = lastTaskID
			}

			//taskID := fmt.Sprintf("task%03d", taskCounter)
			taskID := utils.GenerateGlobalIDForTask()
			stack = append(stack, taskID)
			//fmt.Println("TASK:", taskID, "ARGS:", arg1, arg2, "OPERATOR:", elem)
			task = models.Task{
				ID:           taskID,
				ExpressionID: expressionID,
				Arg1:         arg1,
				Arg2:         arg2,
				Operation:    elem,
				Status:       "PENDING",
			}
			tasks = append(tasks, task)
			//fmt.Println("TASK:", task)
		} else {
			fmt.Errorf("Error. Invalid token: %s", elem)
		}
	}
	fmt.Println("ALL TASKS:", tasks)
	for _, task := range tasks {
		SaveTaskToBD(task)
	}
}
