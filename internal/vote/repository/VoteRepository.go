package repository

import (
	"github.com/NikitaLobaev/BMSTU-DB/internal/models"
)

type VoteRepository interface {
	InsertOrUpdateByThreadSlug(string, *models.Vote) (*models.Vote, error)
	InsertOrUpdateByThreadId(uint32, *models.Vote) (*models.Vote, error)
}
