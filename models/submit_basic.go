package models

import (
	"gorm.io/gorm"
)

type SubmitBasic struct {
	gorm.Model
	Identity        string       `gorm:"column:identity;type:varchar(36);" json:"identity"`                                  // 分类的唯一标识
	ProblemIdentity string       `gorm:"column:problem_identity;type:varchar(36);comment:'问题的唯一标识'" json:"problem_identity"` // 分类的唯一标识
	ProblemBasic    ProblemBasic `gorm:"foreignKey:identity;references:problem_identity"`

	UserIdentity string    `gorm:"column:user_identity;type:varchar(36);comment:'用户的唯一标识'" json:"user_identity"` // 分类的唯一标识
	UserBasic    UserBasic `gorm:"foreignKey:identity;references:user_identity;"`

	Path   string `gorm:"column:path;type:varchar(255);comment:'代码路径'" json:"path"`
	Status int    `gorm:"column:status;type:tinyint(1);comment:'-1-带判断，1-答案正确，2-答案错误，3-运行超时，4-运行超内存，5编译错误'" json:"status"`
}

func (SubmitBasic) TableName() string {
	return "submit_basic"
}

func GetSubmitList(userIdentity, problemIdentity string, status int) *gorm.DB {
	tx := DB.Debug().Model(&SubmitBasic{}).Preload("ProblemBasic", func(db *gorm.DB) *gorm.DB {
		return db.Omit("content")
	}).Preload("UserBasic", func(db *gorm.DB) *gorm.DB {
		return db.Omit("password")
	})
	if userIdentity != "" {
		tx = tx.Where("user_identity=?", userIdentity)
	}
	if problemIdentity != "" {
		tx = tx.Where("problem_identity=?", problemIdentity)
	}

	if status != 0 {
		tx = tx.Where("status = ?", status)
	}
	return tx
}

func CreateSubmit(sb *SubmitBasic) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(&sb).Error
		if err != nil {
			return err
		}
		// 更新user表的pass_num，submit_num字段
		m := map[string]interface{}{
			"submit_num": gorm.Expr("submit_num+?", 1),
		}
		if sb.Status == 1 {
			m["pass_num"] = gorm.Expr("pass_num+?", 1)
		}
		err = tx.Model(&UserBasic{}).Where("identity=?", sb.UserIdentity).Updates(m).Error
		if err != nil {
			return err
		}

		// 更新problem表的pass_num，submit_num字段
		err = tx.Model(&ProblemBasic{}).Where("identity=?", sb.ProblemIdentity).Updates(m).Error
		if err != nil {
			return err
		}
		return nil
	})
}
