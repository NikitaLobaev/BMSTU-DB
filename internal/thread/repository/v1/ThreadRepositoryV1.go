package v1

import (
	"database/sql"
	"github.com/NikitaLobaev/BMSTU-DB/internal/models"
	"github.com/labstack/gommon/log"
	"strconv"
)

type ThreadRepositoryV1 struct {
	dbConnection *sql.DB
}

func NewThreadRepositoryV1(dbConnection *sql.DB) *ThreadRepositoryV1 {
	return &ThreadRepositoryV1{
		dbConnection: dbConnection,
	}
}

func (threadRepositoryV1 *ThreadRepositoryV1) Insert(thread *models.Thread) (*models.Thread, error) {
	const query = "INSERT INTO thread (title, user_nickname, forum_slug, message, slug, created) SELECT $1, user_.nickname, forum.slug, $4, $5, $6 FROM user_, forum WHERE user_.nickname = $2 AND forum.slug = $3 RETURNING thread.id, thread.title, thread.user_nickname, thread.forum_slug, thread.message, thread.slug, thread.created"
	var slug sql.NullString
	if err := threadRepositoryV1.dbConnection.QueryRow(query, thread.Title, thread.UserNickname, thread.ForumSlug,
		thread.Message, thread.Slug, thread.Created).Scan(&thread.Id, &thread.Title, &thread.UserNickname,
		&thread.ForumSlug, &thread.Message, &slug, &thread.Created); err != nil {
		return nil, err
	}
	if slug.Valid {
		thread.Slug = slug.String
	}
	return thread, nil
}

func (threadRepositoryV1 *ThreadRepositoryV1) selectBySlugOrId(isSlug bool, slugOrId string) (*models.Thread, error) {
	const query1 = "SELECT id, title, user_nickname, forum_slug, message, votes, slug, created FROM thread WHERE slug = $1"
	const query2 = "SELECT id, title, user_nickname, forum_slug, message, votes, slug, created FROM thread WHERE id = $1"
	var query string
	if isSlug {
		query = query1
	} else {
		query = query2
	}

	thread := new(models.Thread)
	var threadSlug sql.NullString
	if err := threadRepositoryV1.dbConnection.QueryRow(query, slugOrId).Scan(&thread.Id, &thread.Title,
		&thread.UserNickname, &thread.ForumSlug, &thread.Message, &thread.Votes, &threadSlug, &thread.Created); err != nil {
		return nil, err
	}

	if threadSlug.Valid {
		thread.Slug = threadSlug.String
	}

	return thread, nil
}

func (threadRepositoryV1 *ThreadRepositoryV1) SelectBySlug(slug string) (*models.Thread, error) {
	return threadRepositoryV1.selectBySlugOrId(true, slug)
}

func (threadRepositoryV1 *ThreadRepositoryV1) SelectById(id uint32) (*models.Thread, error) {
	return threadRepositoryV1.selectBySlugOrId(false, strconv.FormatUint(uint64(id), 10))
}

func (threadRepositoryV1 *ThreadRepositoryV1) updateBySlugOrId(isSlug bool, slugOrId string, threadUpdate *models.ThreadUpdate) (*models.Thread, error) {
	if threadUpdate.Message == "" && threadUpdate.Title == "" {
		return threadRepositoryV1.selectBySlugOrId(isSlug, slugOrId)
	}

	const query1 = "UPDATE thread SET "
	const query2 = "message = $2"
	const query3 = "title = $2"
	const query4 = "message = $2, title = $3"
	const queryWhere = " WHERE "
	const querySlug = "slug"
	const queryId = "id"
	const query5 = " = $1 RETURNING id, title, user_nickname, forum_slug, message, votes, slug, created"

	var row *sql.Row

	queryEnd := queryWhere
	if isSlug {
		queryEnd += querySlug
	} else {
		queryEnd += queryId
	}
	queryEnd += query5

	query := query1
	if threadUpdate.Message != "" && threadUpdate.Title != "" {
		query += query4 + queryEnd
		row = threadRepositoryV1.dbConnection.QueryRow(query, slugOrId, threadUpdate.Message, threadUpdate.Title)
	} else if threadUpdate.Message != "" {
		query += query2 + queryEnd
		row = threadRepositoryV1.dbConnection.QueryRow(query, slugOrId, threadUpdate.Message)
	} else { //if threadUpdate.Title != ""
		query += query3 + queryEnd
		row = threadRepositoryV1.dbConnection.QueryRow(query, slugOrId, threadUpdate.Title)
	}

	thread := new(models.Thread)
	var threadSlug sql.NullString
	if err := row.Scan(&thread.Id, &thread.Title, &thread.UserNickname, &thread.ForumSlug, &thread.Message,
		&thread.Votes, &threadSlug, &thread.Created); err != nil {
		return nil, err
	}

	if threadSlug.Valid {
		thread.Slug = threadSlug.String
	}

	return thread, nil
}

