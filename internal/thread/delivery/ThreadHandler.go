package delivery

import (
	"github.com/NikitaLobaev/BMSTU-DB/internal/models"
	"github.com/NikitaLobaev/BMSTU-DB/internal/thread/usecase"
	. "github.com/NikitaLobaev/BMSTU-DB/internal/tools/responser"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
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
	echoWS.POST("/api/thread/:slug_or_id/create", threadHandler.HandlerPostsCreate())
	echoWS.GET("/api/thread/:slug_or_id/details", threadHandler.HandlerThreadGetOne())
	echoWS.POST("/api/thread/:slug_or_id/details", threadHandler.HandlerThreadUpdate())
	echoWS.GET("/api/thread/:slug_or_id/posts", threadHandler.HandlerThreadGetPosts())
	echoWS.POST("/api/thread/:slug_or_id/vote", threadHandler.HandlerThreadVote())
}

func (threadHandler *ThreadHandler) HandlerThreadGetOne() echo.HandlerFunc {
	return func(context echo.Context) error {
		slugOrId := context.Param("slug_or_id")
		return Respond(context, threadHandler.threadUsecase.GetBySlugOrId(slugOrId))
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

func (threadHandler *ThreadHandler) HandlerPostsCreate() echo.HandlerFunc {
	return func(context echo.Context) error {
		slugOrId := context.Param("slug_or_id")
		posts := new(models.Posts)
		if err := context.Bind(posts); err != nil {
			return context.NoContent(http.StatusServiceUnavailable)
		}
		return Respond(context, threadHandler.threadUsecase.CreatePosts(slugOrId, posts))
	}
}

func (threadHandler *ThreadHandler) HandlerThreadGetPosts() echo.HandlerFunc {
	return func(context echo.Context) error {
		slugOrId := context.Param("slug_or_id")

		postParams := new(models.PostParams)
		if limit, err := strconv.ParseUint(context.QueryParam("limit"), 10, 32); err == nil {
			postParams.SetLimit(uint32(limit))
		}
		if since, err := strconv.ParseUint(context.QueryParam("since"), 10, 64); err == nil {
			postParams.SetSince(since)
		}
		if sort := context.QueryParam("sort"); sort != "" {
			postParams.SetSort(sort)
		}
		if desc, err := strconv.ParseBool(context.QueryParam("desc")); err == nil {
			postParams.SetDesc(desc)
		}

		return Respond(context, threadHandler.threadUsecase.GetPosts(slugOrId, postParams))
	}
}
