package delivery

import (
	"../../models"
	. "../../tools/responser"
	"../usecase"
	"github.com/labstack/echo/v4"
	"net/http"
)

type PostHandler struct {
	postUsecase *usecase.PostUsecase
}

func NewPostHandler(postUsecase *usecase.PostUsecase) *PostHandler {
	return &PostHandler{
		postUsecase: postUsecase,
	}
}

func (postHandler *PostHandler) Configure(echoWS *echo.Echo) {
	echoWS.GET("/api/post/:id/details", nil)
	echoWS.POST("/api/post/:id/details", nil)
}

func (postHandler *PostHandler) HandlerUserCreate() echo.HandlerFunc {
	return func(context echo.Context) error {
		user := new(models.User)
		if err := context.Bind(user); err != nil {
			return context.NoContent(http.StatusServiceUnavailable)
		}
		user.Nickname = context.Param("nickname")
		return Respond(context, postHandler.postUsecase.Create(user))
	}
}

func (postHandler *PostHandler) HandlerUserGetOne() echo.HandlerFunc {
	return func(context echo.Context) error {
		nickname := context.Param("nickname")
		return Respond(context, postHandler.postUsecase.GetByNickname(nickname))
	}
}

func (postHandler *PostHandler) HandlerUpdate() echo.HandlerFunc {
	return func(context echo.Context) error {
		userUpdate := new(models.UserUpdate)
		if err := context.Bind(userUpdate); err != nil {
			return context.NoContent(http.StatusServiceUnavailable)
		}
		nickname := context.Param("nickname")
		return Respond(context, postHandler.postUsecase.Update(nickname, userUpdate))
	}
}
