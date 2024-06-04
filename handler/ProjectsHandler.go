package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/techiemohitjangra/portfolio/view/pages"
)

type ProjectsHandler struct{}

func (h ProjectsHandler) HandleProjectsShow(c echo.Context) error {
	return render(c, pages.ProjectsPage())
}
