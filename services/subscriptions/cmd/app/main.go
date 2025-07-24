package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"subscriptions/internal/handlers"
	"subscriptions/internal/router"
	"subscriptions/internal/service"
	"subscriptions/internal/storage"
)

func main() {
	log.Print("Привет")
	fmt.Print("ПРивет")
	ports := os.Getenv("LISTEN_AND_SERVE_PORTS")

	storage := storage.NewStorage()

	storage.RunMigrations()

	mux := http.NewServeMux()

	s := service.NewService(storage)

	w := service.StartWorkerPool(6, s)

	h := handlers.NewHandler(w)

	router := router.NewRouter(h)
	router.InitRoutes(mux)
	wrapped := router.WrapMiddle(mux)
	fmt.Print("aboba")
	log.Printf("Слушаю на %v", ports)
	if err := http.ListenAndServe(ports, wrapped); err != nil {
    log.Fatalf("HTTP server failed: %v", err)
}

}
