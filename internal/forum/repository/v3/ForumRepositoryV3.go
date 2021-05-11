package v3

import (
	"database/sql"
	"github.com/NikitaLobaev/BMSTU-DB/internal/models"
)

type ForumRepositoryV3 struct {
	dbConnection *sql.DB
}

func NewForumRepositoryV3(dbConnection *sql.DB) *ForumRepositoryV3 {
	return &ForumRepositoryV3{
		dbConnection: dbConnection,
	}
}

func (forumRepositoryV3 *ForumRepositoryV3) SelectBySlug(slug string) (*models.Forum, error) {
	const query = "SELECT title, user_nickname, slug, posts, threads FROM forum WHERE slug = $1"
	forum := new(models.Forum)
	if err := forumRepositoryV3.dbConnection.QueryRow(query, slug).Scan(&forum.Title, &forum.UserNickname, &forum.Slug,
		&forum.Posts, &forum.Threads); err == sql.ErrNoRows {
		return nil, err
	}
	return forum, nil
}

func (forumRepositoryV3 *ForumRepositoryV3) Insert(forum *models.Forum) (*models.Forum, error) {
	const query = "INSERT INTO forum (title, user_nickname, slug) SELECT $1, user_.nickname, $3 FROM user_ WHERE user_.nickname = $2 RETURNING forum.user_nickname"
	if err := forumRepositoryV3.dbConnection.QueryRow(query, forum.Title, forum.UserNickname, forum.Slug).
		Scan(&forum.UserNickname); err == sql.ErrNoRows {
		return nil, err
	}
	return forum, nil
}
