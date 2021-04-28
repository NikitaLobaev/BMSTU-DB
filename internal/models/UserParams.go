package models

type UserParams struct {
	Limit uint32
	Since string
	Desc  bool

	isLimitSet bool
	isSinceSet bool
	isDescSet  bool
}

func (userParams *UserParams) SetLimit(limit uint32) {
	userParams.Limit = limit
	userParams.isLimitSet = true
}

func (userParams *UserParams) SetSince(since string) {
	userParams.Since = since
	userParams.isSinceSet = true
}

func (userParams *UserParams) SetDesc(desc bool) {
	userParams.Desc = desc
	userParams.isDescSet = true
}

func (userParams *UserParams) IsLimitSet() bool {
	return userParams.isLimitSet
}

func (userParams *UserParams) IsSinceSet() bool {
	return userParams.isSinceSet
}

func (userParams *UserParams) IsDescSet() bool {
	return userParams.isDescSet
}
