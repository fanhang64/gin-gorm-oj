package test

import (
	"fmt"
	"gin_gorm_oj/models"
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestGormTest(t *testing.T) {
	dsn := "root:123@tcp(127.0.0.1:3306)/gin_gorm_oj?charset=utf8mb4&parseTime=true&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Printf("db: %v\n", db)

	data := make([]*models.ProblemBasic, 0)

	if err := db.Find(&data).Error; err != nil {
		t.Fatal(err)
	}

	for _, v := range data {
		fmt.Printf("v: %+v\n", v)
	}
}
