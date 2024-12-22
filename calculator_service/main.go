package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func CalculateHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if rec := recover(); rec != nil {
			http.Error(w, `{"error":"Internal server error"}`, http.StatusInternalServerError)
			log.Printf("panic occurred: %v", rec)
		}
	}()

	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var req CalcRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	result, err := Calc(req.Expression)
	if err != nil {
		http.Error(w, `{"error":"expression is not valid"}`, http.StatusUnprocessableEntity)
		return
	}

	resp := CalcResponse{Result: result}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func main() {
	http.HandleFunc("/api/v1/calculate", CalculateHandler)
	log.Println("server is running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Error starting server:", err)
	}
}
