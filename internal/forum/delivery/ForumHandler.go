package delivery

import (
	"github.com/NikitaLobaev/BMSTU-DB/internal/forum/usecase"
	"github.com/NikitaLobaev/BMSTU-DB/internal/models"
	. "github.com/NikitaLobaev/BMSTU-DB/internal/tools/responser"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"time"
)

type ForumHandler struct {
	forumUsecase *usecase.ForumUsecase
}

func NewForumHandler(forumUsecase *usecase.ForumUsecase) *ForumHandler {
	return &ForumHandler{
		forumUsecase: forumUsecase,
	}
}

func (forumHandler *ForumHandler) Configure(echoWS *echo.Echo) {
	echoWS.POST("/api/forum/create", forumHandler.HandlerForumCreate())
	echoWS.GET("/api/forum/:slug/details", forumHandler.HandlerForumGetOne())
	echoWS.POST("/api/forum/:slug/create", forumHandler.HandlerThreadCreate())
	echoWS.GET("/api/forum/:slug/users", forumHandler.HandlerForumGetUsers())
	echoWS.GET("/api/forum/:slug/threads", forumHandler.HandlerForumGetThreads())
}

func (forumHandler *ForumHandler) HandlerForumCreate() echo.HandlerFunc {
	return func(context echo.Context) error {
		forum := new(models.Forum)
		if err := context.Bind(forum); err != nil {
			return context.NoContent(http.StatusServiceUnavailable)
		}

		return Respond(context, forumHandler.forumUsecase.Create(forum))
	}
}

func (forumHandler *ForumHandler) HandlerForumGetOne() echo.HandlerFunc {
	return func(context echo.Context) error {
		slug := context.Param("slug")
		return Respond(context, forumHandler.forumUsecase.GetBySlug(slug))
	}
}

func (forumHandler *ForumHandler) HandlerThreadCreate() echo.HandlerFunc {
	return func(context echo.Context) error {
		thread := new(models.Thread)
		if err := context.Bind(thread); err != nil {
			return context.NoContent(http.StatusServiceUnavailable)
		}
		thread.ForumSlug = context.Param("slug") //TODO: здесь раньше была ошибка, сейчас тоже может быть, работало только при slug_

		return Respond(context, forumHandler.forumUsecase.ThreadCreate(thread))
	}
}

func (forumHandler *ForumHandler) HandlerForumGetThreads() echo.HandlerFunc {
	return func(context echo.Context) error {
		slug := context.Param("slug")

		forumParams := new(models.ForumParams)
		if limit, err := strconv.ParseUint(context.QueryParam("limit"), 10, 32); err == nil {
			forumParams.SetLimit(uint32(limit))
		}
		if since, err := time.Parse(time.RFC3339, context.QueryParam("since")); err == nil {
			forumParams.SetSince(since)
		}
		if desc, err := strconv.ParseBool(context.QueryParam("desc")); err == nil {
			forumParams.SetDesc(desc)
		}

		return Respond(context, forumHandler.forumUsecase.GetThreadsBySlug(slug, forumParams))
	}
}

func (forumHandler *ForumHandler) HandlerForumGetUsers() echo.HandlerFunc {
	return func(context echo.Context) error {
		slug := context.Param("slug")

		userParams := new(models.UserParams)
		if limit, err := strconv.ParseUint(context.QueryParam("limit"), 10, 32); err == nil {
			userParams.SetLimit(uint32(limit))
		}
		if since := context.QueryParam("since"); since != "" {
			userParams.SetSince(since)
		}
		if desc, err := strconv.ParseBool(context.QueryParam("desc")); err == nil {
			userParams.SetDesc(desc)
		}

		return Respond(context, forumHandler.forumUsecase.GetUsersBySlug(slug, userParams))
	}
}
