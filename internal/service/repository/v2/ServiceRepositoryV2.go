package v2

import (
	"database/sql"
	"github.com/NikitaLobaev/BMSTU-DB/internal/models"
)

type ServiceRepositoryV2 struct {
	dbConnection *sql.DB
}

func NewServiceRepositoryV2(dbConnection *sql.DB) *ServiceRepositoryV2 {
	return &ServiceRepositoryV2{
		dbConnection: dbConnection,
	}
}

func (serviceRepositoryV2 *ServiceRepositoryV2) Truncate() error {
	const query = "TRUNCATE TABLE user_ RESTART IDENTITY CASCADE"
	_, err := serviceRepositoryV2.dbConnection.Exec(query)
	return err
}

func (serviceRepositoryV2 *ServiceRepositoryV2) Select() (*models.Status, error) {
	const query = "SELECT (SELECT COUNT(*) FROM user_), (SELECT COUNT(*) FROM forum), (SELECT COUNT(*) FROM thread), (SELECT COUNT(*) FROM post)"
	status := new(models.Status)
	if err := serviceRepositoryV2.dbConnection.QueryRow(query).Scan(&status.Users, &status.Forums, &status.Threads,
		&status.Posts); err == sql.ErrNoRows {
		return nil, err
	}
	return status, nil
}
