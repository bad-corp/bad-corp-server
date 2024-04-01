package main

import (
	"time"
)

type User struct {
	Id        uint64    `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Account   string    `json:"account" gorm:"column:account;type:string;not null"`
	Password  string    `json:"password" gorm:"column:password;type:string;not null"`
	Gender    uint8     `json:"gender" gorm:"column:gender;type:tinyint;not null"`
	Name      string    `json:"name" gorm:"column:name;type:string;not null"`
	Province  uint32    `json:"province" gorm:"column:province;type:int;not null"`
	City      uint32    `json:"city" gorm:"column:city;type:int;not null"`
	District  uint32    `json:"district" gorm:"column:district;type:int;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;type:date;not null"`
}

func (User) TableName() string {
	return "user"
}

type Subject struct {
	Id        int64     `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Name      string    `json:"name" gorm:"column:name;type:string;not null"`
	CreatedBy uint64    `json:"created_by" gorm:"column:created_by;type:bigint;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;type:date;not null"`
}

func (Subject) TableName() string {
	return "subject"
}

type SubjectComment struct {
	Id               int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	UserId           uint64 `json:"user_id" gorm:"column:user_id;type:bigint;not null"`
	SubjectId        int64  `json:"subject_id" gorm:"column:subject_id;type:bigint;not null"`
	SubjectIdCreator uint64 `json:"subject_id_creator" gorm:"column:subject_id_creator;type:bigint;not null"`
	Score            int8   `json:"score" gorm:"column:score;type:tinyint;not null"`
	// TODO: 设计变更
	C int32 `json:"c" gorm:"column:c;type:int;not null"`
}

func (SubjectComment) TableName() string {
	return "subject_comment"
}

type UserScore struct {
	UserId int64 `json:"user_id" gorm:"column:user_id;primaryKey;type:bigint"`
	Score  int64 `json:"score" gorm:"column:score;type:bigint;not null"`
}

func (UserScore) TableName() string {
	return "user_score"
}

type SubjectBloom struct {
	UserId uint64 `json:"user_id" gorm:"column:user_id;primaryKey;type:bigint"`
	Bloom  []byte `json:"score" gorm:"column:bloom;type:blob;not null"`
}

func (SubjectBloom) TableName() string {
	return "subject_bloom"
}

type Area struct {
	Id uint64 `json:"id" gorm:"column:id;primaryKey"`
}

func (Area) TableName() string {
	return "area"
}
