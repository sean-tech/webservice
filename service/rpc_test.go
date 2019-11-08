package service

import (
	"fmt"
	"testing"
)

type User struct {
	UserId int						`json:"user_id" validate:"required,min=1"`
	UserName string					`json:"user_name" validate:"required,len=11"`
	Email string					`json:"email" validate:"required,email"`
}

type GoodsPayInfoParameter struct {
	GoodsId int			`json:"goods_id" validate:"required,min=1"`
	GoodsName string	`json:"goods_name" validate:"required,gte=1"`
	GoodsAmount int		`json:"goods_amount" validate:"required,min=1"`
	Remark string 		`json:"remark" validate:"gte=0"`
}

type GoodsPayParameter struct {
	UserInfo *User					`json:"user_info" validate:"required"`
	Goods []*GoodsPayInfoParameter	`json:"goods" validate:"required,gte=1,dive,required"`
	GoodsIds []int				`json:"goods_ids" validate:"required,gte=1,dive,min=1"`
}

// 正确数据
var parameter1 *GoodsPayParameter = &GoodsPayParameter{
	UserInfo: &User{
		UserId:   101,
		UserName: "12345678901",
		Email:    "123456@qq.com",
	},
	Goods:    []*GoodsPayInfoParameter{&GoodsPayInfoParameter{
		GoodsId:     1001,
		GoodsName:   "三只松鼠干果巧克力100g包邮",
		GoodsAmount: 1,
		Remark:      "",
	}},
	GoodsIds: []int{1},
}
// 错误数据
var parameter2 *GoodsPayParameter = &GoodsPayParameter{
	UserInfo: &User{
		UserId:   0, // 错误1
		UserName: "12345678901000", // 错误2
		Email:    "123456@qq.com",
	},
	Goods:    []*GoodsPayInfoParameter{&GoodsPayInfoParameter{
		GoodsId:     1001,
		GoodsName:   "三只松鼠干果巧克力100g包邮",
		GoodsAmount: 1,
		Remark:      "",
	}},
	GoodsIds: []int{1},
}
// 错误数据
var parameter3 *GoodsPayParameter = &GoodsPayParameter{
	UserInfo: &User{
		UserId:   101,
		UserName: "12345678901",
		Email:    "123456@qq.com",
	},
	Goods:    []*GoodsPayInfoParameter{}, // 错误
	GoodsIds: []int{1},
}
// 错误数据
var parameter4 *GoodsPayParameter = &GoodsPayParameter{
	UserInfo: &User{
		UserId:   101,
		UserName: "12345678901",
		Email:    "123456@qq.com",
	},
	Goods:    []*GoodsPayInfoParameter{&GoodsPayInfoParameter{
		GoodsId:     1001,
		GoodsName:   "三只松鼠干果巧克力100g包邮",
		GoodsAmount: 1,
		Remark:      "",
	}},
	GoodsIds: []int{0}, // 错误
}
// 错误数据
var parameter5 *GoodsPayParameter = &GoodsPayParameter{
	UserInfo: &User{
		UserId:   101,
		UserName: "12345678901",
		Email:    "123456@qq.com",
	},
	Goods:    []*GoodsPayInfoParameter{&GoodsPayInfoParameter{
		GoodsId:     1001,
		GoodsName:   "三只松鼠干果巧克力100g包邮",
		GoodsAmount: 1,
		Remark:      "",
	}},
	GoodsIds: nil, // 错误
}

func TestValidParameter(t *testing.T) {

	var dataArray = []*GoodsPayParameter{parameter1, parameter2, parameter3, parameter4, parameter5}
	for idx, parameter := range dataArray {
		err := ValidParameter(parameter)
		if err != nil {
			fmt.Printf("parameter%d validate failed: %v\n", idx, err)
		} else  {
			fmt.Printf("parameter%d validate success\n", idx)
		}
	}

}