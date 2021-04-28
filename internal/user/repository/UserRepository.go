package repository

import (
	"github.com/NikitaLobaev/BMSTU-DB/internal/models"
)

type UserRepository interface {
	Insert(*models.User) (*models.User, error)
	SelectByNicknameOrEmail(string, string) (*models.Users, error)
	SelectByNickname(string) (*models.User, error)
	SelectByEmail(string) (*models.User, error)
	Update(string, *models.UserUpdate) (*models.User, error)
	SelectUsersByForumSlug(string, *models.UserParams) (*models.Users, error)
}