func (threadRepositoryV1 *ThreadRepositoryV1) UpdateBySlug(slug string, threadUpdate *models.ThreadUpdate) (*models.Thread, error) {
	return threadRepositoryV1.updateBySlugOrId(true, slug, threadUpdate)
}

func (threadRepositoryV1 *ThreadRepositoryV1) UpdateById(id uint32, threadUpdate *models.ThreadUpdate) (*models.Thread, error) {
	return threadRepositoryV1.updateBySlugOrId(false, strconv.FormatUint(uint64(id), 10), threadUpdate)
}

func (threadRepositoryV1 *ThreadRepositoryV1) SelectThreadsBySlug(slug string, forumParams *models.ForumParams) (*models.Threads, error) {
	const query1 = "SELECT id, title, user_nickname, forum_slug, message, votes, slug, created FROM thread WHERE forum_slug = $1"
	const queryCreated1 = " AND created <= $2"
	const queryCreated2 = " AND created >= $2"
	const queryOrderBy = " ORDER BY created"
	const queryDesc = " DESC"
	const queryLimit1 = " LIMIT $2"
	const queryLimit2 = " LIMIT $3"

	query := query1
	var rows *sql.Rows
	var err error
	if forumParams.IsSinceSet() {
		if forumParams.IsDescSet() && forumParams.Desc {
			query += queryCreated1 + queryOrderBy + queryDesc
		} else {
			query += queryCreated2 + queryOrderBy
		}
		if forumParams.IsLimitSet() {
			query += queryLimit2
			rows, err = threadRepositoryV1.dbConnection.Query(query, slug, forumParams.Since, forumParams.Limit)
		} else {
			rows, err = threadRepositoryV1.dbConnection.Query(query, slug, forumParams.Since)
		}
	} else {
		query += queryOrderBy
		if forumParams.IsDescSet() && forumParams.Desc {
			query += queryDesc
		}
		if forumParams.IsLimitSet() {
			query += queryLimit1
			rows, err = threadRepositoryV1.dbConnection.Query(query, slug, forumParams.Limit)
		} else {
			rows, err = threadRepositoryV1.dbConnection.Query(query, slug)
		}
	}

	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Error(err)
		}
	}()

	threads := make(models.Threads, 0)
	for rows.Next() {
		thread := new(models.Thread)
		var threadSlug sql.NullString
		if err = rows.Scan(&thread.Id, &thread.Title, &thread.UserNickname, &thread.ForumSlug, &thread.Message,
			&thread.Votes, &thread.Slug, &thread.Created); err != nil {
			return nil, err
		}

		if threadSlug.Valid {
			thread.Slug = threadSlug.String
		}

		threads = append(threads, thread)
	}

	return &threads, nil
}

func (threadRepositoryV1 *ThreadRepositoryV1) InsertPosts(thread *models.Thread, posts *models.Posts) (*models.Posts, error) {
	const query = "INSERT INTO post (post_parent_id, user_nickname, message, forum_slug, thread_id, created) SELECT $1, user_.nickname, $3, $4, $5, $6 FROM user_ WHERE user_.nickname = $2 RETURNING post.id, post.user_nickname"

	tx, err := threadRepositoryV1.dbConnection.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if tx == nil {
			return
		}
		if err := tx.Rollback(); err != nil {
			log.Error(err)
		}
	}()

	for _, post := range *posts {
		if err = tx.QueryRow(query, post.ParentPostId, post.UserNickname, post.Message, thread.ForumSlug, thread.Id,
			post.Created).Scan(&post.Id, &post.UserNickname); err != nil {
			return nil, err
		}
		post.ThreadId = thread.Id
		post.ForumSlug = thread.ForumSlug
	}

	if err = tx.Commit(); err != nil {
		log.Error(err)
		return nil, err
	}
	tx = nil

	return posts, nil
}

