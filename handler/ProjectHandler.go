package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/techiemohitjangra/portfolio/view/pages"
)

type ProjectHandler struct{}

func (h ProjectHandler) HandleProjectShow(c echo.Context) error {
	return render(c, pages.ProjectPage())
}
