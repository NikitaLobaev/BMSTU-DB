package usecase

import "../repository"

type ForumUsecase struct {
	forumRepository *repository.ForumRepository
}

func NewForumUsecase(forumRepository *repository.ForumRepository) *ForumUsecase {
	return &ForumUsecase{
		forumRepository: forumRepository,
	}
}
