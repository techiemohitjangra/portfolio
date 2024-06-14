package main

import (
	"github.com/labstack/echo/v4"
	handler "github.com/techiemohitjangra/portfolio/handler"
)

func SetupRoutes(router *echo.Echo) {
	homeHandler := handler.HomeHandler{}
	router.GET("/", homeHandler.HandleHomeShow)
	router.POST("/contact", homeHandler.HandleContact)

	aboutHandler := handler.AboutHandler{}
	router.GET("/about", aboutHandler.HandleAboutShow)

	resumeHandler := handler.ResumeHandler{}
	router.GET("/resume", resumeHandler.HandleResumeShow)

	projectsHandler := handler.ProjectsHandler{}
	router.GET("/projects", projectsHandler.HandleProjectsShow)
	router.GET("/project/:id", projectsHandler.HandleProjectShow)
	router.POST("/project", projectsHandler.HandleProjectCreate)
	router.PATCH("/project/:id", projectsHandler.HandleProjectUpdate)
	// router.PATCH("/project/title/:id", projectsHandler.HandleProjectsShow)
	// router.PATCH("/project/subtitle/:id", projectsHandler.HandleProjectsShow)
	// router.PATCH("/project/content/:id", projectsHandler.HandleProjectsShow)
	router.DELETE("/project/:id", projectsHandler.HandleProjectDelete)

	blogsHandler := handler.BlogsHandler{}
	router.GET("/blogs", blogsHandler.HandleBlogsShow)
	router.GET("/blog/:id", blogsHandler.HandleBlogShow)
	router.POST("/blog", blogsHandler.HandleBlogCreate)
	router.PATCH("/project/:id", blogsHandler.HandleBlogUpdate)
	router.DELETE("/project/:id", blogsHandler.HandleBlogDelete)

	userHandler := handler.UserHandler{}
	router.POST("/register", userHandler.HandleRegister)
	router.POST("/login", userHandler.HandleLogin)
	router.POST("/logout", userHandler.HandleLogout)
}
