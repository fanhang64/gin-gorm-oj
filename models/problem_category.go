package models

import "gorm.io/gorm"

type ProblemCategory struct {
	gorm.Model
	ProblemId uint `gorm:"column:problem_id;" json:"problem_id"` // 一个问题多个分类的。

	CategoryId    uint          `gorm:"column:category_id;" json:"category_id"`
	CategoryBasic CategoryBasic `gorm:"foreignKey:id;references:category_id"` // CategoryBasic 表的外键id关联，本表的category_id
}

func (ProblemCategory) TableName() string {
	return "problem_category"
}
