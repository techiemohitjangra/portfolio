package main

import (
	"log/slog"

	"github.com/labstack/echo/v4"
	handler "github.com/techiemohitjangra/portfolio/handler"
)

func main() {
	router := echo.New()

	homeHandler := handler.HomeHandler{}
	router.GET("/", homeHandler.HandleHomeShow)

	aboutHandler := handler.AboutHandler{}
	router.GET("/about", aboutHandler.HandleAboutShow)

	resumeHandler := handler.ResumeHandler{}
	router.GET("/resume", resumeHandler.HandleResumeShow)

	projectsHandler := handler.ProjectsHandler{}
	router.GET("/projects", projectsHandler.HandleProjectsShow)

	blogsHandler := handler.BlogsHandler{}
	router.GET("/blogs", blogsHandler.HandleBlogsShow)

	listenAddr := ":3000"
	slog.Info("HTTP server started", "listenAddr", listenAddr)
	router.Start(listenAddr)
}
