package models

import "gorm.io/gorm"

// TestCase 测试案例表
type TestCase struct {
	gorm.Model
	Identity string `gorm:"column:identity;type:varchar(36);not null" json:"identity"`

	ProblemIdentity string `gorm:"column:problem_identity;type:varchar(36)" json:"problem_identity"`
	Input           string `gorm:"column:input;type:varchar()"`
	Output          string `gorm:"column:output;type:varchar()"`
}

func (TestCase) TableName() string {
	return "test_case"
}
