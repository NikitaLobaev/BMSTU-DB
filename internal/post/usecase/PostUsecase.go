package usecase

import (
	"../../models"
	. "../../tools/response"
	"../repository"
	"database/sql"
	"net/http"
)

type PostUsecase struct {
	postRepository *repository.PostRepository
}

func NewPostUsecase(postRepository *repository.PostRepository) *PostUsecase {
	return &PostUsecase{
		postRepository: postRepository,
	}
}

func (postUsecase *PostUsecase) Create(user *models.User) *Response {
	users, err := postUsecase.postRepository.SelectByNicknameOrEmail(user.Nickname, user.Email)
	if err != nil {
		return NewResponse(http.StatusServiceUnavailable, nil)
	}

	if len(users) > 0 {
		return NewResponse(http.StatusConflict, users)
	}

	if err := postUsecase.postRepository.Insert(user); err != nil {
		return NewResponse(http.StatusServiceUnavailable, nil)
	}
	return NewResponse(http.StatusOK, user)
}

func (postUsecase *PostUsecase) GetByNickname(nickname string) *Response {
	user, err := postUsecase.postRepository.SelectByNickname(nickname)
	if err != nil {
		if err == sql.ErrNoRows {
			return NewResponse(http.StatusNotFound, models.Error{
				Message: "Can't find user with nickname " + nickname,
			})
		}
		return NewResponse(http.StatusServiceUnavailable, nil)
	}
	return NewResponse(http.StatusOK, user)
}

func (postUsecase *PostUsecase) Update(nickname string, userUpdate *models.UserUpdate) *Response {
	user, err := postUsecase.postRepository.Update(nickname, userUpdate)
	if err != nil {
		return NewResponse(http.StatusServiceUnavailable, nil)
	}
	return NewResponse(http.StatusOK, user)
}
