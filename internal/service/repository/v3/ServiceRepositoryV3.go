package v3

import (
	"database/sql"
	"github.com/NikitaLobaev/BMSTU-DB/internal/models"
)

type ServiceRepositoryV3 struct {
	dbConnection *sql.DB
}

func NewServiceRepositoryV3(dbConnection *sql.DB) *ServiceRepositoryV3 {
	return &ServiceRepositoryV3{
		dbConnection: dbConnection,
	}
}

func (serviceRepositoryV3 *ServiceRepositoryV3) Truncate() error {
	const query = "TRUNCATE TABLE user_ RESTART IDENTITY CASCADE"
	_, err := serviceRepositoryV3.dbConnection.Exec(query)
	return err
}

func (serviceRepositoryV3 *ServiceRepositoryV3) Select() (*models.Status, error) {
	const query = "SELECT (SELECT COUNT(*) FROM user_), (SELECT COUNT(*) FROM forum), (SELECT COUNT(*) FROM thread), (SELECT COUNT(*) FROM post)"
	status := new(models.Status)
	if err := serviceRepositoryV3.dbConnection.QueryRow(query).Scan(&status.Users, &status.Forums, &status.Threads,
		&status.Posts); err == sql.ErrNoRows {
		return nil, err
	}
	return status, nil
}
