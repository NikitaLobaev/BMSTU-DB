package repository

import (
	"../../models"
	"database/sql"
	"strconv"
)

type VoteRepository struct {
	dbConnection *sql.DB
}

func NewVoteRepository(dbConnection *sql.DB) *VoteRepository {
	return &VoteRepository{
		dbConnection: dbConnection,
	}
}

func (voteRepository *VoteRepository) insertOrUpdateByThreadSlugOrId(isSlug bool, slugOrId string, vote *models.Vote) error {
	const query1 = "INSERT INTO vote (profile_id, thread_id, voice) SELECT profile.id, thread.id, $3 FROM profile, thread WHERE profile.nickname = $1 AND thread.slug = $2 ON CONFLICT (profile_id, thread_id) DO UPDATE SET voice = $3"
	const query2 = "INSERT INTO vote (profile_id, thread_id, voice) SELECT profile.id, thread.id, $3 FROM profile, thread WHERE profile.nickname = $1 AND thread.id = $2 ON CONFLICT (profile_id, thread_id) DO UPDATE SET voice = $3"
	var query string
	if isSlug {
		query = query1
	} else {
		query = query2
	}

	_, err := voteRepository.dbConnection.Exec(query, slugOrId, vote.UserNickname, vote.Voice)
	return err
}

func (voteRepository *VoteRepository) InsertOrUpdateByThreadSlug(slug string, vote *models.Vote) error {
	return voteRepository.insertOrUpdateByThreadSlugOrId(true, slug, vote)
}

func (voteRepository *VoteRepository) InsertOrUpdateByThreadId(id uint32, vote *models.Vote) error {
	return voteRepository.insertOrUpdateByThreadSlugOrId(false, strconv.FormatUint(uint64(id), 10), vote)
}
