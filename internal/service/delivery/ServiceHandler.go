package delivery

import (
	"github.com/NikitaLobaev/BMSTU-DB/internal/service/usecase"
	. "github.com/NikitaLobaev/BMSTU-DB/internal/tools/responser"
	"github.com/labstack/echo/v4"
)

type ServiceHandler struct {
	serviceUsecase *usecase.ServiceUsecase
}

func NewServiceHandler(userUsecase *usecase.ServiceUsecase) *ServiceHandler {
	return &ServiceHandler{
		serviceUsecase: userUsecase,
	}
}

func (serviceHandler *ServiceHandler) Configure(echoWS *echo.Echo) {
	echoWS.POST("/api/service/clear", serviceHandler.HandlerServiceClear())
	echoWS.GET("/api/service/status", serviceHandler.HandlerServiceStatus())
}

func (serviceHandler *ServiceHandler) HandlerServiceClear() echo.HandlerFunc {
	return func(context echo.Context) error {
		return Respond(context, serviceHandler.serviceUsecase.Clear())
	}
}

func (serviceHandler *ServiceHandler) HandlerServiceStatus() echo.HandlerFunc {
	return func(context echo.Context) error {
		return Respond(context, serviceHandler.serviceUsecase.Status())
	}
}
