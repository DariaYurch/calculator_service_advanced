package handlers

import (
	"calculator_advanced_v2/orchestrator/models"
	"calculator_advanced_v2/orchestrator/services"
	"calculator_advanced_v2/utils"
	"encoding/json"
	"net/http"
	"strings"
)

func AddExpression(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if r.Method == "POST" {
		var expression struct {
			Expression string `json:expression`
		}
		err := json.NewDecoder(r.Body).Decode(&expression)
		if utils.ValidateExpression(expression.Expression) != nil {
			http.Error(w, "Error: The expression contains invalid symbols", http.StatusUnprocessableEntity) //422
			return
		}
		if err != nil || expression.Expression == "" {
			http.Error(w, "Expression can not be empty", http.StatusUnprocessableEntity) //422
			return
		}

		taskID, err := utils.GenerateID()
		if err != nil {
			http.Error(w, "Error generating the id", http.StatusInternalServerError) //500
			return
		}

		expressionForDB := models.Expression{
			ID:         taskID,
			Expression: expression.Expression,
			Status:     "PENDING",
		}

		services.SaveExpressionToDB(expressionForDB)
		services.ParseExpressionAndMakeTasks(taskID, expression.Expression)
		w.WriteHeader(http.StatusCreated) //201
		json.NewEncoder(w).Encode(map[string]string{"id": taskID})

	} else {
		http.Error(w, "Method is not supported", http.StatusMethodNotAllowed)
	}
}

func HandlerGetExpressions(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if r.Method == "GET" {
		expressions := services.GetAllExpressions()
		if expressions == nil {
			http.Error(w, "No expressions found", http.StatusInternalServerError)
			return
		}
		type ExpressionResponse struct {
			ID     string  `json:"id"`
			Status string  `json:"status"`
			Result float64 `json:"result"`
		}
		var response []ExpressionResponse
		for _, expr := range expressions {
			response = append(response, ExpressionResponse{
				ID:     expr.ID,
				Status: expr.Status,
				Result: expr.Result,
			})
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"expressions": response,
		})
	} else {
		http.Error(w, "Method is not supported", http.StatusMethodNotAllowed)
	}
}

func HandlerGetExpressionByID(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if r.Method == "GET" {
		type ExpressionResponse struct {
			ID     string  `json:"id"`
			Status string  `json:"status"`
			Result float64 `json:"result"`
		}
		var response []ExpressionResponse
		path := r.URL.Path
		id := strings.TrimPrefix(path, "/api/v1/expressions/")
		if id == "" || strings.Contains(id, "/") {
			http.Error(w, "Invalid id", http.StatusInternalServerError) //500
			return
		}
		expression, found := services.GetExpressionByID(id)
		response = append(response, ExpressionResponse{
			ID:     expression.ID,
			Status: expression.Status,
			Result: expression.Result,
		})
		if !found {
			http.Error(w, "Expression not found", http.StatusNotFound) //404
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"expression": response,
		})
	} else {
		http.Error(w, "Method is not supported", http.StatusMethodNotAllowed)
	}
}
