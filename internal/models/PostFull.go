package models

type PostFull struct {
	Post   Post    `json:"post"`
	User   *User   `json:"author,omitempty"`
	Thread *Thread `json:"thread,omitempty"`
	Forum  *Forum  `json:"forum,omitempty"`
}
