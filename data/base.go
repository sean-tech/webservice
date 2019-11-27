package data

import (
	"time"
)

type BaseParameter int32
type BaseResp int32

type BaseModel struct {
	CreateTime time.Time	`db:"create_time" json:"createTime"`
	UpdateTime time.Time	`db:"update_time" json:"updateTime"`
	UpdateUser string		`db:"update_user" json:"updateUser"`
	TbStatus string			`db:"tb_status" json:"-"`
}

type Error struct {
	Code int
	Msg string
}

func (this *Error) Error() string {
	return this.Msg
}

func NewError(code int, msg string) *Error {
	return &Error{
		Code: code,
		Msg:  msg,
	}
}
