package models

type Forum struct {
	Title        string `json:"title"`
	UserNickname string `json:"user"`
	Slug         string `json:"slug"`
	Posts        uint64 `json:"posts,omitempty"`
	Threads      uint32 `json:"threads,omitempty"`
}
