package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/techiemohitjangra/portfolio/view/pages"
)

type BlogsHandler struct{}

func (h BlogsHandler) HandleBlogsShow(c echo.Context) error {
	return render(c, pages.BlogsPage())
}
