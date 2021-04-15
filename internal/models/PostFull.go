package models

type PostFull struct {
	Post   Post    `json:"post"`
	User   *User   `json:"user,omitempty"`
	Thread *Thread `json:"thread,omitempty"`
	Forum  *Forum  `json:"forum,omitempty"`
}
