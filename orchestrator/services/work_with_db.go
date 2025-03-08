package services

import (
	"calculator_advanced_v2/orchestrator/models"
	"calculator_advanced_v2/orchestrator/services/storage"
	"calculator_advanced_v2/utils"
	"database/sql"
	"fmt"
	"log"
)

func SaveExpressionToDB(expression models.Expression) {
	//fmt.Println("Saving expression to DB:", expression.ID, expression.Expression, expression.Status)
	db := storage.GetDB()
	_, err := db.Exec("INSERT INTO expressions (id, expression, status) VALUES (?, ?, ?)", expression.ID, expression.Expression, expression.Status)
	if err != nil {
		log.Println("Error inserting expression to DB to SaveExpressionToDB", err)
		return
	}
}

func GetAllExpressions() []models.Expression {
	db := storage.GetDB()
	rows, err := db.Query("SELECT id, status, result FROM expressions")
	if err != nil {
		log.Println("Error select expressions from DB to GetAllExpression")
		return nil
	}
	defer rows.Close()
	var expressions []models.Expression
	for rows.Next() {
		var expr models.Expression
		var result sql.NullFloat64
		err := rows.Scan(&expr.ID, &expr.Status, &result)
		if err == nil {
			if result.Valid {
				expr.Result = result.Float64
			} else {
				expr.Result = 0
			}
			expressions = append(expressions, expr)
		} else {
			log.Println("Error getting data from DB to GetAllExpressions")
		}
		if len(expressions) == 0 {
			log.Println("DB is empty")
			return nil
		}
	}
	return expressions
}

func GetExpressionByID(id string) (models.Expression, bool) {
	db := storage.GetDB()
	var expr models.Expression
	var result sql.NullFloat64
	err := db.QueryRow("SELECT id, expression, status, result FROM expressions WHERE id = ?", id).Scan(&expr.ID, &expr.Expression, &expr.Status, &result)
	if err != nil {
		log.Println("Error select data from DB to GetExpressionByID", err)
		return models.Expression{}, false
	} else {
		if result.Valid {
			expr.Result = result.Float64
		} else {
			expr.Result = 0
		}
		return expr, true
	}
}

func SaveTaskToBD(task models.Task) {
	db := storage.GetDB()
	_, err := db.Exec("INSERT INTO tasks (id, expression_id, arg1, arg2, operation, status, result) VALUES (?, ?, ?, ?, ?, ?, ?)", task.ID, task.ExpressionID, task.Arg1, task.Arg2, task.Operation, task.Status, task.Result)
	if err != nil {
		log.Println("Error inserting task to DB to SaveTaskToDB:", err)
	}
}

type TaskModified struct {
	ID             string `json:"id"`
	Arg1           string `json:"arg1"`
	Arg2           string `json:"arg2"`
	Operation      string `json:"operation"`
	Operation_time int    `json:"operation_time"`
}

func GetTaskFromDB() (TaskModified, bool) {
	db := storage.GetDB()
	var expressionID string
	var task TaskModified
	err := db.QueryRow("SELECT id, arg1, arg2, operation FROM tasks WHERE status = 'PENDING' LIMIT 1").Scan(&task.ID, &task.Arg1, &task.Arg2, &task.Operation)
	task.Operation_time = utils.GetTimeByOperator(task.Operation)
	if err != nil {
		return TaskModified{}, false
	} else {
		log.Println("Error select data from DB to GetTaskFromDB", err)
	}
	err = db.QueryRow("SELECT expression_id FROM tasks WHERE id = ?", task.ID).Scan(&expressionID)
	if err != nil {
		log.Println("Error to select data from DB to GetTaskFromDB", err)
	}
	//log.Println("EXPRESSION_ID:", expressionID)
	_, err = db.Exec("UPDATE tasks SET status = 'IN_PROGRESS' WHERE id = ?", task.ID)
	if err != nil {
		log.Println("Error to update data from DB (tasks table)", err)
	}
	_, err = db.Exec("UPDATE expressions SET status = 'IN_PROGRESS' WHERE id = ?", expressionID)
	if err != nil {
		log.Println("Error to update data from DB (expressions table)", err)
	}
	return task, true
}

func UpdateTaskResult(id string, result float64) error {
	db := storage.GetDB()
	_, err := db.Exec("UPDATE tasks SET status = 'COMPLETED', result = ? WHERE id = ?", result, id)
	if err != nil {
		return fmt.Errorf("Request execution error: ", err)
	}
	return nil
}

func GetResultTaskByID(taskID string) float64 {
	db := storage.GetDB()
	var result sql.NullFloat64
	err := db.QueryRow("SELECT result FROM tasks WHERE id = ?", taskID).Scan(&result)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("Task is not found", taskID)
		} else {
			log.Println("Error response to DB", err)
		}
		return 0
	}
	if result.Valid {
		return result.Float64
	}
	log.Println("Is not result in DB", taskID)
	return 0
}
