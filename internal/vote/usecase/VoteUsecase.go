package usecase

import (
	"database/sql"
	"github.com/NikitaLobaev/BMSTU-DB/internal/models"
	. "github.com/NikitaLobaev/BMSTU-DB/internal/tools/response"
	"github.com/NikitaLobaev/BMSTU-DB/internal/vote/repository"
	"github.com/labstack/gommon/log"
	"net/http"
	"strconv"
)

type VoteUsecase struct {
	voteRepository repository.VoteRepository
}

func NewVoteUsecase(voteRepository repository.VoteRepository) *VoteUsecase {
	return &VoteUsecase{
		voteRepository: voteRepository,
	}
}

func (voteUsecase *VoteUsecase) Vote(threadSlugOrId string, vote *models.Vote) *Response {
	id, err := strconv.ParseUint(threadSlugOrId, 10, 32)
	var err2 error
	if err == nil {
		_, err2 = voteUsecase.voteRepository.InsertOrUpdateByThreadId(uint32(id), vote)
	} else {
		_, err2 = voteUsecase.voteRepository.InsertOrUpdateByThreadSlug(threadSlugOrId, vote)
	}

	if err2 != nil {
		if err2 == sql.ErrNoRows {
			return NewResponse(http.StatusNotFound, models.Error{
				Message: "Can't find user with nickname " + vote.UserNickname + " or thread with slug or id " + threadSlugOrId,
			})
		}
		log.Error(err2)
		return NewResponse(http.StatusServiceUnavailable, nil)
	}

	return NewResponse(http.StatusCreated, vote)
}
