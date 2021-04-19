package delivery

import (
	. "../../tools/responser"
	"../usecase"
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
	echoWS.POST("/service/clear", serviceHandler.HandlerServiceClear())
	echoWS.GET("/service/status", serviceHandler.HandlerServiceStatus())
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
