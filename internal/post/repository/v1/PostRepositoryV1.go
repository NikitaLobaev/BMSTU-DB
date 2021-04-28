package v1

import (
	"database/sql"
	"github.com/NikitaLobaev/BMSTU-DB/internal/models"
)

type PostRepositoryV1 struct {
	dbConnection *sql.DB
}

func NewPostRepositoryV1(dbConnection *sql.DB) *PostRepositoryV1 {
	return &PostRepositoryV1{
		dbConnection: dbConnection,
	}
}

func (postRepositoryV1 *PostRepositoryV1) SelectById(id uint64) (*models.Post, error) {
	const query = "SELECT id, post_parent_id, user_nickname, message, is_edited, forum_slug, thread_id, created FROM post WHERE id = $1"
	post := new(models.Post)
	var postParentId sql.NullInt64
	if err := postRepositoryV1.dbConnection.QueryRow(query, id).Scan(&post.Id, &postParentId, &post.UserNickname,
		&post.Message, &post.IsEdited, &post.ForumSlug, &post.ThreadId, &post.Created); err != nil {
		return nil, err
	}

	if postParentId.Valid {
		post.ParentPostId = uint64(postParentId.Int64)
	}

	return post, nil
}

func (postRepositoryV1 *PostRepositoryV1) Update(id uint64, postUpdate *models.PostUpdate) (*models.Post, error) {
	if postUpdate.Message == "" {
		return postRepositoryV1.SelectById(id)
	}

	const query = "UPDATE post SET message = $1 WHERE id = $2 RETURNING id, post_parent_id, user_nickname, message, is_edited, forum_slug, thread_id, created"
	post := new(models.Post)
	var postParentId sql.NullInt64
	if err := postRepositoryV1.dbConnection.QueryRow(query, postUpdate.Message, id).Scan(&post.Id, &postParentId,
		&post.UserNickname, &post.Message, &post.IsEdited, &post.ForumSlug, &post.ThreadId, &post.Created); err != nil {
		return nil, err
	}

	if postParentId.Valid {
		post.ParentPostId = uint64(postParentId.Int64)
	}

	return post, nil
}
