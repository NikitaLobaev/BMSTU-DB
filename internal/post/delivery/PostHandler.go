package delivery

import (
	"github.com/NikitaLobaev/BMSTU-DB/internal/models"
	"github.com/NikitaLobaev/BMSTU-DB/internal/post/usecase"
	. "github.com/NikitaLobaev/BMSTU-DB/internal/tools/responser"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"strings"
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
	echoWS.GET("/api/post/:id/details", postHandler.HandlerPostGetOne())
	echoWS.POST("/api/post/:id/details", postHandler.HandlerPostUpdate())
}

func (postHandler *PostHandler) HandlerPostGetOne() echo.HandlerFunc {
	return func(context echo.Context) error {
		userId, err := strconv.ParseUint(context.Param("id"), 10, 64)
		if err != nil {
			return context.NoContent(http.StatusServiceUnavailable)
		}

		var user, forum, thread bool
		for _, related := range strings.Split(context.QueryParam("related"), ",") {
			switch related {
			case "user":
				user = true
				break
			case "forum":
				forum = true
			case "thread":
				thread = true
				break
			}
		}

		return Respond(context, postHandler.postUsecase.GetPostFullById(userId, user, forum, thread))
	}
}

func (postHandler *PostHandler) HandlerPostUpdate() echo.HandlerFunc {
	return func(context echo.Context) error {
		userId, err := strconv.ParseUint(context.Param("id"), 10, 64)
		if err != nil {
			return context.NoContent(http.StatusServiceUnavailable)
		}

		postUpdate := new(models.PostUpdate)
		if err := context.Bind(postUpdate); err != nil {
			return context.NoContent(http.StatusServiceUnavailable)
		}

		return Respond(context, postHandler.postUsecase.Update(userId, postUpdate))
	}
}
