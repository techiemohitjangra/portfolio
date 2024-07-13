package main

import (
	"log"
	"log/slog"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	router := echo.New()
	router.Static("/static", "static")

	env, err := godotenv.Read(".env")
	if err != nil {
		log.Fatal("failed to load .env")
	}

	SetupRoutes(router)

	listenAddr := ":" + env["PORT"]
	slog.Info("HTTP server started", "listenAddr", listenAddr)
	router.Start(listenAddr)
}
