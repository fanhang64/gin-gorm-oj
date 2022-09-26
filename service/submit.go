package service

import (
	"bytes"
	"fmt"
	"gin_gorm_oj/constants"
	"gin_gorm_oj/helper"
	"gin_gorm_oj/models"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"log"
	"os/exec"
	"runtime"
	"strconv"
	"sync"
	"time"
)

// GetSubmitListHandler
// @Tags 公共的方法
// @Summary 获取提交列表
// @Param page query int false "page"
// @Param size query int false "size"
// @Param problem_identity query string false "problem_identity"
// @Param user_identity query string false "user_identity"
// @Param status query int false "status"
// @Success 200 {string} json "{"code":"200", "msg":"", "data":""}"
// @Router /submit-list [get]
func GetSubmitListHandler(ctx *gin.Context) {
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
	userIdentity := ctx.Query("user_identity")
	problemIdentity := ctx.Query("problem_identity")

	status, _ := strconv.Atoi(ctx.Query("status"))

	tx := models.GetSubmitList(userIdentity, problemIdentity, status)
	offset := (page - 1) * size
	var count int64
	var submitBasicList []models.SubmitBasic

	if err := tx.Count(&count).Offset(offset).Limit(size).Find(&submitBasicList).Error; err != nil {
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
			"list":  submitBasicList,
			"count": count,
		},
	})
}

// SubmitHandler
// @Tags 用户的方法
// @Summary 提交代码
// @Param authorization header string true "authorization"
// @Param problem_identity query string true "problem_identity"
// @Param code body string true "code"
// @Success 200 {string} json "{"code":"200", "msg":"", "data":""}"
// @Router /user/submit [post]
func SubmitHandler(ctx *gin.Context) {
	problemIdentity := ctx.Query("problem_identity")

	bs, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": -1,
			"msg":  err.Error(),
		})
		return
	}
	// 代码保存到本地
	filePath, err := helper.CodeSave(bs)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": -1,
			"msg":  err.Error(),
		})
		return
	}
	// 提交
	user, exists := ctx.Get("user")
	claims := user.(*helper.UserClaims)
	if !exists {
		ctx.JSON(200, gin.H{
			"code": -1,
			"msg":  err.Error(),
		})
		return
	}
	sb := models.SubmitBasic{
		Identity:        helper.GetUUID(),
		ProblemIdentity: problemIdentity,
		UserIdentity:    claims.Identity,
		Path:            filePath,
	}

	// 代码判断
	pb, err := models.GetProblemPreload(problemIdentity)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": -1,
			"msg":  err.Error(),
		})
		return
	}
	wa := make(chan int)  // 答案错误
	oom := make(chan int) // 内存溢出
	ce := make(chan int)  // 编译错误
	passCount := 0
	var mtx sync.Mutex
	var info string

	for _, testCase := range pb.TestCases {
		testCase := testCase
		go func() {
			// 执行测试
			cmd := exec.Command("go", "run", filePath)
			var out, stdErr bytes.Buffer
			cmd.Stdout = &out
			cmd.Stderr = &stdErr

			pipe, err := cmd.StdinPipe()
			if err != nil {
				log.Println(err)
				return
			}
			io.WriteString(pipe, testCase.Input)

			var start, end runtime.MemStats // 计算内存
			runtime.ReadMemStats(&start)
			// 根据测试的输入案例，进行运行，拿到输出结果和标准的输出结果进行比对
			if err := cmd.Run(); err != nil {
				log.Println(err, stdErr.String())
				if err.Error() == "exit status 2" {
					info = stdErr.String()
					ce <- 1
					return
				}
			}
			// 运行超内存
			runtime.ReadMemStats(&end)
			if end.Alloc/1024-(start.Alloc/1024) > uint64(pb.MaxMem) {
				info = "运行超内存"
				oom <- 1
				return
			}
			// 答案错误
			fmt.Printf("out.String(): %v\n", out.String())
			if testCase.Output != out.String() {
				info = "答案错误"
				wa <- 1
				return
			}

			mtx.Lock()
			passCount++
			mtx.Unlock()
		}()
	}
	select {
	case <-wa:
		sb.Status = 2
	case <-oom:
		sb.Status = 4
	case <-time.After(time.Millisecond * time.Duration(pb.MaxRuntime)):
		if len(pb.TestCases) == passCount { // 答案正确
			sb.Status = 1
		} else {
			sb.Status = 3
		}
	case <-ce:
		sb.Status = 5 // 编译错误
	}
	// 创建submit并更新user和problem的pass_num等字段
	err = models.CreateSubmit(&sb)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": -1,
			"msg":  err.Error(),
		})
		return
	}
	ctx.JSON(200, gin.H{
		"code": 200,
		"msg":  "提交代码成功！",
		"data": map[string]interface{}{
			"status": sb.Status,
			"info":   info,
		},
	})
}
