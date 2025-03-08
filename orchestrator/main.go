package main

import (
	"calculator_advanced_v2/orchestrator/handlers"
	"calculator_advanced_v2/orchestrator/services/storage"
	"log"
	"net/http"
)

func main() {
	storage.InitConnectionDB()
	http.HandleFunc("/api/v1/calculate", handlers.AddExpression)
	http.HandleFunc("/api/v1/expressions", handlers.HandlerGetExpressions)
	http.HandleFunc("/api/v1/expressions/", handlers.HandlerGetExpressionByID)

	http.HandleFunc("/api/internal/task", handlers.TaskHandler)

	log.Println("Orchestrator is running")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
