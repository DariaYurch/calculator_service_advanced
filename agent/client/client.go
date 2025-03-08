package client

import (
	"calculator_advanced_v2/orchestrator/services/storage"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

const orchestratorURL = "http://localhost:8080"

type Task struct {
	ID            string `json:"id"`
	Arg1          string `json:"arg1"`
	Arg2          string `json:"arg2"`
	Operation     string `json:"operation"`
	OperationTime int    `json:"operation_time"`
}

func GetTask() (*Task, error) {
	resp, err := http.Get(orchestratorURL + "/api/internal/task")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.New("No avialable tasks")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error getting tasks", resp.StatusCode)
	}
	var task struct {
		Task Task `json:"task"`
	}
	err = json.NewDecoder(resp.Body).Decode(&task)
	if err != nil {
		return nil, err
	}
	//log.Println("TASK TAKEN")
	return &task.Task, nil
}

func SendTask(taskID string, result float64) error {
	var expressionID string
	db := storage.GetDB()
	if db == nil {
		log.Println("Error to connect with DB")
		return sql.ErrConnDone
	}

	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM tasks WHERE id = ?)", taskID).Scan(&exists)
	if err != nil {
		log.Println("Task is not exists", err)
		return err
	}

	if !exists {
		log.Println("Error with task", taskID)
		return sql.ErrNoRows
	}

	_, err = db.Exec("UPDATE tasks SET result = ?, status = 'COMPLETED' WHERE id = ?", result, taskID)
	if err != nil {
		log.Println("Error to update task in DB", err)
		return err
	}
	err = db.QueryRow("SELECT expression_id FROM tasks WHERE id = ?", taskID).Scan(&expressionID)
	if err != nil {
		log.Println("Error to get data from expressions table")
	}
	_, err = db.Exec("UPDATE expressions SET result = ?, status = 'COMPLETED' WHERE id = ?", result, expressionID)
	if err != nil {
		log.Println("Error to update data in expressions table")
	}
	log.Println("COMPLETED send result to db", taskID, ":", result)
	return nil
}
