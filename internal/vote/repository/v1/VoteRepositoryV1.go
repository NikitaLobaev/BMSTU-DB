package v1

import (
	"database/sql"
	"github.com/NikitaLobaev/BMSTU-DB/internal/models"
	"strconv"
)

type VoteRepositoryV1 struct {
	dbConnection *sql.DB
}

func NewVoteRepositoryV1(dbConnection *sql.DB) *VoteRepositoryV1 {
	return &VoteRepositoryV1{
		dbConnection: dbConnection,
	}
}

//TODO: к каждому методу можно сделать "аннотации", какие поля в переданных аргументах "trusted" (доверенные, уже гарантированно извлечённые из бд, а не пришли от клиента) - чтобы лишний раз не ходить в субд
func (voteRepositoryV1 *VoteRepositoryV1) insertOrUpdateByThreadSlugOrId(isSlug bool, slugOrId string, vote *models.Vote) (*models.Vote, error) {
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

	if err := voteRepositoryV1.dbConnection.QueryRow(query, vote.UserNickname, slugOrId, vote.Voice).Scan(&vote.Voice); err != nil {
		return nil, err
	}
	return vote, nil
}

func (voteRepositoryV1 *VoteRepositoryV1) InsertOrUpdateByThreadSlug(slug string, vote *models.Vote) (*models.Vote, error) {
	return voteRepositoryV1.insertOrUpdateByThreadSlugOrId(true, slug, vote)
}

func (voteRepositoryV1 *VoteRepositoryV1) InsertOrUpdateByThreadId(id uint32, vote *models.Vote) (*models.Vote, error) {
	return voteRepositoryV1.insertOrUpdateByThreadSlugOrId(false, strconv.FormatUint(uint64(id), 10), vote)
}
