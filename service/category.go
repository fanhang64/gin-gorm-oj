package service

import (
	"fmt"
	"gin_gorm_oj/constants"
	"gin_gorm_oj/helper"
	"gin_gorm_oj/models"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

// GetCategoryListHandler
// @Tags 管理员方法
// @Summary 获取分类列表
// @Param authorization header string true "authorization"
// @Param page query int false "page"
// @Param size query int false "size"
// @Param keyword query string false "keyword"
// @Success 200 {string} json "{"code":"200", "msg":"", "data":""}"
// @Router /admin/category-list [get]
func GetCategoryListHandler(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.DefaultQuery("page", constants.DefaultPage))
	if err != nil {
		log.Printf("convert string to int err:%v\n", err.Error())
		return
	}
	size, err := strconv.Atoi(ctx.DefaultQuery("size", constants.DefaultSize))
	if err != nil {
		log.Printf("convert string to int err:%v\n", err.Error())
		return
	}
	keyword := ctx.Query("keyword")

	offset := (page - 1) * size
	list, err := models.GetCategoryList(offset, size, keyword)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": -1,
			"msg":  err.Error(),
		})
		return
	}
	ctx.JSON(200, gin.H{
		"code": 200,
		"msg":  "ok",
		"data": map[string]interface{}{
			"list":  list,
			"count": len(list),
		},
	})
}

// CreateCategoryHandler
// @Tags 管理员方法
// @Summary 创建分类
// @Param authorization header string true "authorization"
// @Param name formData string true "name"
// @Param parent_id formData int false "parent_id"
// @Success 200 {string} json "{"code":"200", "msg":"", "data":""}"
// @Router /admin/category-add [post]
func CreateCategoryHandler(ctx *gin.Context) {
	var categoryBasic models.CategoryBasic
	if err := ctx.ShouldBind(&categoryBasic); err != nil {
		ctx.JSON(200, gin.H{
			"code": -1,
			"msg":  fmt.Sprintf("bind category data err:%v\n", err.Error()),
		})
		return
	}
	categoryBasic.Identity = helper.GetUUID()
	err := models.CreateCategory(&categoryBasic)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": -1,
			"msg":  err.Error(),
		})
		return
	}
	ctx.JSON(200, gin.H{
		"code": 200,
		"msg":  "数据创建成功！",
		"data": categoryBasic,
	})
}

// DeleteCategoryHandler
// @Tags 管理员方法
// @Summary 删除分类
// @Param authorization header string true "authorization"
// @Param identity query string true "identity"
// @Success 200 {string} json "{"code":"200", "msg":"", "data":""}"
// @Router /admin/category-delete [delete]
func DeleteCategoryHandler(ctx *gin.Context) {
	identity := ctx.Query("identity")
	if identity == "" {
		ctx.JSON(200, gin.H{
			"code": -1,
			"msg":  "identity 不能为空！",
		})
		return
	}
	err := models.DeleteCategory(identity)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": -1,
			"msg":  fmt.Sprintf("删除数据失败: %v", err.Error()),
		})
		return
	}
	ctx.JSON(200, gin.H{
		"code": 200,
		"msg":  "删除成功！",
	})
}

// UpdateCategoryHandler
// @Tags 管理员方法
// @Summary 更新分类
// @Param authorization header string true "authorization"
// @Param identity formData string true "identity"
// @Param name formData string false "name"
// @Param parent_id formData int false "parent_id"
// @Success 200 {string} json "{"code":"200", "msg":"", "data":""}"
// @Router /admin/category-update [put]
func UpdateCategoryHandler(ctx *gin.Context) {
	var basic models.CategoryBasic

	if err := ctx.ShouldBind(&basic); err != nil {
		ctx.JSON(200, gin.H{
			"code": -1,
			"msg":  err.Error(),
		})
		return
	}
	err := models.UpdateCategory(&basic)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": -1,
			"msg":  err.Error(),
		})
		return
	}
	ctx.JSON(200, gin.H{
		"code": 200,
		"msg":  "更新成功！",
		"data": basic,
	})
}
