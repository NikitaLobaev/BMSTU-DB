package v3

import (
	"database/sql"
	"github.com/NikitaLobaev/BMSTU-DB/internal/models"
)

type PostRepositoryV3 struct {
	dbConnection *sql.DB
}

func NewPostRepositoryV3(dbConnection *sql.DB) *PostRepositoryV3 {
	return &PostRepositoryV3{
		dbConnection: dbConnection,
	}
}

func (postRepositoryV3 *PostRepositoryV3) SelectById(id uint64) (*models.Post, error) {
	const query = "SELECT id, post_parent_id, user_nickname, message, is_edited, forum_slug, thread_id, created FROM post WHERE id = $1"
	post := new(models.Post)
	var postParentId sql.NullInt64
	if err := postRepositoryV3.dbConnection.QueryRow(query, id).Scan(&post.Id, &postParentId, &post.UserNickname,
		&post.Message, &post.IsEdited, &post.ForumSlug, &post.ThreadId, &post.Created); err != nil {
		return nil, err
	}

	if postParentId.Valid {
		post.ParentPostId = uint64(postParentId.Int64)
	}

	return post, nil
}

func (postRepositoryV3 *PostRepositoryV3) Update(id uint64, postUpdate *models.PostUpdate) (*models.Post, error) {
	if postUpdate.Message == "" {
		return postRepositoryV3.SelectById(id)
	}

	const query = "UPDATE post SET message = $1 WHERE id = $2 RETURNING id, post_parent_id, user_nickname, message, is_edited, forum_slug, thread_id, created"
	post := new(models.Post)
	var postParentId sql.NullInt64
	if err := postRepositoryV3.dbConnection.QueryRow(query, postUpdate.Message, id).Scan(&post.Id, &postParentId,
		&post.UserNickname, &post.Message, &post.IsEdited, &post.ForumSlug, &post.ThreadId, &post.Created); err != nil {
		return nil, err
	}

	if postParentId.Valid {
		post.ParentPostId = uint64(postParentId.Int64)
	}

	return post, nil
}
