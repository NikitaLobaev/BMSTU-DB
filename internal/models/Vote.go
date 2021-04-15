package models

type Vote struct {
	UserNickname string `json:"nickname"`
	Voice        int32  `json:"voice"`
}
