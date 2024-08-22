package handler

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	model "github.com/techiemohitjangra/portfolio/model"
	pages "github.com/techiemohitjangra/portfolio/view/pages"
)

type ProjectsHandler struct {
	DB  *sql.DB
	Env map[string]string
}

// TODO: add authentication check
func (handler ProjectsHandler) HandleProjectShow(ctx echo.Context) error {
	if handler.DB == nil {
		db, err := model.OpenDB(handler.Env["DB"])
		if err != nil {
			log.Println("failed to openDB: ", err)
			return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error_message": "database not available",
			})
		}
		handler.DB = db
	}

	projectID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		log.Println("failed to parse project ID from path parameter: ", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error_message": "failed to parse project ID from path parameter",
		})
	}

	project, err := model.GetProject(handler.DB, projectID)
	if err != nil {
		log.Println("failed to get project from ID: ", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error_message": "failed to get project from ID",
		})
	}

	acceptHeader := ctx.Request().Header.Get("Accept")
	if strings.Contains(acceptHeader, "application/json") {
		return ctx.JSON(http.StatusOK, project)
	} else if strings.Contains(acceptHeader, "text/html") {
		return render(ctx, pages.ProjectPage(project))
	}
	return ctx.JSON(http.StatusNotAcceptable, map[string]interface{}{
		"error_message": "Unsupported format",
	})
}

// TODO: add authentication check
func (handler ProjectsHandler) HandleProjectsShow(ctx echo.Context) error {
	if handler.DB == nil {
		db, err := model.OpenDB(handler.Env["DB"])
		if err != nil {
			log.Println("failed to open DB: ", err)
			return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error_message": "database not available",
			})
		}
		handler.DB = db
	}

	projects, err := model.GetProjects(handler.DB)
	if err != nil {
		log.Println("failed to get projects: ", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error_message": "failed to get projects",
		})
	}

	acceptHeader := ctx.Request().Header.Get("Accept")
	if strings.Contains(acceptHeader, "application/json") {
		return ctx.JSON(http.StatusOK, projects)
	} else if strings.Contains(acceptHeader, "text/html") {
		return render(ctx, pages.ProjectsPage(projects))
	}
	return ctx.JSON(http.StatusNotAcceptable, map[string]interface{}{
		"error_message": "Unsupported format",
	})
}

// TODO: add authentication check
func (handler ProjectsHandler) HandleProjectCreate(ctx echo.Context) error {
	if handler.DB == nil {
		db, err := model.OpenDB(handler.Env["DB"])
		if err != nil {
			log.Println("failed to open DB: ", err)
			return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error_message": "database not available",
			})
		}
		handler.DB = db
	}

	var project model.Project
	err := ctx.Bind(project)
	if err != nil {
		log.Println("failed to bind project to create project")
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"error_message": "failed to parse data",
		})
	}

	acceptHeader := ctx.Request().Header.Get("Accept")
	if strings.Contains(acceptHeader, "application/json") {
		_, err := model.AddProject(handler.DB, project)
		if err != nil {
			log.Println("failed to add project: ", err)
			return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error_message": "failed to add project",
			})
		}
		return ctx.JSON(http.StatusCreated, map[string]interface{}{
			"message": "project created successfully.",
		})
	} else if strings.Contains(acceptHeader, "text/html") {
		projectID, err := model.AddProject(handler.DB, project)
		if err != nil {
			log.Println("failed to add project: ", err)
			return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error_message": "failed to add project",
			})
		}

		return ctx.Redirect(http.StatusCreated, fmt.Sprintf("/project/%d", projectID))
	}

	return ctx.JSON(http.StatusNotAcceptable, map[string]interface{}{
		"error_message": "Unsupported format",
	})
}

// TODO: add authentication check
func (handler ProjectsHandler) HandleProjectUpdate(ctx echo.Context) error {
	if handler.DB == nil {
		db, err := model.OpenDB(handler.Env["DB"])
		if err != nil {
			log.Println("failed to open DB: ", err)
			return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error_message": "database not available",
			})
		}
		handler.DB = db
	}

	var project model.Project
	err := ctx.Bind(project)
	if err != nil {
		log.Println("failed to bind project to create project")
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"error_message": "failed to parse data",
		})
	}

	acceptHeader := ctx.Request().Header.Get("Accept")
	if strings.Contains(acceptHeader, "application/json") {
		_, err := model.UpdateProject(handler.DB, project)
		if err != nil {
			log.Println("failed to add project: ", err)
			return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error_message": "failed to add project",
			})
		}
		return ctx.JSON(http.StatusOK, map[string]interface{}{
			"message": "project created successfully.",
		})
	} else if strings.Contains(acceptHeader, "text/html") {
		projectID, err := model.UpdateProject(handler.DB, project)
		if err != nil {
			log.Println("failed to add project: ", err)
			return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error_message": "failed to add project",
			})
		}

		return ctx.Redirect(http.StatusCreated, fmt.Sprintf("/project/%d", projectID))
	}

	return ctx.JSON(http.StatusNotAcceptable, map[string]interface{}{
		"error_message": "Unsupported format",
	})
}

// TODO: add authentication check
func (handler ProjectsHandler) HandleProjectDelete(ctx echo.Context) error {
	if handler.DB == nil {
		db, err := model.OpenDB(handler.Env["DB"])
		if err != nil {
			log.Println("failed to open DB: ", err)
			return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error_message": "database not available",
			})
		}
		handler.DB = db
	}

	projectID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		log.Println("failed to project ID.")
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"error_message": "failed to parse project ID.",
		})
	}

	acceptHeader := ctx.Request().Header.Get("Accept")
	if strings.Contains(acceptHeader, "application/json") {
		err := model.DeleteProject(handler.DB, projectID)
		if err != nil {
			log.Println("failed to add project: ", err)
			return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error_message": "failed to add project",
			})
		}
		return ctx.JSON(http.StatusOK, map[string]interface{}{
			"message": "project created successfully.",
		})
	} else if strings.Contains(acceptHeader, "text/html") {
		err := model.DeleteProject(handler.DB, projectID)
		if err != nil {
			log.Println("failed to add project: ", err)
			return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error_message": "failed to add project",
			})
		}

		return ctx.Redirect(http.StatusOK, fmt.Sprintf("/projects/"))
	}

	return ctx.JSON(http.StatusNotAcceptable, map[string]interface{}{
		"error_message": "Unsupported format",
	})
}
