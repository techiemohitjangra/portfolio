package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/techiemohitjangra/portfolio/view/pages"
)

type HomeHandler struct{}

func (h HomeHandler) HandleHomeShow(c echo.Context) error {
	return render(c, pages.HomePage())
}
