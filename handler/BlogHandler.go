package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/techiemohitjangra/portfolio/view/pages"
)

type BlogHandler struct{}

func (h BlogHandler) HandleBlogShow(c echo.Context) error {
	return render(c, pages.BlogPage())
}
