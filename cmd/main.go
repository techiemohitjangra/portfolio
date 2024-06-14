package main

import (
	"log/slog"

	"github.com/labstack/echo/v4"
)

func main() {
	router := echo.New()
	router.Static("/static", "static")

	SetupRoutes(router)

	listenAddr := ":3000"
	slog.Info("HTTP server started", "listenAddr", listenAddr)
	router.Start(listenAddr)
}
