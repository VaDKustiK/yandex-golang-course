package main

import (
	"log"
	"net/http"
	"os"

	"github.com/VaDKustiK/yandex-golang-course/calculator_service/cmd/agent"
	"github.com/VaDKustiK/yandex-golang-course/calculator_service/cmd/orchestrator"
)

func main() {
	go func() {
		http.HandleFunc("/api/v1/calculate", agent.CalculateHandler)
		http.HandleFunc("/api/v1/expressions", orchestrator.ListExpressionsHandler)
		http.HandleFunc("/api/v1/expressions/", orchestrator.GetExpressionHandler)
		http.HandleFunc("/internal/task", orchestrator.InternalTaskHandler)

		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}
		log.Printf("Сервер запущен на http://localhost:%s\n", port)
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			log.Fatal("Ошибка при запуске сервера:", err)
		}
	}()

	go agent.RunAgent()

	select {}
}
