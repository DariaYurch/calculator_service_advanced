package worker

import (
	"calculator_advanced_v2/agent/client"
	"calculator_advanced_v2/orchestrator/services"
	"calculator_advanced_v2/utils"
	"log"
	"strconv"
	"time"
)

func StartWorker(id int) {
	log.Print("Worker is started", id)

	for {
		task, err := client.GetTask()
		if err != nil {
			log.Print(err)
			time.Sleep(2 * time.Second)
			continue
		}

		if task == nil {
			log.Print("No available tasks", id)
			time.Sleep(2 * time.Second)
			continue
		}

		executionTime := utils.GetTimeByOperator(task.Operation)
		//log.Println("The worker received the task:", task.ID, task.Arg1, task.Operation, task.Arg2, executionTime)
		time.Sleep(time.Duration(executionTime) * time.Millisecond)
		result := performCalculation(task.Operation, task.Arg1, task.Arg2)
		//log.Println("RESULT OF OPERATION ", task.ID, " IS ", result)
		err = client.SendTask(task.ID, result)
		if err != nil {
			log.Println("Error to send result")
		}
	}
}

func performCalculation(op string, arg1, arg2 string) float64 {
	//log.Println("Getting arguments:", arg1, arg2)
	num1, err1 := strconv.ParseFloat(arg1, 64)
	if err1 != nil {
		num1 = services.GetResultTaskByID(arg1)
	}
	//log.Println("NUM1:", num1)

	num2, err2 := strconv.ParseFloat(arg2, 64)
	if err2 != nil {
		num2 = services.GetResultTaskByID(arg2)
	}
	//log.Println("NUM2:", num2)

	switch op {
	case "+":
		return num1 + num2
	case "-":
		return num1 - num2
	case "*":
		return num1 * num2
	case "/":
		if num2 != 0 {
			return num1 / num2
		}
		log.Println("Error: division by zero")
		return 0
	default:
		log.Println("Unknown operation:", op)
		return 0
	}
}
