// @title           Subscriptions API
// @version         1.0
// @description     API для управления подписками.
// @host      localhost:8080
// @BasePath  /
package main

import (
	"log"
	"net/http"
	"os"
	"subscriptions/internal/handlers"
	"subscriptions/internal/router"
	"subscriptions/internal/service"
	"subscriptions/internal/storage"
	_ "subscriptions/internal/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	
	ports := os.Getenv("LISTEN_AND_SERVE_PORTS")

	storage := storage.NewStorage()

	log.Print("connected to db")

	storage.RunMigrations()

	log.Print("migrations accepted")

	mux := http.NewServeMux()

	s := service.NewService(storage)

	w := service.StartWorkerPool(6, s)

	h := handlers.NewHandler(w)

	router := router.NewRouter(h)
	router.InitRoutes(mux)
	wrapped := router.WrapMiddle(mux)

	log.Printf("listening on %v", ports)
	if err := http.ListenAndServe(ports, wrapped); err != nil {
    log.Fatalf("HTTP server failed: %v", err)
}

}
