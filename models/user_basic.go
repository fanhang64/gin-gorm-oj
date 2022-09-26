package models

import (
	"errors"
	"fmt"
	"gin_gorm_oj/helper"
	"gorm.io/gorm"
)

type UserBasic struct {
	gorm.Model
	Identity string `gorm:"column:identity;type=varchar(36);" json:"identity"` // 用户的唯一标识
	Name     string `gorm:"column:name;type:varchar(100);comment:'名称'" json:"name" form:"name" binding:"required" msg:"用户名不能为空"`
	Password string `gorm:"column:password;type:varchar(32);comment:'密码'" json:"password" form:"password" binding:"required" msg:"密码不能为空"`
	Phone    string `gorm:"column:phone;type:varchar(20);comment:'手机号'" json:"phone" form:"phone"`
	Email    string `gorm:"column:email;type:varchar(100);comment:'邮箱'" json:"email" form:"email" binding:"required" msg:"邮箱不能为空"`

	PassNum   int64 `gorm:"column:pass_num;type:int(11);comment:'完成问题个数'" json:"pass_num"`
	SubmitNum int64 `gorm:"column:submit_num;type:int(11);comment:'提交次数'" json:"submit_num"`

	IsAdmin int `gorm:"column:is_admin;type:tinyint(1);comment:'是否管理员,0否，1是'" json:"is_admin"`
}

func (UserBasic) TableName() string {
	return "user_basic"
}

func GetUserDetail(identity string) (*UserBasic, error) {
	var userBasic UserBasic
	if err := DB.Omit("password").Where("identity = ?", identity).First(&userBasic).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("查询不到用户详情！")
		} else {
			return nil, err
		}
	}
	return &userBasic, nil
}

func GetUserInfo(username, password string) (*UserBasic, error) {
	var user UserBasic
	err := DB.Where("name=? ", username).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("找不到用户信息！")
		} else {
			return nil, err
		}
	} else {
		if user.Password != password {
			return nil, errors.New("用户名与密码不匹配！")
		}
	}
	fmt.Printf("user: %v\n", user)
	return &user, nil
}

func CreateUser(basic *UserBasic) error {
	basic.Identity = helper.GetUUID()
	basic.Password = helper.GetMd5(basic.Password)
	return DB.Create(basic).Error
}

func CheckUserExist(email string) error {
	var count int64
	err := DB.Model(&UserBasic{}).Where("email=?", email).Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("用户名已存在！")
	}
	return nil
}

func GetRankList() *gorm.DB {
	return DB.Model(&UserBasic{}).Order("pass_num desc, submit_num asc")
}
