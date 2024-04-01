package main

import "time"

type Resource struct {
	Id   uint64 `json:"id"`
	Name string `json:"name,omitempty"`
}

type SignUpDTO struct {
	Password string `json:"password"`
}

type SignInDTO struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

type AuthToken struct {
	Token       string `json:"token"`
	Account     string `json:"account"`
	ExpiredTime int64  `json:"expired_time"`
}

type CreateCommentDTO struct {
	Score  uint8 `json:"score" validate:"required,min=1,max=10"`
	Score2 uint8 `json:"score2" validate:"required,min=1,max=10"`
	Score3 uint8 `json:"score3" validate:"required,min=1,max=10"`
}

type UserDTO struct {
	Account   string    `json:"account"`
	CreatedAt time.Time `json:"created_at"`
}
