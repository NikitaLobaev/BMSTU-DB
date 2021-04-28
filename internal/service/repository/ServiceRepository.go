package repository

import (
	"github.com/NikitaLobaev/BMSTU-DB/internal/models"
)

type ServiceRepository interface {
	Truncate() error
	Select() (*models.Status, error)
}
