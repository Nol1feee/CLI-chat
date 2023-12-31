package model

import "time"

type UserUpdate struct {
	UserInfo *UserInfo
	Id       int64
}

type UserInfo struct {
	Name  string
	Email string
	Role  string
}

type User struct {
	UserInfo  *UserInfo
	CreatedAt time.Time
	UpdatedAt time.Time
}
