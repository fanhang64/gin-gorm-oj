package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"gin_gorm_oj/helper"
	"gorm.io/gorm"
	"strconv"
)

type ProblemBasic struct {
	gorm.Model
	Identity   string `gorm:"column:identity;type:varchar(36);" json:"identity" form:"identity"`           // 问题的唯一标识
	Title      string `gorm:"column:title;type:varchar(255)" form:"title" json:"title" binding:"required"` // 文章标题
	Content    string `gorm:"column:content;type:text;" form:"content" json:"content" binding:"required"`  // 文章内容
	MaxRuntime int64  `gorm:"column:max_runtime;type:int(11);comment:'最大运行时'" form:"max_runtime" json:"max_runtime" binding:"required"`
	MaxMem     int64  `gorm:"column:max_mem;type:int(11);comment:'最大内存'" form:"max_mem" json:"max_mem" binding:"required"`

	// 一个问题，多个分类
	ProblemCategories []ProblemCategory `gorm:"foreignKey:problem_id;references:id;"` // 关联问题分类表

	// 一个问题，多个测试案例
	TestCases []TestCase `gorm:"foreignKey:ProblemIdentity;references:identity"`

	PassNum   int64 `gorm:"column:pass_num;type:int(11);comment:'完成问题个数'" json:"pass_num"`
	SubmitNum int64 `gorm:"column:submit_num;type:int(11);comment:'提交次数'" json:"submit_num"`
}

func (ProblemBasic) TableName() string {
	return "problem_basic"
}

func GetProblemList(keyword, categoryIdentity string) *gorm.DB {
	tx := DB.Model(&ProblemBasic{}).Preload("ProblemCategories.CategoryBasic").
		Where("title like ? or content like ?", "%"+keyword+"%", "%"+keyword+"%")
	if categoryIdentity != "" {
		tx = tx.Joins("right join problem_category pc on pc.problem_id = problem_basic.id").
			Where("pc.category_id = (select id from category_basic where identity = ?)", categoryIdentity)
	}
	return tx
}

func GetProblem(identity string) (*ProblemBasic, error) {
	var problemBasic ProblemBasic
	err := DB.Where("identity = ?", identity).First(&problemBasic).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("查询问题记录不存在！")
		} else {
			return nil, err
		}
	}
	return &problemBasic, nil
}

func CreteProblem(problem *ProblemBasic, categoryIds []string, testCases []string) error {
	problem.Identity = helper.GetUUID()
	//err := DB.Create(problem).Error
	//if err != nil {
	//	return err
	//}
	problemCategories := make([]ProblemCategory, 0, 10)

	for _, cId := range categoryIds {
		newCId, _ := strconv.ParseUint(cId, 10, 64)
		problemCategories = append(problemCategories, ProblemCategory{
			ProblemId:  problem.ID, // 没有创建对象也可以拿到ID
			CategoryId: uint(newCId),
		})
	}
	problem.ProblemCategories = problemCategories

	testCSli := make([]TestCase, 0, 10)
	for _, t := range testCases {
		caseMap := make(map[string]string)
		fmt.Printf("t: %v\n", t)
		e := json.Unmarshal([]byte(t), &caseMap)
		if e != nil {
			return e
		}
		if _, ok := caseMap["input"]; !ok {
			return errors.New("input不存在")
		}
		if _, ok := caseMap["output"]; !ok {
			return errors.New("output不存在")
		}
		testCSli = append(testCSli, TestCase{
			Identity:        helper.GetUUID(),
			ProblemIdentity: problem.Identity,
			Input:           caseMap["input"],
			Output:          caseMap["output"],
		})
	}
	problem.TestCases = testCSli
	return DB.Create(problem).Error
}

func UpdateProblem(problem *ProblemBasic, categoryIds []string, testCases []string) error {
	err := DB.Transaction(func(tx *gorm.DB) error {
		e := tx.Model(&ProblemBasic{}).Where("identity=?", problem.Identity).Updates(problem).Error
		if e != nil {
			return e
		}
		e = tx.Where("identity=?", problem.Identity).First(problem).Error
		if e != nil {
			return e
		}
		// 问题分类更新
		e = tx.Delete(&ProblemCategory{}, "problem_id=?", problem.ID).Error
		if e != nil {
			return e
		}

		problemCateList := make([]ProblemCategory, 0)
		for _, cId := range categoryIds {
			newCId, _ := strconv.Atoi(cId)
			problemCateList = append(problemCateList, ProblemCategory{
				ProblemId:  problem.ID,
				CategoryId: uint(newCId),
			})
		}
		e = tx.Create(&problemCateList).Error
		if e != nil {
			return e
		}

		//problem.ProblemCategories = problemCateList
		// 问题案例更新
		e = tx.Delete(&TestCase{}, "problem_identity=?", problem.Identity).Error
		if e != nil {
			return e
		}
		testCaseSli := make([]TestCase, 0)
		for _, testCase := range testCases {
			caseMap := map[string]string{}
			e = json.Unmarshal([]byte(testCase), &caseMap)
			if e != nil {
				return e
			}
			if _, ok := caseMap["input"]; !ok {
				return errors.New("input 不存在！")
			}
			if _, ok := caseMap["output"]; !ok {
				return errors.New("output不存在！")
			}
			testCaseSli = append(testCaseSli, TestCase{
				Identity:        helper.GetUUID(),
				ProblemIdentity: problem.Identity,
				Input:           caseMap["input"],
				Output:          caseMap["output"],
			})
		}
		//problem.TestCases = testCaseSli
		e = tx.Create(&testCaseSli).Error
		if e != nil {
			return e
		}
		return nil
	})
	return err
}

func GetProblemPreload(identity string) (*ProblemBasic, error) {
	var problemBasic ProblemBasic
	err := DB.Where("identity = ?", identity).Preload("TestCases").First(&problemBasic).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("查询问题记录不存在！")
		} else {
			return nil, err
		}
	}
	return &problemBasic, nil
}
