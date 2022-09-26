package router

import (
	"gin_gorm_oj/middlewares"
	"gin_gorm_oj/service"

	_ "gin_gorm_oj/docs"

	"github.com/gin-gonic/gin"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middlewares.Cors())
	// 配置路由规则
	r.GET("/ping", service.PingHandler)
	// 问题
	r.GET("/problem-list", service.GetProblemListHandler)
	r.GET("/problem-detail", service.GetProblemDetailHandler)
	r.GET("/category-list", service.GetCategoryListHandler)

	// 用户
	r.GET("/user-detail", service.GetUserDetailHandler)
	r.POST("/login", service.LoginHandler)
	r.POST("/verify-code", service.SendVerifyCodeHandler)
	r.POST("/register", service.RegisterHandler)

	// 排行榜
	r.GET("/rank-list", service.GetRankListHandler)

	// 提交记录
	r.GET("/submit-list", service.GetSubmitListHandler)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// 管理员方法
	admin := r.Group("/admin", middlewares.AuthAdminCheck())
	{
		admin.POST("/problem-add", service.CreateProblemHandler)
		admin.POST("/category-add", service.CreateCategoryHandler)
		admin.DELETE("/category-delete", service.DeleteCategoryHandler)
		admin.PUT("/category-update", service.UpdateCategoryHandler)
		admin.PUT("/problem-update", service.UpdateProblemHandler)
	}

	user := r.Group("/user", middlewares.AuthUserCheck())
	{
		user.POST("/submit", service.SubmitHandler)
	}
	return r
}
