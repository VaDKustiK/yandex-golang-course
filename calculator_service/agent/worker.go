package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/VaDKustiK/yandex-golang-course/calculator_service/common"
)

// RunAgent запускает рабочих (воркеры) агента.
func RunAgent() {
	cpStr := os.Getenv("COMPUTING_POWER")
	cp, err := strconv.Atoi(cpStr)
	if err != nil || cp < 1 {
		cp = 1
	}
	log.Printf("Запуск агента с %d рабочими горутинами", cp)
	for i := 0; i < cp; i++ {
		go worker(i)
	}
	select {} // Бесконечное ожидание.
}

func worker(id int) {
	for {
		resp, err := http.Get("http://localhost:8080/internal/task")
		if err != nil {
			log.Printf("Worker %d: ошибка при получении задачи: %v", id, err)
			time.Sleep(2 * time.Second)
			continue
		}
		if resp.StatusCode == http.StatusNotFound {
			// Задач пока нет – ждём.
			time.Sleep(2 * time.Second)
			resp.Body.Close()
			continue
		}
		var data struct {
			Task common.Task `json:"task"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			log.Printf("Worker %d: ошибка декодирования задачи: %v", id, err)
			resp.Body.Close()
			continue
		}
		resp.Body.Close()
		task := data.Task
		log.Printf("Worker %d: получена задача %+v", id, task)
		// Имитация задержки выполнения операции.
		time.Sleep(time.Duration(task.OperationTime) * time.Millisecond)
		result, err := compute(task.Arg1, task.Arg2, task.Operation)
		if err != nil {
			log.Printf("Worker %d: ошибка вычисления: %v", id, err)
			continue
		}
		reqData := common.TaskResultRequest{
			ID:     task.ID,
			Result: result,
		}
		buf, _ := json.Marshal(reqData)
		resp2, err := http.Post("http://localhost:8080/internal/task", "application/json", bytes.NewBuffer(buf))
		if err != nil {
			log.Printf("Worker %d: ошибка отправки результата: %v", id, err)
			continue
		}
		resp2.Body.Close()
		log.Printf("Worker %d: задача %d выполнена, результат: %f", id, task.ID, result)
	}
}

func compute(arg1, arg2 float64, op string) (float64, error) {
	switch op {
	case "+":
		return arg1 + arg2, nil
	case "-":
		return arg1 - arg2, nil
	case "*":
		return arg1 * arg2, nil
	case "/":
		if arg2 == 0 {
			return 0, errors.New("division by zero")
		}
		return arg1 / arg2, nil
	default:
		return 0, errors.New("unknown operator")
	}
}
