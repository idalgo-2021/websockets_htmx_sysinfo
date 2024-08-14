package main

import (
	"log"
	"net/http"
	"time"
	"websockets_htmx_sysinfo/internal/server"
	"websockets_htmx_sysinfo/internal/service"
)

func main() {

	// Создание сервиса для сбора данных о системе
	hs := service.NewHardwareService()

	// Инициализация WS-сервера
	srv := server.NewServer(hs)

	// Запуск HTTP-сервера
	go func() {
		log.Println("Starting server on :8080")
		if err := http.ListenAndServe(":8080", srv.Router()); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Периодически публикуем информацию
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if err := srv.PublishSystemData(); err != nil {
			log.Printf("Error publishing system data: %v", err)
		}
	}
}
