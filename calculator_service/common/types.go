package common

// CalcRequest – структура запроса на вычисление выражения.
type CalcRequest struct {
	Expression string `json:"expression"`
}

// CalcResponse – структура ответа после вычисления выражения.
type CalcResponse struct {
	Result float64 `json:"result,omitempty"`
	Error  string  `json:"error,omitempty"`
}

// Expression представляет арифметическое выражение, его статус и связанные задачи.
type Expression struct {
	ID         int      `json:"id"`
	Expression string   `json:"expression"`
	Status     string   `json:"status"`           // "pending", "in_progress", "completed"
	Result     *float64 `json:"result,omitempty"` // итоговый результат, если вычислено
	TaskIDs    []int    `json:"-"`                // идентификаторы задач, связанных с выражением
}

// Task представляет отдельную вычислительную задачу (операцию).
type Task struct {
	ID            int      `json:"id"`
	ExprID        int      `json:"expr_id"`
	Arg1          float64  `json:"arg1"`
	Arg2          float64  `json:"arg2"`
	Operation     string   `json:"operation"` // "+", "-", "*", "/"
	Result        *float64 `json:"result,omitempty"`
	Status        string   `json:"status"`         // "pending", "completed"
	OperationTime int      `json:"operation_time"` // время выполнения операции в миллисекундах
}

// TaskResultRequest – структура для передачи результата выполнения задачи.
type TaskResultRequest struct {
	ID     int     `json:"id"`
	Result float64 `json:"result"`
}
