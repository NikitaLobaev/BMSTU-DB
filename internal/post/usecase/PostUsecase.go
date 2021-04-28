package usecase

import (
	"database/sql"
	ForumUsecase "github.com/NikitaLobaev/BMSTU-DB/internal/forum/usecase"
	"github.com/NikitaLobaev/BMSTU-DB/internal/models"
	"github.com/NikitaLobaev/BMSTU-DB/internal/post/repository"
	ThreadUsecase "github.com/NikitaLobaev/BMSTU-DB/internal/thread/usecase"
	. "github.com/NikitaLobaev/BMSTU-DB/internal/tools/response"
	UserUsecase "github.com/NikitaLobaev/BMSTU-DB/internal/user/usecase"
	"github.com/labstack/gommon/log"
	"net/http"
	"strconv"
)

type PostUsecase struct {
	postRepository repository.PostRepository
	forumUsecase   *ForumUsecase.ForumUsecase
	threadUsecase  *ThreadUsecase.ThreadUsecase
	userUsecase    *UserUsecase.UserUsecase
}

func NewPostUsecase(postRepository repository.PostRepository, forumUsecase *ForumUsecase.ForumUsecase,
	threadUsecase *ThreadUsecase.ThreadUsecase, userUsecase *UserUsecase.UserUsecase) *PostUsecase {
	return &PostUsecase{
		postRepository: postRepository,
		forumUsecase:   forumUsecase,
		threadUsecase:  threadUsecase,
		userUsecase:    userUsecase,
	}
}

func (postUsecase *PostUsecase) GetPostFullById(id uint64, user bool, forum bool, thread bool) *Response {
	post, err := postUsecase.postRepository.SelectById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return NewResponse(http.StatusNotFound, models.Error{
				Message: "Can't find post with id " + strconv.FormatUint(id, 10),
			})
		}
		log.Error(err)
		return NewResponse(http.StatusServiceUnavailable, nil)
	}

	postFull := new(models.PostFull)

	postFull.Post = *post

	if user {
		responseUser := postUsecase.userUsecase.GetByNickname(postFull.Post.UserNickname)
		if responseUser.Code != http.StatusOK {
			log.Error(err)
			return NewResponse(http.StatusServiceUnavailable, nil)
		}
		postFull.User = responseUser.JSONObject.(*models.User)
	}

	if forum {
		responseForum := postUsecase.forumUsecase.GetBySlug(postFull.Post.ForumSlug)
		if responseForum.Code != http.StatusOK {
			log.Error(err)
			return NewResponse(http.StatusServiceUnavailable, nil)
		}
		postFull.Forum = responseForum.JSONObject.(*models.Forum)
	}

	if thread {
		responseThread := postUsecase.threadUsecase.GetBySlugOrId(strconv.FormatUint(uint64(postFull.Post.ThreadId), 10))
		if responseThread.Code != http.StatusOK {
			log.Error(err)
			return NewResponse(http.StatusServiceUnavailable, nil)
		}
		postFull.Thread = responseThread.JSONObject.(*models.Thread)
	}

	return NewResponse(http.StatusOK, postFull)
}

func (postUsecase *PostUsecase) Update(id uint64, postUpdate *models.PostUpdate) *Response {
	post, err := postUsecase.postRepository.Update(id, postUpdate)
	if err != nil {
		if err == sql.ErrNoRows {
			return NewResponse(http.StatusNotFound, models.Error{
				Message: "Can't find post with id " + strconv.FormatUint(id, 10),
			})
		}
		log.Error(err)
		return NewResponse(http.StatusServiceUnavailable, nil)
	}
	return NewResponse(http.StatusOK, post)
}
