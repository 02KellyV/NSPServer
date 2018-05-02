package models

import (
	"time"
)

//User d
type User struct {
	ID        uint      `json:"id"         xorm:"id bigint not null autoincr pk"`
	FirstName string    `json:"first_name" xorm:"first_name varchar(100) not null"    valid:"required,alphaSpaces"`
	LastName  string    `json:"last_name"  xorm:"last_name varchar(100) not null"     valid:"required,alphaSpaces"`
	Email     string    `json:"email"      xorm:"email varchar(20) not null unique" valid:"required,email"`
	Password  string    `json:"-"          xorm:"password varchar(100) not null"       valid:"required,password,encript"`
	CreateAt  time.Time `json:"create_at"`
	UpdateAt  time.Time `json:"update_at"`
}

//TableName d
func (that User) TableName() string {
	return "users"
}
