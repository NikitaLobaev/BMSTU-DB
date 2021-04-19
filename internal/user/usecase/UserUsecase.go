package usecase

import (
	"../../models"
	. "../../tools/response"
	"../repository"
	"database/sql"
	"net/http"
)

type UserUsecase struct {
	userRepository *repository.UserRepository
}

func NewUserUsecase(userRepository *repository.UserRepository) *UserUsecase {
	return &UserUsecase{
		userRepository: userRepository,
	}
}

func (userUsecase *UserUsecase) Create(user *models.User) *Response {
	users, err := userUsecase.userRepository.SelectByNicknameOrEmail(user.Nickname, user.Email)
	if err != nil {
		return NewResponse(http.StatusServiceUnavailable, nil)
	}

	if len(users) > 0 {
		return NewResponse(http.StatusConflict, users)
	}

	newUser, err := userUsecase.userRepository.Insert(user)
	if err != nil {
		return NewResponse(http.StatusServiceUnavailable, nil)
	}

	return NewResponse(http.StatusOK, newUser)
}

func (userUsecase *UserUsecase) GetByNickname(nickname string) *Response {
	user, err := userUsecase.userRepository.SelectByNickname(nickname)
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

func (userUsecase *UserUsecase) Update(nickname string, userUpdate *models.UserUpdate) *Response {
	user, err := userUsecase.userRepository.Update(nickname, userUpdate)
	if err != nil {
		return NewResponse(http.StatusServiceUnavailable, nil)
	}
	return NewResponse(http.StatusOK, user)
}
