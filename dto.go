package main

import "time"

type Resource struct {
	Id   uint64 `json:"id"`
	Name string `json:"name,omitempty"`
}

type SignUpDTO struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

type SignInDTO struct {
	SignUpDTO
}

type AuthToken struct {
	Token       string `json:"token"`
	ExpiredTime int64  `json:"expired_time"`
}

type CreateSubjectDTO struct {
	Name string `json:"name"`
}

type CreateCommentDTO struct {
	Score int8 `json:"score" validate:"required,min=1,max=10"`
}

type UserDTO struct {
	Id        uint64    `json:"id"`
	Account   string    `json:"account"`
	Gender    uint8     `json:"gender"`
	Name      string    `json:"name"`
	Province  Resource  `json:"province"`
	City      Resource  `json:"city"`
	District  Resource  `json:"district"`
	CreatedAt time.Time `json:"created_at"`
}
