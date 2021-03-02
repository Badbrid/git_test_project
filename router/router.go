package router

import (
	"github.com/gin-gonic/gin"
	"github.com/git_test_project/controller"
	"github.com/git_test_project/middleware"
)

func CollectRouter(r *gin.Engine) *gin.Engine {
	r.Use(middleware.CORSMiddleware())
	r.POST("/api/auth/register", controller.Register)

	r.POST("/api/auth/login", controller.Login)

	r.GET("/api/auth/info", middleware.AutoMiddleware(), controller.Info)

	return r
}
