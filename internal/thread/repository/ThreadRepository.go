package repository

import (
	"../../models"
	"database/sql"
	"strconv"
)

type ThreadRepository struct {
	dbConnection *sql.DB
}

func NewThreadRepository(dbConnection *sql.DB) *ThreadRepository {
	return &ThreadRepository{
		dbConnection: dbConnection,
	}
}

func (threadRepository *ThreadRepository) Insert(thread *models.Thread) (*models.Thread, error) {
	const query = "INSERT INTO thread (title, profile_nickname, forum_slug, message, slug, created) SELECT $1, profile.nickname, forum.slug, $4, $5, $6 FROM profile, forum WHERE profile.nickname = $1 AND forum.slug = $3 RETURNING thread.id, thread.profile_nickname, thread.forum_slug"
	if err := threadRepository.dbConnection.QueryRow(query, thread.Title, thread.UserNickname, thread.ForumSlug,
		thread.Message, thread.Slug, thread.Created).Scan(&thread.Title, &thread.UserNickname, &thread.ForumSlug,
		&thread.Message, &thread.Slug, &thread.Created); err != nil {
		return nil, err
	}
	return thread, nil
}

//TODO: возможно, этот метод использоваться не будет и достаточно SelectById...
func (threadRepository *ThreadRepository) selectBySlugOrId(isSlug bool, slugOrId string) (*models.Thread, error) {
	const query1 = "SELECT id, title, profile_nickname, forum_slug, message, votes, slug, created FROM thread WHERE slug = $1"
	const query2 = "SELECT id, title, profile_nickname, forum_slug, message, votes, slug, created FROM thread WHERE id = $1"
	var query string
	if isSlug {
		query = query1
	} else {
		query = query2
	}

	thread := new(models.Thread)
	var threadSlug sql.NullString
	if err := threadRepository.dbConnection.QueryRow(query, slugOrId).Scan(&thread.Id, &thread.Title, &thread.UserNickname,
		&thread.ForumSlug, &thread.Message, &thread.Votes, &threadSlug, &thread.Created); err != nil {
		return nil, err
	}

	if threadSlug.Valid {
		thread.Slug = threadSlug.String
	}

	return thread, nil
}

func (threadRepository *ThreadRepository) SelectBySlug(slug string) (*models.Thread, error) {
	return threadRepository.selectBySlugOrId(true, slug)
}

func (threadRepository *ThreadRepository) SelectById(id uint32) (*models.Thread, error) {
	return threadRepository.selectBySlugOrId(false, strconv.FormatUint(uint64(id), 10))
}

func (threadRepository *ThreadRepository) updateBySlugOrId(isSlug bool, slugOrId string, threadUpdate *models.ThreadUpdate) (*models.Thread, error) {
	const query1 = "UPDATE thread SET message = $2, title = $3 WHERE slug = $1 RETURNING id, title, profile_nickname, forum_slug, message, votes, slug, created"
	const query2 = "UPDATE thread SET message = $2, title = $3 WHERE id = $1 RETURNING id, title, profile_nickname, forum_slug, message, votes, slug, created"
	var query string
	if isSlug {
		query = query1
	} else {
		query = query2
	}

	thread := new(models.Thread)
	var threadSlug sql.NullString
	if err := threadRepository.dbConnection.QueryRow(query, slugOrId, threadUpdate.Message, threadUpdate.Title).
		Scan(&thread.Id, &thread.Title, &thread.UserNickname, &thread.ForumSlug, &thread.Message, &thread.Votes,
			&threadSlug, &thread.Created); err != nil {
		return nil, err
	}

	if threadSlug.Valid {
		thread.Slug = threadSlug.String
	}

	return thread, nil
}

func (threadRepository *ThreadRepository) UpdateBySlug(slug string, threadUpdate *models.ThreadUpdate) (*models.Thread, error) {
	return threadRepository.updateBySlugOrId(false, slug, threadUpdate)
}

func (threadRepository *ThreadRepository) UpdateById(id uint32, threadUpdate *models.ThreadUpdate) (*models.Thread, error) {
	return threadRepository.updateBySlugOrId(true, strconv.FormatUint(uint64(id), 10), threadUpdate)
}
