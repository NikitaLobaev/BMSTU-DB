package models

import "time"

type Thread struct {
	Id           uint32    `json:"id"`
	Title        string    `json:"title"`
	UserNickname string    `json:"author"`
	ForumSlug    string    `json:"forum"`
	Message      string    `json:"message"`
	Votes        int32     `json:"votes"`
	Slug         string    `json:"slug"`
	Created      time.Time `json:"created"`
}

//easyjson:json
type Threads []*Thread
