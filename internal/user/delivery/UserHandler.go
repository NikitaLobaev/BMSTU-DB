package delivery

import (
	"../../models"
	. "../../tools/responser"
	"../usecase"
	"github.com/labstack/echo/v4"
	"net/http"
)

type UserHandler struct {
	userUsecase *usecase.UserUsecase
}

func NewUserHandler(userUsecase *usecase.UserUsecase) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
	}
}

func (userHandler *UserHandler) Configure(echoWS *echo.Echo) {
	echoWS.POST("/api/user/:nickname/create", userHandler.HandlerUserCreate())
}

func (userHandler *UserHandler) HandlerUserCreate() echo.HandlerFunc {
	return func(context echo.Context) error {
		user := new(models.User)
		if err := context.Bind(user); err != nil {
			return context.NoContent(http.StatusServiceUnavailable)
		}
		user.Nickname = context.Param("nickname")

		return Respond(context, userHandler.userUsecase.Create(user))
	}
}

func (userHandler *UserHandler) HandlerUserGetOne() echo.HandlerFunc {
	return func(context echo.Context) error {
		nickname := context.Param("nickname")
		return Respond(context, userHandler.userUsecase.GetByNickname(nickname))
	}
}

//func (userHandler *UserHandler)