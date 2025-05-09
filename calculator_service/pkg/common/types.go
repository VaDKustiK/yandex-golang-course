package common

import (
	"gorm.io/gorm"
)

type CalcRequest struct {
	Expression string `json:"expression"`
}

type CalcResponse struct {
	Result float64 `json:"result,omitempty"`
	Error  string  `json:"error,omitempty"`
}

type Expression struct {
	gorm.Model
	Expression string   `json:"expression"`
	Status     string   `json:"status"`
	Result     *float64 `json:"result,omitempty"`
	UserID     uint     `json:"-"`
	Tasks      []Task   `gorm:"foreignKey:ExprID"`
}

type Task struct {
	gorm.Model
	ExprID        uint     `json:"expr_id"`
	Expression    string   `json:"expression"`
	Arg1          float64  `json:"arg1"`
	Arg2          float64  `json:"arg2"`
	Operation     string   `json:"operation"`
	Result        *float64 `json:"result,omitempty"`
	Status        string   `json:"status"`
	OperationTime int      `json:"operation_time"`
}

type TaskResultRequest struct {
	ID     uint    `json:"id"`
	Result float64 `json:"result"`
}

type User struct {
	gorm.Model
	Login    string `gorm:"uniqueIndex"`
	Password string
}
