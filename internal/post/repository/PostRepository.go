package repository

import (
	"github.com/NikitaLobaev/BMSTU-DB/internal/models"
)

type PostRepository interface {
	SelectById(uint64) (*models.Post, error)
	Update(uint64, *models.PostUpdate) (*models.Post, error)
}
