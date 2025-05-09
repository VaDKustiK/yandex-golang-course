package main

import (
	"log"
	"net/http"
	"os"

	"github.com/VaDKustiK/yandex-golang-course/calculator_service/pkg/common"
	"github.com/VaDKustiK/yandex-golang-course/calculator_service/pkg/orchestrator"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	http.HandleFunc("/api/v1/calculate", orchestrator.AddExpressionHandler)
	http.HandleFunc("/api/v1/expressions", orchestrator.ListExpressionsHandler)
	http.HandleFunc("/api/v1/expressions/", orchestrator.GetExpressionHandler)
	http.HandleFunc("/internal/task", orchestrator.InternalTaskHandler)

	db, err := gorm.Open(sqlite.Open("calc.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&common.User{}, &common.Expression{}, &common.Task{})
	orchestrator.SetDB(db)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Оркестратор запущен на http://localhost:%s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Ошибка при запуске сервера:", err)
	}
}
