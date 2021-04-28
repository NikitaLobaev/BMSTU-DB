package v1

import (
	"database/sql"
	"github.com/NikitaLobaev/BMSTU-DB/internal/models"
)

type ServiceRepositoryV1 struct {
	dbConnection *sql.DB
}

func NewServiceRepositoryV1(dbConnection *sql.DB) *ServiceRepositoryV1 {
	return &ServiceRepositoryV1{
		dbConnection: dbConnection,
	}
}

func (serviceRepositoryV1 *ServiceRepositoryV1) Truncate() error {
	const query = "TRUNCATE TABLE user_ RESTART IDENTITY CASCADE"
	_, err := serviceRepositoryV1.dbConnection.Exec(query)
	return err
}

func (serviceRepositoryV1 *ServiceRepositoryV1) Select() (*models.Status, error) {
	const query = "SELECT (SELECT COUNT(*) FROM user_), (SELECT COUNT(*) FROM forum), (SELECT COUNT(*) FROM thread), (SELECT COUNT(*) FROM post)"
	status := new(models.Status)
	if err := serviceRepositoryV1.dbConnection.QueryRow(query).Scan(&status.Users, &status.Forums, &status.Threads,
		&status.Posts); err == sql.ErrNoRows {
		return nil, err
	}
	return status, nil
}
