package usecase

import (
	"database/sql"
	"github.com/NikitaLobaev/BMSTU-DB/internal/models"
	"github.com/NikitaLobaev/BMSTU-DB/internal/thread/repository"
	. "github.com/NikitaLobaev/BMSTU-DB/internal/tools/response"
	VoteUsecase "github.com/NikitaLobaev/BMSTU-DB/internal/vote/usecase"
	"github.com/labstack/gommon/log"
	"github.com/lib/pq"
	"net/http"
	"strconv"
	"time"
)

type ThreadUsecase struct {
	threadRepository repository.ThreadRepository
	voteUsecase      *VoteUsecase.VoteUsecase
}

func NewThreadUsecase(threadRepository repository.ThreadRepository, voteUsecase *VoteUsecase.VoteUsecase) *ThreadUsecase {
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
		log.Error(err)
		return NewResponse(http.StatusServiceUnavailable, nil)
	}
	return NewResponse(http.StatusCreated, thread2)
}

func (threadUsecase *ThreadUsecase) GetBySlugOrId(slugOrId string) *Response {
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
		log.Error(err)
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
		log.Error(err)
		return NewResponse(http.StatusServiceUnavailable, nil)
	}
	return NewResponse(http.StatusOK, thread)
}

func (threadUsecase *ThreadUsecase) Vote(slugOrId string, vote *models.Vote) *Response {
	responseVote := threadUsecase.voteUsecase.Vote(slugOrId, vote)
	if responseVote.Code != http.StatusCreated {
		return responseVote
	}

	result := threadUsecase.GetBySlugOrId(slugOrId)
	return result
}

func (threadUsecase *ThreadUsecase) GetThreadsByForumSlug(forumSlug string, forumParams *models.ForumParams) *Response {
	threads, err := threadUsecase.threadRepository.SelectThreadsBySlug(forumSlug, forumParams)
	if err != nil {
		log.Error(err)
		return NewResponse(http.StatusServiceUnavailable, nil)
	}
	return NewResponse(http.StatusOK, threads)
}

func (threadUsecase *ThreadUsecase) CreatePosts(slugOrId string, posts *models.Posts) *Response {
	responseThread := threadUsecase.GetBySlugOrId(slugOrId)
	if responseThread.Code != http.StatusOK {
		return responseThread
	}

	location, _ := time.LoadLocation("UTC")
	now := time.Now().In(location).Round(time.Microsecond)
	for _, post := range *posts {
		if post.Created.IsZero() {
			post.Created = now
		}
	}

	posts, err := threadUsecase.threadRepository.InsertPosts(responseThread.JSONObject.(*models.Thread), posts)
	if err != nil {
		if err == sql.ErrNoRows {
			return NewResponse(http.StatusNotFound, models.Error{
				Message: "Can't find one of users",
			})
		} else if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "P0001" {
			return NewResponse(http.StatusConflict, models.Error{
				Message: "Can't find one of parent posts or it was created in another thread",
			})
		}
		log.Error(err)
		return NewResponse(http.StatusServiceUnavailable, nil)
	}
	return NewResponse(http.StatusCreated, posts)
}

func (threadUsecase *ThreadUsecase) GetPosts(slugOrId string, postParams *models.PostParams) *Response {
	responseThread := threadUsecase.GetBySlugOrId(slugOrId)
	if responseThread.Code != http.StatusOK {
		return responseThread
	}

	posts, err := threadUsecase.threadRepository.SelectPosts(responseThread.JSONObject.(*models.Thread), postParams)
	if err != nil {
		log.Error(err)
		return NewResponse(http.StatusServiceUnavailable, nil)
	}
	return NewResponse(http.StatusOK, posts)
}
