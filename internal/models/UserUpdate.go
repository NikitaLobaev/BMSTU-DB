package models

type UserUpdate struct {
	FullName string `json:"fullname"`
	About    string `json:"about"`
	Email    string `json:"email"`
}
