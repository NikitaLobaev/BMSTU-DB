package models

type PostParams struct {
	Limit uint32
	Since uint64
	Sort  string
	Desc  bool

	isLimitSet bool
	isSinceSet bool
	isSortSet  bool
	isDescSet  bool
}

func (postParams *PostParams) SetLimit(limit uint32) {
	postParams.Limit = limit
	postParams.isLimitSet = true
}

func (postParams *PostParams) SetSince(since uint64) {
	postParams.Since = since
	postParams.isSinceSet = true
}

func (postParams *PostParams) SetSort(sort string) {
	postParams.Sort = sort
	postParams.isSortSet = true
}

func (postParams *PostParams) SetDesc(desc bool) {
	postParams.Desc = desc
	postParams.isDescSet = true
}

func (postParams *PostParams) IsLimitSet() bool {
	return postParams.isLimitSet
}

func (postParams *PostParams) IsSinceSet() bool {
	return postParams.isSinceSet
}

func (postParams *PostParams) IsSortSet() bool {
	return postParams.isSortSet
}

func (postParams *PostParams) IsDescSet() bool {
	return postParams.isDescSet
}
