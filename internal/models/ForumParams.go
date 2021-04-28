package models

import "time"

type ForumParams struct {
	Limit uint32
	Since time.Time
	Desc  bool

	isLimitSet bool
	isSinceSet bool
	isDescSet  bool
}

func (forumParams *ForumParams) SetLimit(limit uint32) {
	forumParams.Limit = limit
	forumParams.isLimitSet = true
}

func (forumParams *ForumParams) SetSince(since time.Time) {
	forumParams.Since = since
	forumParams.isSinceSet = true
}

func (forumParams *ForumParams) SetDesc(desc bool) {
	forumParams.Desc = desc
	forumParams.isDescSet = true
}

func (forumParams *ForumParams) IsLimitSet() bool {
	return forumParams.isLimitSet
}

func (forumParams *ForumParams) IsSinceSet() bool {
	return forumParams.isSinceSet
}

func (forumParams *ForumParams) IsDescSet() bool {
	return forumParams.isDescSet
}
