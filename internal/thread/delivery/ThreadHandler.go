package delivery

import (
	"../../models"
	. "../../tools/responser"
	"../usecase"
	"github.com/labstack/echo/v4"
	"net/http"
)

type ThreadHandler struct {
	threadUsecase *usecase.ThreadUsecase
}

func NewThreadHandler(userUsecase *usecase.ThreadUsecase) *ThreadHandler {
	return &ThreadHandler{
		threadUsecase: userUsecase,
	}
}

func (threadHandler *ThreadHandler) Configure(echoWS *echo.Echo) {
	echoWS.POST("/api/thread/:slug_or_id/create", nil)
	echoWS.GET("/api/thread/:slug_or_id/details", threadHandler.HandlerThreadGetOne())
	echoWS.POST("/api/thread/:slug_or_id/details", threadHandler.HandlerThreadUpdate())
	echoWS.GET("/api/thread/:slug_or_id/posts", nil)
	echoWS.POST("/api/thread/:slug_or_id/vote", threadHandler.HandlerThreadVote())
}

func (threadHandler *ThreadHandler) HandlerThreadGetOne() echo.HandlerFunc {
	return func(context echo.Context) error {
		slugOrId := context.Param("slug_or_id")
		return Respond(context, threadHandler.threadUsecase.GetDetails(slugOrId))
	}
}

func (threadHandler *ThreadHandler) HandlerThreadUpdate() echo.HandlerFunc {
	return func(context echo.Context) error {
		slugOrId := context.Param("slug_or_id")
		threadUpdate := new(models.ThreadUpdate)
		if err := context.Bind(threadUpdate); err != nil {
			return context.NoContent(http.StatusServiceUnavailable)
		}
		return Respond(context, threadHandler.threadUsecase.UpdateDetails(slugOrId, threadUpdate))
	}
}

func (threadHandler *ThreadHandler) HandlerThreadVote() echo.HandlerFunc {
	return func(context echo.Context) error {
		slugOrId := context.Param("slug_or_id")
		vote := new(models.Vote)
		if err := context.Bind(vote); err != nil {
			return context.NoContent(http.StatusServiceUnavailable)
		}
		return Respond(context, threadHandler.threadUsecase.Vote(slugOrId, vote))
	}
}
