package models

type Expression struct {
	ID         string  `json:"id"`
	Expression string  `json:"expression"`
	Status     string  `json:"status"`
	Result     float64 `json:"result"`
}
type Task struct {
	ID           string  `json:"id"`            // Уникальный ID задачи
	ExpressionID string  `json:"expression_id"` // ID выражения, которому принадлежит задача
	Arg1         string  `json:"arg1"`          // Первый операнд (может быть числом или ID предыдущей задачи)
	Arg2         string  `json:"arg2"`          // Второй операнд (может быть числом или ID предыдущей задачи)
	Operation    string  `json:"operation"`     // Операция (+, -, *, /)
	Status       string  `json:"status"`        // Статус выполнения (PENDING, COMPLETED)
	Result       float64 `json:"result"`        // Результат (если задача выполнена)
}
