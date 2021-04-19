package delivery

import (
	"../../models"
	. "../../tools/responser"
	"../usecase"
	"github.com/labstack/echo/v4"
	"net/http"
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
	echoWS.POST("/api/forum/:slug/create", nil)
	echoWS.GET("/api/forum/:slug/users", nil)
	echoWS.GET("/api/forum/:slug/threads", nil)
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
