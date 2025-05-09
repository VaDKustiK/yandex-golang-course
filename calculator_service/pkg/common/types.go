package common

type CalcRequest struct {
	Expression string `json:"expression"`
}

type CalcResponse struct {
	Result float64 `json:"result,omitempty"`
	Error  string  `json:"error,omitempty"`
}

type Expression struct {
	ID         uint     `gorm:"primaryKey" json:"id"`
	UserID     uint     `gorm:"index;not null" json:"-"`
	Expression string   `gorm:"not null" json:"expression"`
	Status     string   `gorm:"not null" json:"status"`
	Result     *float64 `json:"result,omitempty"`
	Tasks      []Task   `gorm:"foreignKey:ExprID"`
}

// type Task struct {
// 	gorm.Model
// 	ExprID        uint     `json:"expr_id"`
// 	Expression    string   `json:"expression"`
// 	Arg1          float64  `json:"arg1"`
// 	Arg2          float64  `json:"arg2"`
// 	Operation     string   `json:"operation"`
// 	Result        *float64 `json:"result,omitempty"`
// 	Status        string   `json:"status"`
// 	OperationTime int      `json:"operation_time"`
// }

type Task struct {
	ID            uint     `gorm:"primaryKey" json:"id"`
	ExprID        uint     `gorm:"index;not null" json:"expr_id"`
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
	ID       uint   `gorm:"primaryKey"`
	Login    string `gorm:"uniqueIndex;not null"`
	Password string `gorm:"not null"`
}
