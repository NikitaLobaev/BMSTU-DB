package repository

import (
	"github.com/NikitaLobaev/BMSTU-DB/internal/models"
)

type ThreadRepository interface {
	Insert(*models.Thread) (*models.Thread, error)
	SelectBySlug(string) (*models.Thread, error)
	SelectById(uint32) (*models.Thread, error)
	UpdateBySlug(string, *models.ThreadUpdate) (*models.Thread, error)
	UpdateById(uint32, *models.ThreadUpdate) (*models.Thread, error)
	SelectThreadsBySlug(string, *models.ForumParams) (*models.Threads, error)
	InsertPosts(*models.Thread, *models.Posts) (*models.Posts, error)
	SelectPosts(*models.Thread, *models.PostParams) (*models.Posts, error)
}
