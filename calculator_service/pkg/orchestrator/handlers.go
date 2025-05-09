package orchestrator

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/VaDKustiK/yandex-golang-course/calculator_service/pkg/common"
)

var (
	expressionsMu sync.Mutex
	tasksMu       sync.Mutex

	expressions = make(map[int]*common.Expression)
	tasks       = make(map[int]*common.Task)

	nextExprID = 1
	nextTaskID = 1
)

func AddExpressionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	var req common.CalcRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "invalid request body"}`, http.StatusBadRequest)
		return
	}
	exprStr := strings.TrimSpace(req.Expression)
	if exprStr == "" {
		http.Error(w, `{"error": "empty expression"}`, http.StatusUnprocessableEntity)
		return
	}

	expressionsMu.Lock()
	exprID := nextExprID
	nextExprID++
	expr := &common.Expression{
		ID:         exprID,
		Expression: exprStr,
		Status:     "pending",
		TaskIDs:    []int{},
	}
	expressions[exprID] = expr
	expressionsMu.Unlock()

	tokens := common.Tokenize(exprStr)
	if len(tokens) < 3 {
		http.Error(w, `{"error": "expression is too short"}`, http.StatusUnprocessableEntity)
		return
	}
	var numbers []float64
	var ops []string
	for i, token := range tokens {
		if i%2 == 0 {
			num, err := strconv.ParseFloat(token, 64)
			if err != nil {
				http.Error(w, `{"error": "invalid number in expression"}`, http.StatusUnprocessableEntity)
				return
			}
			numbers = append(numbers, num)
		} else {
			ops = append(ops, token)
		}
	}

	tasksMu.Lock()
	for i, op := range ops {
		opTime := getOperationTime(op)
		task := &common.Task{
			ID:            nextTaskID,
			ExprID:        exprID,
			Expression:    exprStr,
			Arg1:          numbers[i],
			Arg2:          numbers[i+1],
			Operation:     op,
			Status:        "pending",
			OperationTime: opTime,
		}
		nextTaskID++
		tasks[task.ID] = task
		expr.TaskIDs = append(expr.TaskIDs, task.ID)
	}
	tasksMu.Unlock()

	expressionsMu.Lock()
	expr.Status = "in_progress"
	expressionsMu.Unlock()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": exprID})
}

func ListExpressionsHandler(w http.ResponseWriter, r *http.Request) {
	expressionsMu.Lock()
	defer expressionsMu.Unlock()
	var list []common.Expression
	for _, expr := range expressions {
		if expr.Status != "completed" {
			allDone := true
			var final float64
			tasksMu.Lock()
			for _, tid := range expr.TaskIDs {
				task, ok := tasks[tid]
				if !ok || task.Status != "completed" || task.Result == nil {
					allDone = false
					break
				}
				final = *task.Result
			}
			tasksMu.Unlock()
			if allDone {
				expr.Status = "completed"
				expr.Result = &final
			}
		}
		list = append(list, *expr)
	}
	json.NewEncoder(w).Encode(map[string][]common.Expression{"expressions": list})
}

func GetExpressionHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/expressions/")
	exprID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error": "invalid id"}`, http.StatusBadRequest)
		return
	}
	expressionsMu.Lock()
	expr, ok := expressions[exprID]
	expressionsMu.Unlock()
	if !ok {
		http.Error(w, `{"error": "expression not found"}`, http.StatusNotFound)
		return
	}
	if expr.Status != "completed" {
		allDone := true
		var final float64
		tasksMu.Lock()
		for _, tid := range expr.TaskIDs {
			task, ok := tasks[tid]
			if !ok || task.Status != "completed" || task.Result == nil {
				allDone = false
				break
			}
			final = *task.Result
		}
		tasksMu.Unlock()
		if allDone {
			expr.Status = "completed"
			expr.Result = &final
		}
	}
	json.NewEncoder(w).Encode(map[string]*common.Expression{"expression": expr})
}

func InternalTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		getTask(w, r)
	} else if r.Method == http.MethodPost {
		postTaskResult(w, r)
	} else {
		http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

func getTask(w http.ResponseWriter, r *http.Request) {
	tasksMu.Lock()
	defer tasksMu.Unlock()
	for _, task := range tasks {
		if task.Status == "pending" {
			json.NewEncoder(w).Encode(map[string]*common.Task{"task": task})
			return
		}
	}
	http.Error(w, `{"error": "no task"}`, http.StatusNotFound)
}

func postTaskResult(w http.ResponseWriter, r *http.Request) {
	var req common.TaskResultRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "invalid request body"}`, http.StatusBadRequest)
		return
	}
	tasksMu.Lock()
	task, ok := tasks[req.ID]
	if !ok {
		tasksMu.Unlock()
		http.Error(w, `{"error": "task not found"}`, http.StatusNotFound)
		return
	}
	if task.Status != "pending" {
		tasksMu.Unlock()
		http.Error(w, `{"error": "task already completed"}`, http.StatusUnprocessableEntity)
		return
	}
	task.Result = &req.Result
	task.Status = "completed"
	tasksMu.Unlock()
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
	msStr := os.Getenv(envVar)
	ms, err := strconv.Atoi(msStr)
	if err != nil {
		return 100
	}
	return ms
}
