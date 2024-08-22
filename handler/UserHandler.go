package handler

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	model "github.com/techiemohitjangra/portfolio/model"
)

type UserHandler struct {
	DB  *sql.DB
	Env map[string]string
}

func (handler UserHandler) HandleRegister(ctx echo.Context) error {
	if handler.DB == nil {
		db, err := model.OpenDB(handler.Env["DB"])
		if err != nil {
			log.Println("failed to open DB file: ", err)
			return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error_message": "database not available",
			})
		}
		handler.DB = db
	}
	var user model.User
	err := ctx.Bind(user)
	if err != nil {
		log.Println("failed to bind user to create user")
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"error_message": "failed to parse user data",
		})
	}

	err = model.AddUser(handler.DB, user)
	if err != nil {
		log.Println("failed to add user")
		return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error_message": "failed to add user",
		})
	}
	return nil
}

func (handler UserHandler) HandleLogin(ctx echo.Context) error {
	if handler.DB == nil {
		db, err := model.OpenDB(handler.Env["DB"])
		if err != nil {
			log.Println("failed to open DB file: ", err)
			return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error_message": "database not available",
			})
		}
		handler.DB = db
	}
	return nil
}

func (handler UserHandler) HandleLogout(ctx echo.Context) error {
	if handler.DB == nil {
		db, err := model.OpenDB(handler.Env["DB"])
		if err != nil {
			log.Println("failed to open DB file: ", err)
			return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error_message": "database not available",
			})
		}
		handler.DB = db
	}
	return nil
}
