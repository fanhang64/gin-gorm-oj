package models

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type CategoryBasic struct {
	gorm.Model
	Identity string `gorm:"column:identity;type=varchar(36);" json:"identity" form:"identity"` // 分类的唯一标识
	Name     string `gorm:"column:name;type:varchar(100);comment:'分类名称'" json:"name" form:"name" binding:"required"`
	ParentId int64  `gorm:"column:parent_id;default:0;comment:'父级id'" json:"parent_id" form:"parent_id"`
}

func (CategoryBasic) TableName() string {
	return "category_basic"
}

func GetCategoryList(offset, size int, keyword string) (categoryBasicList []CategoryBasic, err error) {
	err = DB.Debug().Model(&CategoryBasic{}).Where("name like ?", "%"+keyword+"%").
		Offset(offset).Limit(size).Find(&categoryBasicList).Error

	fmt.Printf("categoryBasicList: %v\n", categoryBasicList)
	return
}

func CreateCategory(categoryBasic *CategoryBasic) error {
	return DB.Create(categoryBasic).Error
}

func UpdateCategory(basic *CategoryBasic) error {
	return DB.Model(basic).Where("identity=?", basic.Identity).Updates(basic).Error
}

func DeleteCategory(identity string) error {
	var count int64
	err := DB.Model(&ProblemCategory{}).
		Where("category_id=(select id from category_basic where identity=?)", identity).Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("获取分类下存在问题，不可删除！")
	}
	return DB.Delete(&CategoryBasic{}, "identity=?", identity).Error
}
