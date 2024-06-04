package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/techiemohitjangra/portfolio/view/pages"
)

type AboutHandler struct{}

func (h AboutHandler) HandleAboutShow(c echo.Context) error {
	return render(c, pages.AboutPage())
}
