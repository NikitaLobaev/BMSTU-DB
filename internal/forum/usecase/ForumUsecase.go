package usecase

import (
	"database/sql"
	"github.com/NikitaLobaev/BMSTU-DB/internal/forum/repository"
	"github.com/NikitaLobaev/BMSTU-DB/internal/models"
	ThreadUsecase "github.com/NikitaLobaev/BMSTU-DB/internal/thread/usecase"
	. "github.com/NikitaLobaev/BMSTU-DB/internal/tools/response"
	UserUsecase "github.com/NikitaLobaev/BMSTU-DB/internal/user/usecase"
	"github.com/labstack/gommon/log"
	"net/http"
)

type ForumUsecase struct {
	forumRepository repository.ForumRepository
	threadUsecase   *ThreadUsecase.ThreadUsecase
	userUsecase     *UserUsecase.UserUsecase
}

func NewForumUsecase(forumRepository repository.ForumRepository, threadUsecase *ThreadUsecase.ThreadUsecase,
	userUsecase *UserUsecase.UserUsecase) *ForumUsecase {
	return &ForumUsecase{
		forumRepository: forumRepository,
		threadUsecase:   threadUsecase,
		userUsecase:     userUsecase,
	}
}

func (forumUsecase *ForumUsecase) Create(forum *models.Forum) *Response {
	if existingForum, err := forumUsecase.forumRepository.SelectBySlug(forum.Slug); existingForum != nil && err == nil {
		return NewResponse(http.StatusConflict, existingForum)
	} else if err != nil && err != sql.ErrNoRows {
		log.Error(err)
		return NewResponse(http.StatusServiceUnavailable, nil)
	}

	newForum, err := forumUsecase.forumRepository.Insert(forum)
	if err != nil {
		if err == sql.ErrNoRows {
			return NewResponse(http.StatusNotFound, models.Error{
				Message: "Can't find user with nickname " + forum.UserNickname,
			})
		}
		log.Error(err)
		return NewResponse(http.StatusServiceUnavailable, nil)
	}
	return NewResponse(http.StatusCreated, newForum)
}

func (forumUsecase *ForumUsecase) GetBySlug(slug string) *Response {
	forum, err := forumUsecase.forumRepository.SelectBySlug(slug)
	if err == sql.ErrNoRows {
		return NewResponse(http.StatusNotFound, models.Error{
			Message: "Can't find forum with slug " + slug,
		})
	} else if err != nil {
		log.Error(err)
		return NewResponse(http.StatusServiceUnavailable, nil)
	}
	return NewResponse(http.StatusOK, forum)
}

func (forumUsecase *ForumUsecase) ThreadCreate(thread *models.Thread) *Response {
	return forumUsecase.threadUsecase.Create(thread)
}

func (forumUsecase *ForumUsecase) GetThreadsBySlug(slug string, forumParams *models.ForumParams) *Response {
	responseForum := forumUsecase.GetBySlug(slug)
	if responseForum.Code != http.StatusOK {
		return responseForum
	}

	return forumUsecase.threadUsecase.GetThreadsByForumSlug(slug, forumParams)
}

func (forumUsecase *ForumUsecase) GetUsersBySlug(slug string, userParams *models.UserParams) *Response {
	responseForum := forumUsecase.GetBySlug(slug)
	if responseForum.Code != http.StatusOK {
		return responseForum
	}

	return forumUsecase.userUsecase.GetUsersByForumSlug(slug, userParams)
}
