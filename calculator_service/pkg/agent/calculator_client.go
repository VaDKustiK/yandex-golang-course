package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func CalculateRemote(expression string) (float64, error) {
	data := map[string]string{"expression": expression}
	buf, _ := json.Marshal(data)

	resp, err := http.Post("http://localhost:8081/api/v1/calculate", "application/json", bytes.NewBuffer(buf))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("calculator returned %d", resp.StatusCode)
	}

	var r struct{ Result float64 }
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return 0, err
	}
	return r.Result, nil
}
