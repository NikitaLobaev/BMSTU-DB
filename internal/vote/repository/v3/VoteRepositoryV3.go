package v3

import (
	"database/sql"
	"github.com/NikitaLobaev/BMSTU-DB/internal/models"
	"strconv"
)

type VoteRepositoryV3 struct {
	dbConnection *sql.DB
}

func NewVoteRepositoryV3(dbConnection *sql.DB) *VoteRepositoryV3 {
	return &VoteRepositoryV3{
		dbConnection: dbConnection,
	}
}

func (voteRepositoryV3 *VoteRepositoryV3) insertOrUpdateByThreadSlugOrId(isSlug bool, slugOrId string, vote *models.Vote) (*models.Vote, error) {
	const query1 = "INSERT INTO vote (user_nickname, thread_id, voice) SELECT user_.nickname, thread.id, $3 FROM user_, thread WHERE user_.nickname = $1"
	const queryThreadSlug = "AND thread.slug = $2"
	const queryThreadId = " AND thread.id = $2"
	const query2 = " ON CONFLICT (user_nickname, thread_id) DO UPDATE SET voice = $3 RETURNING voice"

	query := query1
	if isSlug {
		query += queryThreadSlug
	} else {
		query += queryThreadId
	}
	query += query2

	if err := voteRepositoryV3.dbConnection.QueryRow(query, vote.UserNickname, slugOrId, vote.Voice).Scan(&vote.Voice); err != nil {
		return nil, err
	}
	return vote, nil
}

func (voteRepositoryV3 *VoteRepositoryV3) InsertOrUpdateByThreadSlug(slug string, vote *models.Vote) (*models.Vote, error) {
	return voteRepositoryV3.insertOrUpdateByThreadSlugOrId(true, slug, vote)
}

func (voteRepositoryV3 *VoteRepositoryV3) InsertOrUpdateByThreadId(id uint32, vote *models.Vote) (*models.Vote, error) {
	return voteRepositoryV3.insertOrUpdateByThreadSlugOrId(false, strconv.FormatUint(uint64(id), 10), vote)
}
