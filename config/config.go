package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

var TimeAdditionMs, TimeSubtractionMs, TimeMultiplicationMs, TimeDivisionMs int

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error to open .env file")
	}
	TimeAdditionMs = getEnvInt("TIME_ADDITION_MS", 1000)
	TimeSubtractionMs = getEnvInt("TIME_SUBTRACTION_MS", 1000)
	TimeMultiplicationMs = getEnvInt("TIME_MULTIPLICATIONS_MS", 2000)
	TimeDivisionMs = getEnvInt("TIME_DIVISIONS_MS", 3000)
}

func getEnvInt(name string, defaultValue int) int {
	val, err := strconv.Atoi(os.Getenv(name))
	if err != nil {
		return defaultValue
	}
	return val
}
