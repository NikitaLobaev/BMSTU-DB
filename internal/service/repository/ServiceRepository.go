package repository

import (
	"../../models"
	"database/sql"
)

type ServiceRepository struct {
	dbConnection *sql.DB
}

func NewServiceRepository(dbConnection *sql.DB) *ServiceRepository {
	return &ServiceRepository{
		dbConnection: dbConnection,
	}
}

func (serviceRepository *ServiceRepository) Truncate() error {
	const query = "TRUNCATE TABLE profile RESTART IDENTITY CASCADE"
	_, err := serviceRepository.dbConnection.Exec(query)
	return err
}

func (serviceRepository *ServiceRepository) Select() (*models.Status, error) {
	const query = "SELECT (SELECT COUNT(*) FROM profile), (SELECT COUNT(*) FROM forum), (SELECT COUNT(*) FROM thread), (SELECT COUNT(*) FROM post)"
	status := new(models.Status)
	if err := serviceRepository.dbConnection.QueryRow(query).Scan(&status.Users, &status.Forums, &status.Threads,
		&status.Posts); err == sql.ErrNoRows {
		return nil, err
	}
	return status, nil
}
