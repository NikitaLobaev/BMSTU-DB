package responser

import (
	. "../response"
	"github.com/labstack/echo/v4"
)

func Respond(context echo.Context, response *Response) error {
	if response.JSONObject == nil {
		return context.NoContent(response.Code)
	}
	return context.JSON(response.Code, response.JSONObject)
}
