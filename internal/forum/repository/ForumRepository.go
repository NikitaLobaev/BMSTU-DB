package repository

import (
	"../../models"
	"database/sql"
)

type ForumRepository struct {
	dbConnection *sql.DB
}

func NewForumRepository(dbConnection *sql.DB) *ForumRepository {
	return &ForumRepository{
		dbConnection: dbConnection,
	}
}

func (forumRepository *ForumRepository) SelectBySlug(slug string) (*models.Forum, error) {
	const query = "SELECT title, profile_nickname, slug, posts, threads FROM forum WHERE slug = $1"
	forum := new(models.Forum)
	if err := forumRepository.dbConnection.QueryRow(query, slug).Scan(&forum.Title, &forum.UserNickname, &forum.Slug,
		&forum.Posts, &forum.Threads); err == sql.ErrNoRows {
		return nil, err
	}
	return forum, nil
}

func (forumRepository *ForumRepository) Insert(forum *models.Forum) (*models.Forum, error) {
	const query = "INSERT INTO forum (title, profile_nickname, slug) SELECT $1, $2, profile.nickname FROM profile WHERE profile.nickname = $3 RETURNING forum.profile_nickname"
	if err := forumRepository.dbConnection.QueryRow(query, forum.Title, forum.UserNickname, forum.Slug).
		Scan(&forum.UserNickname); err == sql.ErrNoRows {
		return nil, err
	}
	return forum, nil
}
