package handler

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	model "github.com/techiemohitjangra/portfolio/model"
	pages "github.com/techiemohitjangra/portfolio/view/pages"
)

type AboutHandler struct {
	DB  *sql.DB
	Env map[string]string
}

func (handler AboutHandler) HandleAboutShow(ctx echo.Context) error {
	if handler.DB == nil {
		db, err := model.OpenDB(handler.Env["DB"])
		if err != nil {
			log.Println("failed to open DB file: ", err)
			return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error_message": "database not available",
			})
		}
		handler.DB = db
	}

	user, err := model.GetUserShow(handler.DB, handler.Env["DEFAULT_USER"])
	if err != nil {
		log.Println("failed to fetch user: ", err)
		return echo.NewHTTPError(http.StatusBadRequest, map[string]interface{}{
			"error_message": "failed to fetch user",
		})
	}
	return render(ctx, pages.AboutPage(user))
}
