package main

import (
	"log"
	"net/http"
	"os"

	"github.com/VaDKustiK/yandex-golang-course/calculator_service/pkg/orchestrator"
)

func main() {
	http.HandleFunc("/api/v1/calculate", orchestrator.AddExpressionHandler)
	http.HandleFunc("/api/v1/expressions", orchestrator.ListExpressionsHandler)
	http.HandleFunc("/api/v1/expressions/", orchestrator.GetExpressionHandler)
	http.HandleFunc("/internal/task", orchestrator.InternalTaskHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Оркестратор запущен на http://localhost:%s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Ошибка при запуске сервера:", err)
	}
}
