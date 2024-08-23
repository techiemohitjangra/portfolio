package handler

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	model "github.com/techiemohitjangra/portfolio/model"
	sample "github.com/techiemohitjangra/portfolio/sample"
)

type TestSetupHandler struct {
	DB  *sql.DB
	Env map[string]string
}

func (handler TestSetupHandler) SetupHandler(ctx echo.Context) error {
	if handler.DB == nil {
		db, err := model.OpenDB(handler.Env["DB"])
		if err != nil {
			log.Println("failed to open test db: ", err)

			return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error_message": "database not available",
			})
		}
		handler.DB = db
	}

	err := model.InitTables(handler.DB)
	if err != nil {
		log.Println("failed to create tables: ", err)
	}

	err = model.AddUser(handler.DB, sample.SampleUser)
	if err != nil {
		log.Println("failed to add sample user: ", err)
	}
	log.Println("successfully added user")

	blogID, err := model.AddBlog(handler.DB, sample.SampleBlog)
	if err != nil {
		log.Println(blogID)
		log.Println("failed to add sample blog: ", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error_message": "failed to add blog in setup",
		})
	}
	log.Println("successfully added blog")

	projectID, err := model.AddProject(handler.DB, sample.SampleProject1)
	if err != nil {
		log.Println(projectID)
		log.Println("failed to add sample project1: ", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error_message": "failed to add project",
		})
	}
	log.Println("successfully added project")

	projectID_2, err := model.AddProject(handler.DB, sample.SampleProject2)
	if err != nil {
		log.Println(projectID_2)
		log.Println("failed to add sample project2: ", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error_message": "failed to add project",
		})
	}
	log.Println("successfully added project 2")

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"time": time.Now(),
	})
}
