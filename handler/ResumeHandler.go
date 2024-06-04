package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/techiemohitjangra/portfolio/view/pages"
)

type ResumeHandler struct{}

func (h ResumeHandler) HandleResumeShow(c echo.Context) error {
	return render(c, pages.ResumePage())
}
