package main

import (
	"log"
	"net/http"
	"os"

	"github.com/VaDKustiK/yandex-golang-course/calculator_service/pkg/agent"
	"github.com/VaDKustiK/yandex-golang-course/calculator_service/pkg/orchestrator"
)

func main() {
	go func() {
		mux := http.NewServeMux()

		mux.HandleFunc("/api/v1/register", orchestrator.RegisterHandler)
		mux.HandleFunc("/api/v1/login", orchestrator.LoginHandler)

		mux.Handle("/api/v1/expressions", orchestrator.AuthMiddleware(http.HandlerFunc(orchestrator.ListExpressionsHandler)))
		mux.Handle("/api/v1/expressions/", orchestrator.AuthMiddleware(http.HandlerFunc(orchestrator.GetExpressionHandler)))
		mux.Handle("/api/v1/calculate", orchestrator.AuthMiddleware(http.HandlerFunc(agent.CalculateHandler)))
		mux.Handle("/internal/task", orchestrator.AuthMiddleware(http.HandlerFunc(orchestrator.InternalTaskHandler)))

		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}
		log.Printf("Сервер запущен на http://localhost:%s\n", port)
		if err := http.ListenAndServe(":"+port, mux); err != nil {
			log.Fatal("Ошибка при запуске сервера:", err)
		}
	}()

	go agent.RunAgent()

	select {}
}
