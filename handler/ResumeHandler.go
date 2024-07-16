package handler

import (
	"database/sql"

	"github.com/labstack/echo/v4"
	pages "github.com/techiemohitjangra/portfolio/view/pages"
)

type ResumeHandler struct {
	DB  *sql.DB
	Env map[string]string
}

func (handler ResumeHandler) HandleResumeShow(ctx echo.Context) error {
	return render(ctx, pages.ResumePage())
}
