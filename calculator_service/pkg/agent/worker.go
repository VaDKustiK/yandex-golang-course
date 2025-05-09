package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/VaDKustiK/yandex-golang-course/calculator_service/pkg/common"
)

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
	select {}
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
			resp.Body.Close()
			time.Sleep(2 * time.Second)
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

		result, err := sendToCalculatorService(task.Expression)
		if err != nil {
			log.Printf("Worker %d: ошибка вычисления выражения: %v", id, err)
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

func sendToCalculatorService(expression string) (float64, error) {
	reqData := map[string]string{
		"expression": expression,
	}

	buf, err := json.Marshal(reqData)
	if err != nil {
		return 0, fmt.Errorf("ошибка сериализации данных: %v", err)
	}

	resp, err := http.Post("http://localhost:8080/api/v1/calculate", "application/json", bytes.NewBuffer(buf))
	if err != nil {
		return 0, fmt.Errorf("ошибка отправки запроса в калькулятор: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("получен неожиданный статус от калькулятора: %v", resp.StatusCode)
	}

	var result struct {
		Result float64 `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("ошибка декодирования результата: %v", err)
	}

	return result.Result, nil
}
