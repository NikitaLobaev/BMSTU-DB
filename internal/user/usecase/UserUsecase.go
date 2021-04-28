package usecase

import (
	"database/sql"
	"github.com/NikitaLobaev/BMSTU-DB/internal/models"
	. "github.com/NikitaLobaev/BMSTU-DB/internal/tools/response"
	"github.com/NikitaLobaev/BMSTU-DB/internal/user/repository"
	"github.com/labstack/gommon/log"
	"net/http"
)

type UserUsecase struct {
	userRepository repository.UserRepository
}

func NewUserUsecase(userRepository repository.UserRepository) *UserUsecase {
	return &UserUsecase{
		userRepository: userRepository,
	}
}

func (userUsecase *UserUsecase) Create(user *models.User) *Response {
	if existingUsers, err := userUsecase.userRepository.SelectByNicknameOrEmail(user.Nickname, user.Email); existingUsers != nil && err == nil && len(*existingUsers) > 0 {
		return NewResponse(http.StatusConflict, existingUsers)
	} else if err != nil && err != sql.ErrNoRows {
		log.Error(err)
		return NewResponse(http.StatusServiceUnavailable, nil)
	}

	newUser, err := userUsecase.userRepository.Insert(user)
	if err != nil {
		log.Error(err)
		return NewResponse(http.StatusServiceUnavailable, nil)
	}

	return NewResponse(http.StatusCreated, newUser)
}

func (userUsecase *UserUsecase) GetByNickname(nickname string) *Response {
	user, err := userUsecase.userRepository.SelectByNickname(nickname)
	if err != nil {
		if err == sql.ErrNoRows {
			return NewResponse(http.StatusNotFound, models.Error{
				Message: "Can't find user with nickname " + nickname,
			})
		}
		log.Error(err)
		return NewResponse(http.StatusServiceUnavailable, nil)
	}
	return NewResponse(http.StatusOK, user)
}

func (userUsecase *UserUsecase) Update(nickname string, userUpdate *models.UserUpdate) *Response {
	responseUser := userUsecase.GetByNickname(nickname)
	if responseUser.Code != http.StatusOK {
		return responseUser
	}

	user := responseUser.JSONObject.(*models.User)

	if userUpdate.Email == "" {
		userUpdate.Email = user.Email
	}
	if userUpdate.About == "" {
		userUpdate.About = user.About
	}
	if userUpdate.FullName == "" {
		userUpdate.FullName = user.FullName
	}

	if user.Email != userUpdate.Email {
		if existingUser, err := userUsecase.userRepository.SelectByEmail(userUpdate.Email); existingUser != nil &&
			err == nil {
			return NewResponse(http.StatusConflict, models.Error{
				Message: "This email is already registered by user with nickname " + existingUser.Nickname,
			})
		} else if err != nil && err != sql.ErrNoRows {
			log.Error(err)
			return NewResponse(http.StatusServiceUnavailable, nil)
		}
	}

	user, err := userUsecase.userRepository.Update(nickname, userUpdate)
	if err != nil {
		log.Error(err)
		return NewResponse(http.StatusServiceUnavailable, nil)
	}
	return NewResponse(http.StatusOK, user)
}

func (userUsecase *UserUsecase) GetUsersByForumSlug(forumSlug string, userParams *models.UserParams) *Response {
	users, err := userUsecase.userRepository.SelectUsersByForumSlug(forumSlug, userParams)
	if err != nil {
		log.Error(err)
		return NewResponse(http.StatusServiceUnavailable, nil)
	}
	return NewResponse(http.StatusOK, users)
}
