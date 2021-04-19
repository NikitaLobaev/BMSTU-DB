package usecase

import (
	"../../models"
	. "../../tools/response"
	VoteUsecase "../../vote/usecase"
	"../repository"
	"database/sql"
	"net/http"
	"strconv"
)

type ThreadUsecase struct {
	threadRepository *repository.ThreadRepository
	voteUsecase      *VoteUsecase.VoteUsecase
}

func NewThreadUsecase(threadRepository *repository.ThreadRepository, voteUsecase *VoteUsecase.VoteUsecase) *ThreadUsecase {
	return &ThreadUsecase{
		threadRepository: threadRepository,
		voteUsecase:      voteUsecase,
	}
}

func (threadUsecase *ThreadUsecase) Create(thread *models.Thread) *Response {
	if existingThread, err := threadUsecase.threadRepository.SelectBySlug(thread.Slug); err == nil {
		return NewResponse(http.StatusConflict, existingThread)
	}

	thread2, err := threadUsecase.threadRepository.Insert(thread)
	if err != nil {
		if err == sql.ErrNoRows {
			return NewResponse(http.StatusNotFound, models.Error{
				Message: "Can't find user with nickname " + thread.UserNickname + " or forum with slug " + thread.ForumSlug,
			})
		}
		return NewResponse(http.StatusServiceUnavailable, nil)
	}
	return NewResponse(http.StatusOK, thread2)
}

func (threadUsecase *ThreadUsecase) GetDetails(slugOrId string) *Response {
	id, err := strconv.ParseUint(slugOrId, 10, 32)
	var thread *models.Thread
	if err == nil {
		thread, err = threadUsecase.threadRepository.SelectById(uint32(id))
	} else {
		thread, err = threadUsecase.threadRepository.SelectBySlug(slugOrId)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return NewResponse(http.StatusNotFound, models.Error{
				Message: "Can't find thread with slug or id " + slugOrId,
			})
		}
		return NewResponse(http.StatusServiceUnavailable, nil)
	}
	return NewResponse(http.StatusOK, thread)
}

func (threadUsecase *ThreadUsecase) UpdateDetails(slugOrId string, threadUpdate *models.ThreadUpdate) *Response {
	id, err := strconv.ParseUint(slugOrId, 10, 32)
	var thread *models.Thread
	if err == nil {
		thread, err = threadUsecase.threadRepository.UpdateById(uint32(id), threadUpdate)
	} else {
		thread, err = threadUsecase.threadRepository.UpdateBySlug(slugOrId, threadUpdate)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return NewResponse(http.StatusNotFound, models.Error{
				Message: "Can't find thread with slug or id " + slugOrId,
			})
		}
		return NewResponse(http.StatusServiceUnavailable, nil)
	}
	return NewResponse(http.StatusOK, thread)
}

func (threadUsecase *ThreadUsecase) Vote(slugOrId string, vote *models.Vote) *Response {
	responseVote := threadUsecase.voteUsecase.Vote(slugOrId, vote)
	if responseVote.Code != http.StatusCreated {
		return responseVote
	}

	return threadUsecase.GetDetails(slugOrId)
}
