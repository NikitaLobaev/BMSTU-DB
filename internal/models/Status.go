package models

type Status struct {
	Users  uint32 `json:"user"`
	Forums uint32 `json:"forum"`
	Thread uint32 `json:"thread"`
	Posts  uint64 `json:"post"`
}
