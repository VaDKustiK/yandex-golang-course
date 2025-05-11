package main

import (
	"log"
	"net"
	"net/http"
	"os"

	"github.com/VaDKustiK/yandex-golang-course/calculator_service/pkg/common"
	"github.com/VaDKustiK/yandex-golang-course/calculator_service/pkg/orchestrator"
	pb "github.com/VaDKustiK/yandex-golang-course/calculator_service/pkg/orchestrator/rpc"
	"github.com/glebarez/sqlite"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

func main() {
	dsn := "file:calc.db?_pragma=foreign_keys(1)"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to open DB: %v", err)
	}
	if err := db.AutoMigrate(&common.User{}, &common.Expression{}, &common.Task{}); err != nil {
		log.Fatalf("auto-migrate failed: %v", err)
	}
	orchestrator.SetDB(db)

	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/register", orchestrator.RegisterHandler)
	mux.HandleFunc("/api/v1/login", orchestrator.LoginHandler)

	protected := func(h http.HandlerFunc) http.Handler {
		return orchestrator.AuthMiddleware(h)
	}
	mux.Handle("/api/v1/calculate", protected(orchestrator.AddExpressionHandler))
	mux.Handle("/api/v1/expressions", protected(orchestrator.ListExpressionsHandler))
	mux.Handle("/api/v1/expressions/", protected(orchestrator.GetExpressionHandler))

	mux.HandleFunc("/internal/task", orchestrator.InternalTaskHandler)

	go func() {
		addr := ":" + getenv("PORT", "8080")
		log.Printf("HTTP server listening on %s", addr)
		log.Fatal(http.ListenAndServe(addr, mux))
	}()

	lis, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatalf("gRPC listen failed: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterOrchestratorServer(grpcServer, orchestrator.NewGRPCServer())

	log.Println("gRPC server listening on :9090")
	log.Fatal(grpcServer.Serve(lis))
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