//TODO: заменить в SQL-запросах " на ` (для переноса на следующую строку...)
func (threadRepositoryV1 *ThreadRepositoryV1) SelectPosts(thread *models.Thread, postParams *models.PostParams) (*models.Posts, error) {
	const query1 = "SELECT id, post_parent_id, user_nickname, message, is_edited, created FROM post WHERE"
	const query2 = " thread_id = $1"
	const queryMore = ">"
	const queryLess = "<"
	const queryLimit1 = " LIMIT $2"
	const queryLimit2 = " LIMIT $3"

	query := query1
	var rows *sql.Rows
	var err error
	switch postParams.Sort {
	case "tree":
		const queryOrderBy = " ORDER BY path_, created, id"
		const queryOrderByDesc = " ORDER BY path_ DESC, created, id"
		query += query2
		if postParams.IsSinceSet() {
			const queryPath1 = " AND path_ "
			const queryPath2 = " (SELECT path_ FROM post WHERE id = $2)"
			query += queryPath1
			if postParams.IsDescSet() && postParams.Desc {
				query += queryLess + queryPath2 + queryOrderByDesc
			} else {
				query += queryMore + queryPath2 + queryOrderBy
			}
			if postParams.IsLimitSet() {
				query += queryLimit2
				rows, err = threadRepositoryV1.dbConnection.Query(query, thread.Id, postParams.Since, postParams.Limit)
			} else {
				rows, err = threadRepositoryV1.dbConnection.Query(query, thread.Id, postParams.Since)
			}
		} else {
			if postParams.IsDescSet() && postParams.Desc {
				query += queryOrderByDesc
			} else {
				query += queryOrderBy
			}
			if postParams.IsLimitSet() {
				query += queryLimit1
				rows, err = threadRepositoryV1.dbConnection.Query(query, thread.Id, postParams.Limit)
			} else {
				rows, err = threadRepositoryV1.dbConnection.Query(query, thread.Id)
			}
		}
		break
	case "parent_tree":
		const queryPostRootId = " post_root_id IN (SELECT id FROM post WHERE post_parent_id IS NULL AND" + query2
		const queryDesc = " DESC"
		query += queryPostRootId
		if postParams.IsSinceSet() {
			const queryPostRootId1 = " AND post_root_id "
			const queryPostRootId2 = " (SELECT post_root_id FROM post WHERE id = $2) ORDER BY id"
			const queryLimit = " LIMIT $3"
			const queryOrderBy1 = ") ORDER BY post_root_id"
			const queryOrderBy2 = ", path_, created, id"
			query += queryPostRootId1
			if postParams.IsDescSet() && postParams.Desc {
				query += queryLess + queryPostRootId2 + queryDesc
				if postParams.IsLimitSet() {
					query += queryLimit
				}
				query += queryOrderBy1 + queryDesc + queryOrderBy2
			} else {
				query += queryMore + queryPostRootId2
				if postParams.IsLimitSet() {
					query += queryLimit
				}
				query += queryOrderBy1 + queryOrderBy2
			}
			if postParams.IsLimitSet() {
				rows, err = threadRepositoryV1.dbConnection.Query(query, thread.Id, postParams.Since, postParams.Limit)
			} else {
				rows, err = threadRepositoryV1.dbConnection.Query(query, thread.Id, postParams.Since)
			}
		} else {
			const queryOrderBy = " ORDER BY id"
			const queryLimit = " LIMIT $2"
			const queryOrderBy1 = ") ORDER BY post_root_id"
			const queryOrderBy2 = ", path_, created, id"
			query += queryOrderBy
			if postParams.IsDescSet() && postParams.Desc {
				query += queryDesc + queryLimit + queryOrderBy1 + queryDesc + queryOrderBy2
			} else {
				query += queryLimit + queryOrderBy1 + queryOrderBy2
			}
			if postParams.IsLimitSet() {
				rows, err = threadRepositoryV1.dbConnection.Query(query, thread.Id, postParams.Limit)
			} else {
				rows, err = threadRepositoryV1.dbConnection.Query(query, thread.Id)
			}
		}
		break
	default: //flat
		const queryOrderBy = " ORDER BY created, id"
		const queryOrderByDesc = " ORDER BY created DESC, id DESC"
		query += query2
		if postParams.IsSinceSet() {
			const queryId1 = " AND id "
			const queryId2 = " $2"
			query += queryId1
			if postParams.IsDescSet() && postParams.Desc {
				query += queryLess + queryId2 + queryOrderByDesc
			} else {
				query += queryMore + queryId2 + queryOrderBy
			}
			if postParams.IsLimitSet() {
				query += queryLimit2
				rows, err = threadRepositoryV1.dbConnection.Query(query, thread.Id, postParams.Since, postParams.Limit)
			} else {
				rows, err = threadRepositoryV1.dbConnection.Query(query, thread.Id, postParams.Since)
			}
		} else {
			if postParams.IsDescSet() && postParams.Desc {
				query += queryOrderByDesc
			} else {
				query += queryOrderBy
			}
			if postParams.IsLimitSet() {
				query += queryLimit1
				rows, err = threadRepositoryV1.dbConnection.Query(query, thread.Id, postParams.Limit)
			} else {
				rows, err = threadRepositoryV1.dbConnection.Query(query, thread.Id)
			}
		}
		break
	}

	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Error(err)
		}
	}()

	posts := make(models.Posts, 0)
	for rows.Next() {
		post := new(models.Post)
		var parentPostId sql.NullInt64
		if err := rows.Scan(&post.Id, &parentPostId, &post.UserNickname, &post.Message, &post.IsEdited,
			&post.Created); err != nil {
			return nil, err
		}

		if parentPostId.Valid {
			post.ParentPostId = uint64(parentPostId.Int64)
		}
		post.ForumSlug = thread.ForumSlug
		post.ThreadId = thread.Id

		posts = append(posts, post)
	}

	return &posts, nil
}
