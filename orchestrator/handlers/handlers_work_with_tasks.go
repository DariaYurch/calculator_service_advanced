package handlers

import (
	"calculator_advanced_v2/orchestrator/services"
	"encoding/json"
	"net/http"
)

func TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		GetTaskHandler(w, r)
	case "POST":
		PostTaskHandler(w, r)
	default:
		http.Error(w, "Method is not supported", http.StatusMethodNotAllowed)
	}
}

func PostTaskHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		ID     string  `json:"id"`
		Result float64 `json:"result"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid data", http.StatusUnprocessableEntity)
		return
	}
	err = services.UpdateTaskResult(request.ID, request.Result)
	if err != nil {
		http.Error(w, "Failed update data in DB", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

}

func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if r.Method == "GET" {
		task, found := services.GetTaskFromDB()
		if !found {
			http.Error(w, "No available tasks", http.StatusNotFound) //404
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"task": task,
		})
	} else {
		http.Error(w, "Method is not supported", http.StatusMethodNotAllowed)
	}
	w.WriteHeader(http.StatusOK)
}
