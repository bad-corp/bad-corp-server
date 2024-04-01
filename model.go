package main

import (
	"time"
)

type User struct {
	Account   string    `json:"account" gorm:"column:account;type:string;not null"`
	Password  string    `json:"password" gorm:"column:password;type:string;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;type:datetime;not null"`
}

func (User) TableName() string {
	return "user"
}

type Corp struct {
	Id        uint64    `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Name      string    `json:"name" gorm:"column:name;type:string;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;type:datetime;not null"`
}

func (Corp) TableName() string {
	return "corp"
}

type CorpComment struct {
	Id          uint64 `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	UserId      string `json:"user_id" gorm:"column:user_id;type:string;not null"`
	CorpId      uint64 `json:"corp_id" gorm:"column:corp_id;type:bigint;not null"`
	Score       uint8  `json:"score" gorm:"column:score;type:tinyint;not null"`
	ScoreCount  uint32 `json:"c" gorm:"column:c;type:int;not null"`
	Score2      uint8  `json:"score2" gorm:"column:score2;type:tinyint;not null"`
	ScoreCount2 uint32 `json:"c2" gorm:"column:c2;type:int;not null"`
	Score3      uint8  `json:"score3" gorm:"column:score3;type:tinyint;not null"`
	ScoreCount3 uint32 `json:"c3" gorm:"column:c3;type:int;not null"`
}

func (CorpComment) TableName() string {
	return "corp_comment"
}

type CorpScore struct {
	CorpId uint64 `json:"corp_id" gorm:"column:corp_id;primaryKey;type:bigint"`
	Score  uint64 `json:"score" gorm:"column:score;type:bigint;not null"`
}

func (CorpScore) TableName() string {
	return "corp_score"
}
