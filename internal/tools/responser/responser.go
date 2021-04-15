package responser

import (
	. "../response"
	"github.com/labstack/echo"
)

func Respond(context echo.Context, response *Response) error {
	if response.JSONObject == nil {
		return context.NoContent(response.Code)
	}
	return context.JSON(response.Code, response.JSONObject)
}
