package repository

import "database/sql"

type ForumRepository struct {
	dbConnection *sql.DB
}

func NewForumRepository(dbConnection *sql.DB) *ForumRepository {
	return &ForumRepository{
		dbConnection: dbConnection,
	}
}
