package agent

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCalculateHandler(t *testing.T) {
	reqBody := []byte(`{"expression": "2+2*2"}`)
	req, err := http.NewRequest("POST", "/api/v1/calculate", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatalf("Ошибка при создании запроса: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CalculateHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Ожидался статус 200, но получили %d", status)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Ошибка при разборе JSON: %v", err)
	}

	if _, ok := resp["result"]; !ok {
		t.Errorf("Ответ должен содержать 'result'")
	}
}
