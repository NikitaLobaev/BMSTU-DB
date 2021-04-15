package delivery

import (
	"../usecase"
	"github.com/labstack/echo/v4"
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

}
