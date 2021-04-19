package usecase

import (
	"../../models"
	ThreadUsecase "../../thread/usecase"
	. "../../tools/response"
	"../repository"
	"database/sql"
	"net/http"
)

type ForumUsecase struct {
	forumRepository *repository.ForumRepository
	threadUsecase   *ThreadUsecase.ThreadUsecase
}

func NewForumUsecase(forumRepository *repository.ForumRepository) *ForumUsecase {
	return &ForumUsecase{
		forumRepository: forumRepository,
	}
}

func (forumUsecase *ForumUsecase) Create(forum *models.Forum) *Response {
	if existingForum, err := forumUsecase.forumRepository.SelectBySlug(forum.Slug); err == nil {
		return NewResponse(http.StatusConflict, existingForum)
	}

	newForum, err := forumUsecase.forumRepository.Insert(forum)
	if err != nil {
		if err == sql.ErrNoRows {
			return NewResponse(http.StatusNotFound, models.Error{
				Message: "Can't find user with nickname " + forum.UserNickname,
			})
		}
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
		return NewResponse(http.StatusServiceUnavailable, nil)
	}
	return NewResponse(http.StatusOK, forum)
}

func (forumUsecase *ForumUsecase) ThreadCreate(thread *models.Thread) *Response {
	return forumUsecase.threadUsecase.Create(thread)
}
