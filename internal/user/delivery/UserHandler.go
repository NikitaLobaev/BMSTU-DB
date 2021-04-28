package delivery

import (
	"github.com/NikitaLobaev/BMSTU-DB/internal/models"
	. "github.com/NikitaLobaev/BMSTU-DB/internal/tools/responser"
	"github.com/NikitaLobaev/BMSTU-DB/internal/user/usecase"
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
	echoWS.GET("/api/user/:nickname/profile", userHandler.HandlerUserGetOne())
	echoWS.POST("/api/user/:nickname/profile", userHandler.HandlerUpdate())
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

func (userHandler *UserHandler) HandlerUpdate() echo.HandlerFunc {
	return func(context echo.Context) error {
		userUpdate := new(models.UserUpdate)
		if err := context.Bind(userUpdate); err != nil {
			return context.NoContent(http.StatusServiceUnavailable)
		}
		nickname := context.Param("nickname")
		return Respond(context, userHandler.userUsecase.Update(nickname, userUpdate))
	}
}
