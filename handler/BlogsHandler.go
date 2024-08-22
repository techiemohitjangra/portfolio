package handler

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	model "github.com/techiemohitjangra/portfolio/model"
	pages "github.com/techiemohitjangra/portfolio/view/pages"
)

type BlogsHandler struct {
	DB  *sql.DB
	Env map[string]string
}

func (handler BlogsHandler) HandleBlogShow(ctx echo.Context) error {
	if handler.DB == nil {
		db, err := model.OpenDB(handler.Env["DB"])
		if err != nil {
			log.Println("failed to open db file: ", err)
			return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error_message": "database not available",
			})
		}
		handler.DB = db
	}

	blogID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		log.Println("failed to parse blog ID: ", err)
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"error_message": "Invalid ID",
		})
	}

	blog, err := model.GetBlog(handler.DB, blogID)
	if err != nil {
		log.Println("failed to get Blog for given id: ", err)
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"error_message": "Failed to fetch blog",
		})
	}

	acceptHeader := ctx.Request().Header.Get("Accept")
	if strings.Contains(acceptHeader, "application/json") {
		return ctx.JSON(http.StatusOK, blog)
	} else if strings.Contains(acceptHeader, "text/html") {
		return render(ctx, pages.BlogPage(blog))
	}
	return ctx.JSON(http.StatusNotAcceptable, map[string]interface{}{
		"error_message": "Unsupported format",
	})
}

func (handler BlogsHandler) HandleBlogsShow(ctx echo.Context) error {
	if handler.DB == nil {
		db, err := model.OpenDB(handler.Env["DB"])
		if err != nil {
			log.Println("failed to open db file: ", err)
			return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error_message": "database not available",
			})
		}
		handler.DB = db
	}

	blogs, err := model.GetBlogs(handler.DB)
	if err != nil {
		log.Println("failed to fetch list of blogs: ", err)
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"error_message": "Failed to fetch list of blogs",
		})
	}

	acceptHeader := ctx.Request().Header.Get("Accept")
	if strings.Contains(acceptHeader, "application/json") {
		return ctx.JSON(http.StatusOK, blogs)
	} else if strings.Contains(acceptHeader, "text/html") {
		return render(ctx, pages.BlogsPage(blogs))
	}
	return ctx.String(http.StatusNotAcceptable, "Unsupported format")
}

func (handler BlogsHandler) HandleBlogCreate(ctx echo.Context) error {
	if handler.DB == nil {
		db, err := model.OpenDB(handler.Env["DB"])
		if err != nil {
			log.Println("failed to open db file: ", err)
			return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error_message": "database not available",
			})
		}
		handler.DB = db
	}

	var blog model.Blog
	err := ctx.Bind(blog)
	if err != nil {
		log.Println("failed to get Form Params from request context(ctx)")
		return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error_message": "failed to parse blog from the request",
		})
	}

	blog.PublishedOn = time.Now()
	blog.LastUpdated = time.Now()

	blogID, err := model.AddBlog(handler.DB, blog)
	if err != nil {
		log.Println("failed to add blog")
		return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error_message": "failed to added blog",
		})
	}
	// TODO: add support for different accept Headers
	log.Println(blogID)

	return ctx.NoContent(http.StatusOK)
}

func (handler BlogsHandler) HandleBlogUpdate(ctx echo.Context) error {
	if handler.DB == nil {
		db, err := model.OpenDB(handler.Env["DB"])
		if err != nil {
			log.Println("failed to open db file: ", err)
			return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error_message": "database not available",
			})
		}
		handler.DB = db
	}

	var blog model.Blog
	err := ctx.Bind(blog)
	if err != nil {
		log.Println("failed to get Form Params from request context(ctx)")
		return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error_message": "failed to parse blog from the request",
		})
	}

	blog.LastUpdated = time.Now()
	err = model.UpdateBlog(handler.DB, blog)

	return nil
}

func (handler BlogsHandler) HandleBlogDelete(ctx echo.Context) error {
	if handler.DB == nil {
		db, err := model.OpenDB(handler.Env["DB"])
		if err != nil {
			log.Println("failed to open db file: ", err)
			return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error_message": "database not available",
			})
		}
		handler.DB = db
	}

	blogID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		log.Println("failed to parse path parameter for HandlerBlogDelete: ", err)
		return ctx.String(http.StatusNotFound, "Invalid Blog ID")
	}

	err = model.DeleteBlog(handler.DB, blogID)
	if err != nil {
		log.Println("failed to delete Blog")
		return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error_message": "Failed to delete Blog",
		})
	}

	acceptHeader := ctx.Request().Header.Get("Accept")
	if strings.Contains(acceptHeader, "application/json") {
		return ctx.JSON(http.StatusOK, map[string]interface{}{
			"message": "blog deleted Successfully",
		})
	} else if strings.Contains(acceptHeader, "text/html") {
		return ctx.JSON(http.StatusOK, map[string]interface{}{
			"message": "blog deleted Successfully",
		})
	}
	return ctx.String(http.StatusNotAcceptable, "Unsupported format")
}
