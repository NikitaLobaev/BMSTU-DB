package usecase

import (
	"../../models"
	. "../../tools/response"
	"../repository"
	"database/sql"
	"net/http"
	"strconv"
)

type VoteUsecase struct {
	voteRepository *repository.VoteRepository
}

func NewVoteUsecase(voteRepository *repository.VoteRepository) *VoteUsecase {
	return &VoteUsecase{
		voteRepository: voteRepository,
	}
}

func (voteUsecase *VoteUsecase) Vote(threadSlugOrId string, vote *models.Vote) *Response {
	id, err := strconv.ParseUint(threadSlugOrId, 10, 32)
	var err2 error
	if err == nil {
		err2 = voteUsecase.voteRepository.InsertOrUpdateByThreadId(uint32(id), vote)
	} else {
		err2 = voteUsecase.voteRepository.InsertOrUpdateByThreadSlug(threadSlugOrId, vote)
	}

	if err2 != nil {
		if err2 == sql.ErrNoRows {
			return NewResponse(http.StatusNotFound, models.Error{
				Message: "Can't find user with nickname " + vote.UserNickname + " or thread with slug or id " + threadSlugOrId,
			})
		}
		return NewResponse(http.StatusServiceUnavailable, nil)
	}

	return NewResponse(http.StatusCreated, vote)
}
