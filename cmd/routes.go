package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	handler "github.com/techiemohitjangra/portfolio/handler"
	model "github.com/techiemohitjangra/portfolio/model"
)

func SetupRoutes(router *echo.Echo) {
	env, err := godotenv.Read(".env")

	db, err := model.OpenDB(env["DB"])
	if err != nil {
		log.Println("failed to open test DB: ", err)
	}
	defer model.CloseDB(db)

	// handlers to setup the test environment and data
	setupHandler := handler.TestSetupHandler{
		DB:  db,
		Env: env,
	}
	router.GET("/setup", setupHandler.SetupHandler)

	// handlers for home page
	homeHandler := handler.HomeHandler{
		DB:  db,
		Env: env,
	}
	router.GET("/", homeHandler.HandleHomeShow)
	router.POST("/contact", homeHandler.HandleContact)

	// handlers for about page
	aboutHandler := handler.AboutHandler{
		DB:  db,
		Env: env,
	}
	router.GET("/about", aboutHandler.HandleAboutShow)

	// handlers for resume page
	resumeHandler := handler.ResumeHandler{
		DB:  db,
		Env: env,
	}
	router.GET("/resume", resumeHandler.HandleResumeShow)

	// handlers for projects
	projectsHandler := handler.ProjectsHandler{
		DB:  db,
		Env: env,
	}
	router.GET("/projects", projectsHandler.HandleProjectsShow)
	router.GET("/project/:id", projectsHandler.HandleProjectShow)
	router.POST("/project", projectsHandler.HandleProjectCreate)
	router.PATCH("/project/:id", projectsHandler.HandleProjectUpdate)
	// router.PATCH("/project/title/:id", projectsHandler.HandleProjectsShow)
	// router.PATCH("/project/subtitle/:id", projectsHandler.HandleProjectsShow)
	// router.PATCH("/project/content/:id", projectsHandler.HandleProjectsShow)
	router.DELETE("/project/:id", projectsHandler.HandleProjectDelete)

	blogsHandler := handler.BlogsHandler{
		DB:  db,
		Env: env,
	}
	router.GET("/blogs", blogsHandler.HandleBlogsShow)
	router.GET("/blog/:id", blogsHandler.HandleBlogShow)
	router.POST("/blog", blogsHandler.HandleBlogCreate)
	router.PATCH("/project/:id", blogsHandler.HandleBlogUpdate)
	router.DELETE("/project/:id", blogsHandler.HandleBlogDelete)

	userHandler := handler.UserHandler{
		DB:  db,
		Env: env,
	}
	router.POST("/register", userHandler.HandleRegister)
	router.POST("/login", userHandler.HandleLogin)
	router.POST("/logout", userHandler.HandleLogout)
}
