package service

import (
	"fmt"
	"gin_gorm_oj/constants"
	"gin_gorm_oj/models"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

// GetProblemListHandler
// @Tags 公共的方法
// @Summary 获取问题列表
// @Param page query int false "page"
// @Param size query int false "size"
// @Param keyword query string false "keyword"
// @Param category_identity query string false "category_identity"
// @Success 200 {string} json "{"code":"200", "msg":"", "data":""}"
// @Router /problem-list [get]
func GetProblemListHandler(ctx *gin.Context) {
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
	categoryIdentity := ctx.Query("category_identity")
	tx := models.GetProblemList(keyword, categoryIdentity)

	offset := (page - 1) * size
	fmt.Printf("offset: %v\n", offset)
	problemList := make([]models.ProblemBasic, 0)

	var count int64
	err = tx.Debug().Count(&count).Omit("content").Offset(offset).Limit(size).Find(&problemList).Error // 这里，执行了两个sql,Count和Find
	if err != nil {
		log.Printf("find problem list err%v\n", err)
		return
	}

	fmt.Printf("count: %v\n", count)
	ctx.JSON(200, gin.H{
		"code": 200,
		"mgs":  "ok",
		"data": map[string]interface{}{"list": problemList, "count": count},
	})
}

// GetProblemDetailHandler
// @Tags 公共的方法
// @Summary 获取问题详情
// @Param identity query string false "identity"
// @Success 200 {string} json "{"code":"200", "msg":"", "data":""}"
// @Router /problem-detail [get]
func GetProblemDetailHandler(ctx *gin.Context) {
	identity := ctx.Query("identity")
	if identity == "" {
		ctx.JSON(200, gin.H{
			"code": -1,
			"msg":  "identity cant be empty",
		})
		return
	}
	problem, err := models.GetProblem(identity)
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
		"data": problem,
	})
}

// CreateProblemHandler
// @Tags 管理员方法
// @Summary 创建问题
// @Param authorization header string true "authorization"
// @Param title formData string true "title"
// @Param content formData string true "content"
// @Param max_runtime formData string true "max_runtime"
// @Param max_mem formData string true "max_mem"
// @Param category_ids formData []string true "category_ids" collectionFormat(multi)
// @Param test_cases formData []string true "test_cases" collectionFormat(multi)
// @Success 200 {string} json "{"code":"200", "msg":"", "data":""}"
// @Router /admin/problem-add [post]
func CreateProblemHandler(ctx *gin.Context) {
	var problem models.ProblemBasic

	if err := ctx.ShouldBind(&problem); err != nil {
		ctx.JSON(200, gin.H{
			"code": 200,
			"msg":  fmt.Sprintf("bind problem err:%v\n", err.Error()),
		})
		return
	}
	categoryIds := ctx.PostFormArray("category_ids")
	testCases := ctx.PostFormArray("test_cases")
	if len(categoryIds) == 0 || len(testCases) == 0 {
		ctx.JSON(200, gin.H{
			"code": 200,
			"msg":  "问题类别或测试用例不能为空！",
		})
		return
	}
	err := models.CreteProblem(&problem, categoryIds, testCases)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": 200,
			"msg":  err.Error(),
		})
		return
	}
	ctx.JSON(200, gin.H{
		"code": 200,
		"msg":  "ok",
	})
}

// UpdateProblemHandler
// @Tags 管理员方法
// @Summary 更新问题
// @Param authorization header string true "authorization"
// @Param identity formData string true "identity"
// @Param title formData string true "title"
// @Param content formData string true "content"
// @Param max_runtime formData int true "max_runtime"
// @Param max_mem formData int true "max_mem"
// @Param category_ids formData []string false "category_ids" collectionFormat(multi)
// @Param test_cases formData []string false "test_cases" collectionFormat(multi)
// @Success 200 {string} json "{"code":"200", "msg":"", "data":""}"
// @Router /admin/problem-update [put]
func UpdateProblemHandler(ctx *gin.Context) {

	var problem models.ProblemBasic
	if err := ctx.ShouldBind(&problem); err != nil {
		ctx.JSON(200, gin.H{
			"code": -1,
			"msg":  "bind problem data err:" + err.Error(),
		})
		return
	}
	categoryIds := ctx.PostFormArray("category_ids")
	testCases := ctx.PostFormArray("test_cases")
	if len(categoryIds) == 0 || len(testCases) == 0 {
		ctx.JSON(200, gin.H{
			"code": 200,
			"msg":  "问题类别或测试用例不能为空！",
		})
		return
	}
	err := models.UpdateProblem(&problem, categoryIds, testCases)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": -1,
			"msg":  "更新problem失败：" + err.Error(),
		})
		return
	}
	ctx.JSON(200, gin.H{
		"code": 200,
		"msg":  "更新problem成功！",
	})
}
