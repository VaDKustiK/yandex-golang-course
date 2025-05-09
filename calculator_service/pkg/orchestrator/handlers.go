// pkg/orchestrator/handlers.go
package orchestrator

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/VaDKustiK/yandex-golang-course/calculator_service/pkg/common"
	"gorm.io/gorm"
)

func AddExpressionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	var req common.CalcRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}
	exprStr := strings.TrimSpace(req.Expression)
	if exprStr == "" {
		http.Error(w, `{"error":"empty expression"}`, http.StatusUnprocessableEntity)
		return
	}

	db := GetDB()
	userID := UserIDFromContext(r.Context())
	expr := &common.Expression{
		Expression: exprStr,
		Status:     "pending",
		UserID:     userID,
	}
	if err := db.Create(expr).Error; err != nil {
		http.Error(w, `{"error":"cannot save expression"}`, http.StatusInternalServerError)
		return
	}

	tokens := common.Tokenize(exprStr)
	if len(tokens) < 3 {
		http.Error(w, `{"error":"expression is too short"}`, http.StatusUnprocessableEntity)
		return
	}

	var numbers []float64
	var ops []string
	for i, token := range tokens {
		if i%2 == 0 {
			num, err := strconv.ParseFloat(token, 64)
			if err != nil {
				http.Error(w, `{"error":"invalid number in expression"}`, http.StatusUnprocessableEntity)
				return
			}
			numbers = append(numbers, num)
		} else {
			ops = append(ops, token)
		}
	}

	for i, op := range ops {
		task := &common.Task{
			ExprID:        expr.ID,
			Expression:    exprStr,
			Arg1:          numbers[i],
			Arg2:          numbers[i+1],
			Operation:     op,
			Status:        "pending",
			OperationTime: getOperationTime(op),
		}
		if err := db.Create(task).Error; err != nil {
			http.Error(w, `{"error":"cannot save task"}`, http.StatusInternalServerError)
			return
		}
	}

	if err := db.Model(expr).Update("status", "in_progress").Error; err != nil {
		http.Error(w, `{"error":"cannot update status"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]uint{"id": expr.ID})
}

func ListExpressionsHandler(w http.ResponseWriter, r *http.Request) {
	db := GetDB()
	var exprs []common.Expression
	userID := UserIDFromContext(r.Context())
	if err := db.
		Where("user_id = ?", userID).
		Preload("Tasks").
		Find(&exprs).
		Error; err != nil {
		http.Error(w, `{"error":"cannot list expressions"}`, http.StatusInternalServerError)
		return
	}
	if err := db.Preload("Tasks").Find(&exprs).Error; err != nil {
		http.Error(w, `{"error":"cannot list expressions"}`, http.StatusInternalServerError)
		return
	}

	for i := range exprs {
		expr := &exprs[i]
		if expr.Status != "completed" {
			allDone := true
			var final float64
			for _, t := range expr.Tasks {
				if t.Status != "completed" || t.Result == nil {
					allDone = false
					break
				}
				final = *t.Result
			}
			if allDone {
				expr.Status = "completed"
				expr.Result = &final
				db.Model(expr).Updates(map[string]interface{}{
					"status": "completed",
					"result": final,
				})
			}
		}
	}

	json.NewEncoder(w).Encode(map[string][]common.Expression{"expressions": exprs})
}

func GetExpressionHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/expressions/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
		return
	}

	expr, err := getExpressionByID(uint(id))
	if err == gorm.ErrRecordNotFound {
		http.Error(w, `{"error":"expression not found"}`, http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
		return
	}

	if expr.Status != "completed" {
		allDone := true
		var final float64
		for _, t := range expr.Tasks {
			if t.Status != "completed" || t.Result == nil {
				allDone = false
				break
			}
			final = *t.Result
		}
		if allDone {
			expr.Status = "completed"
			expr.Result = &final
			GetDB().Model(expr).Updates(map[string]interface{}{
				"status": "completed",
				"result": final,
			})
		}
	}

	json.NewEncoder(w).Encode(map[string]*common.Expression{"expression": expr})
}

func InternalTaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getTask(w, r)
	case http.MethodPost:
		postTaskResult(w, r)
	default:
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

func getTask(w http.ResponseWriter, r *http.Request) {
	db := GetDB()
	var task common.Task
	if err := db.Where("status = ?", "pending").First(&task).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, `{"error":"no task"}`, http.StatusNotFound)
		} else {
			http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
		}
		return
	}
	json.NewEncoder(w).Encode(map[string]common.Task{"task": task})
}

func postTaskResult(w http.ResponseWriter, r *http.Request) {
	var req common.TaskResultRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}
	db := GetDB()
	var task common.Task
	if err := db.First(&task, req.ID).Error; err != nil {
		http.Error(w, `{"error":"task not found"}`, http.StatusNotFound)
		return
	}
	if task.Status != "pending" {
		http.Error(w, `{"error":"task already completed"}`, http.StatusUnprocessableEntity)
		return
	}
	if err := db.Model(&task).Updates(map[string]interface{}{
		"status": "completed",
		"result": req.Result,
	}).Error; err != nil {
		http.Error(w, `{"error":"cannot update task"}`, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func getOperationTime(op string) int {
	var envVar string
	switch op {
	case "+":
		envVar = "TIME_ADDITION_MS"
	case "-":
		envVar = "TIME_SUBTRACTION_MS"
	case "*":
		envVar = "TIME_MULTIPLICATIONS_MS"
	case "/":
		envVar = "TIME_DIVISIONS_MS"
	default:
		return 100
	}
	if ms, err := strconv.Atoi(os.Getenv(envVar)); err == nil {
		return ms
	}
	return 100
}
