package service

import (
	"fmt"
	"gin_gorm_oj/constants"
	"gin_gorm_oj/helper"
	"gin_gorm_oj/models"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"time"
)

// GetUserDetailHandler
// @Tags 公共的方法
// @Summary 获取用户详情
// @Param identity query string false "identity"
// @Success 200 {string} json "{"code":"200", "msg":"", "data":""}"
// @Router /user-detail [get]
func GetUserDetailHandler(ctx *gin.Context) {
	identity := ctx.Query("identity")
	if identity == "" {
		ctx.JSON(200, gin.H{
			"code": 200,
			"msg":  "identity cant be empty",
		})
		return
	}
	detail, err := models.GetUserDetail(identity)
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
		"data": detail,
	})
}

// LoginHandler
// @Tags 公共的方法
// @Summary 用户登录
// @Param username formData string false "username"
// @Param password formData string false "password"
// @Success 200 {string} json "{"code":"200", "msg":"", "data":""}"
// @Router /login [post]
func LoginHandler(ctx *gin.Context) {
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")
	if username == "" || password == "" {
		ctx.JSON(200, gin.H{
			"code": -1,
			"msg":  "用户名或密码不能为空！",
		})
		return
	}
	password = helper.GetMd5(password)
	userInfo, err := models.GetUserInfo(username, password)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": -1,
			"msg":  err.Error(),
		})
		return
	}
	token, err := helper.GenerateToken(userInfo.Identity, userInfo.Name, userInfo.IsAdmin)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": -1,
			"msg":  fmt.Errorf("generate error: %v\n", err.Error()),
		})
		return
	}
	ctx.JSON(200, gin.H{
		"code": 200,
		"msg":  "ok",
		"data": map[string]interface{}{
			"token": token,
		},
	})
}

// SendVerifyCodeHandler
// @Tags 公共的方法
// @Summary 获取验证码
// @Param email formData string true "email"
// @Success 200 {string} json "{"code":"200", "msg":"", "data":""}"
// @Router /verify-code [post]
func SendVerifyCodeHandler(ctx *gin.Context) {
	email := ctx.PostForm("email")
	if email == "" {
		ctx.JSON(200, gin.H{
			"code": -1,
			"msg":  "email 不能为空！",
		})
		return
	}
	var verifyCode = helper.GenerateVerifyCode()
	models.RedisDB.Set(ctx, email, verifyCode, time.Minute*5) // 5分钟
	err := helper.SendVerifyCode(email, verifyCode)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": -1,
			"msg":  fmt.Sprintf("send verify code err:%v\n", err.Error()),
		})
		return
	}
	ctx.JSON(200, gin.H{
		"code": 200,
		"msg":  "验证码发送成功！",
	})
}

// RegisterHandler
// @Tags 公共的方法
// @Summary 用户注册
// @Param email formData string true "email"
// @Param name formData string true "name"
// @Param password formData string true "password"
// @Param phone formData string false "phone"
// @Param code formData string true "code"
// @Success 200 {string} json "{"code":"200", "msg":"", "data":""}"
// @Router /register [post]
func RegisterHandler(ctx *gin.Context) {
	var userBasic models.UserBasic
	if err := ctx.ShouldBind(&userBasic); err != nil {
		ctx.JSON(200, gin.H{
			"code": -1,
			"msg":  err.Error(),
		})
		return
	}
	verifyCode := ctx.PostForm("code")
	if verifyCode == "" {
		ctx.JSON(200, gin.H{
			"code": -1,
			"msg":  "验证码为必填",
		})
		return
	}
	result, err := models.RedisDB.Get(ctx, userBasic.Email).Result()
	if err != nil {
		log.Printf("get verify code err:%v\n", err.Error())
		ctx.JSON(200, gin.H{
			"code": -1,
			"msg":  "验证码过期，请重新获取验证码",
		})
		return
	}
	if result != verifyCode {
		ctx.JSON(200, gin.H{
			"code": -1,
			"msg":  "验证码验证失败，输入有误！",
		})
		return
	}
	err = models.CheckUserExist(userBasic.Email)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": -1,
			"msg":  err.Error(),
		})
		return
	}
	err = models.CreateUser(&userBasic)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": -1,
			"msg":  "用户注册失败！",
		})
		return
	}
	// 生成token
	token, err := helper.GenerateToken(userBasic.Identity, userBasic.Name, userBasic.IsAdmin)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": -1,
			"msg":  "用户注册失败！",
		})
		return
	}
	ctx.JSON(200, gin.H{
		"code": 200,
		"msg":  "用户注册成功",
		"data": map[string]interface{}{
			"token": token,
		},
	})
}

// GetRankListHandler
// @Tags 公共的方法
// @Summary 获取排行榜列表
// @Param page query int false "page"
// @Param size query int false "size"
// @Success 200 {string} json "{"code":"200", "msg":"", "data":""}"
// @Router /rank-list [get]
func GetRankListHandler(ctx *gin.Context) {
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
	tx := models.GetRankList()
	offset := (page - 1) * size
	var userList = make([]models.UserBasic, 0)
	var count int64
	err = tx.Count(&count).Offset(offset).Limit(size).Find(&userList).Error
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": -1,
			"msg":  fmt.Sprintf("get rank list err:%v\n", err.Error()),
		})
		return
	}
	ctx.JSON(200, gin.H{
		"code": 200,
		"msg":  "ok",
		"data": map[string]interface{}{
			"list":  userList,
			"count": count,
		},
	})
}
