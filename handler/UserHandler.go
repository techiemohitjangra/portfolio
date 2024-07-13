package handler

import (
	"database/sql"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	DB  *sql.DB
	Env map[string]string
}

func (handler UserHandler) HandleRegister(ctx echo.Context) error {
	return nil
}

func (handler UserHandler) HandleLogin(ctx echo.Context) error {
	return nil
}

func (handler UserHandler) HandleLogout(ctx echo.Context) error {
	return nil
}
