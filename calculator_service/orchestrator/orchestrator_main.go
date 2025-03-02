package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	//эндпоинты  в handlers.go.
	http.HandleFunc("/api/v1/calculate", AddExpressionHandler)
	http.HandleFunc("/api/v1/expressions", ListExpressionsHandler)
	http.HandleFunc("/api/v1/expressions/", GetExpressionHandler)
	http.HandleFunc("/internal/task", InternalTaskHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Оркестратор запущен на http://localhost:%s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Ошибка при запуске сервера:", err)
	}
}
