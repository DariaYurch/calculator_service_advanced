package utils

import (
	"calculator_advanced_v2/config"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"strings"
	"sync/atomic"
	"unicode"
)

var taskCounter int32 = 0

func GenerateGlobalIDForTask() string {
	id := atomic.AddInt32(&taskCounter, 1)
	return fmt.Sprintf("task_%03d", id)
}

func GetTimeByOperator(operator string) int {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Download error .env file")
	}
	switch operator {
	case "+":
		return config.TimeAdditionMs
	case "-":
		return config.TimeSubtractionMs
	case "*":
		return config.TimeMultiplicationMs
	case "/":
		return config.TimeDivisionMs
	default:
		log.Println("Unknown operator:", operator)
		return 1000
	}
}

func GenerateID() (string, error) {
	bytes := make([]byte, 8)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("Error generating a new ID")
	}
	return hex.EncodeToString(bytes), nil
}

func ValidateExpression(expression string) error {
	expression = strings.ReplaceAll(expression, " ", "")
	var bracketCnt int

	for _, elem := range expression {
		if !(unicode.IsDigit(elem) || strings.ContainsRune("+-*/()", elem)) {
			return errors.New("Error: The expression contains invalid characters")
		}
		if string(elem) == "(" {
			bracketCnt++
		} else if string(elem) == ")" {
			bracketCnt--
		}
	}
	if bracketCnt != 0 {
		return errors.New("Error: Unbalanced brackets")
	}
	return nil
}
