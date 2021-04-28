package repository

import (
	"github.com/NikitaLobaev/BMSTU-DB/internal/models"
)

type ForumRepository interface {
	SelectBySlug(string) (*models.Forum, error)
	Insert(*models.Forum) (*models.Forum, error)
}
