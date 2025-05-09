package common

type CalcRequest struct {
	Expression string `json:"expression"`
}

type CalcResponse struct {
	Result float64 `json:"result,omitempty"`
	Error  string  `json:"error,omitempty"`
}

type Expression struct {
	ID         int      `json:"id"`
	Expression string   `json:"expression"`
	Status     string   `json:"status"`
	Result     *float64 `json:"result,omitempty"`
	TaskIDs    []int    `json:"-"`
}

type Task struct {
	ID            int      `json:"id"`
	ExprID        int      `json:"expr_id"`
	Expression    string   `json:"expression"`
	Arg1          float64  `json:"arg1"`
	Arg2          float64  `json:"arg2"`
	Operation     string   `json:"operation"`
	Result        *float64 `json:"result,omitempty"`
	Status        string   `json:"status"`
	OperationTime int      `json:"operation_time"`
}

type TaskResultRequest struct {
	ID     int     `json:"id"`
	Result float64 `json:"result"`
}
