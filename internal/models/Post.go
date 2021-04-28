package models

import "time"

type Post struct {
	Id           uint64    `json:"id"`
	ParentPostId uint64    `json:"parent,omitempty"`
	UserNickname string    `json:"author"`
	Message      string    `json:"message"`
	IsEdited     bool      `json:"isEdited"`
	ForumSlug    string    `json:"forum"`
	ThreadId     uint32    `json:"thread"`
	Created      time.Time `json:"created"`
}

type Posts []*Post
