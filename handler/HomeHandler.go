package handler

import (
	"database/sql"

	"github.com/labstack/echo/v4"
	pages "github.com/techiemohitjangra/portfolio/view/pages"
)

type HomeHandler struct {
	DB  *sql.DB
	Env map[string]string
}

func (handler HomeHandler) HandleHomeShow(ctx echo.Context) error {
	return render(ctx, pages.HomePage())
}

func (handler HomeHandler) HandleContact(ctx echo.Context) error {
	return nil
}
