package data

import "time"

type BaseParameter int32
type BaseResp int32

type BaseModel struct {
	CreateTime time.Time	`db:"create_time" json:"createTime"`
	UpdateTime time.Time	`db:"update_time" json:"updateTime"`
	UpdateUser string		`db:"update_user" json:"updateUser"`
	TbStatus string			`db:"tb_status" json:"-"`
}

