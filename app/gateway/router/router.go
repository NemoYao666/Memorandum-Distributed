package router

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "micro-todoList-k8s/app/docs"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"

	"micro-todoList-k8s/app/gateway/http"
	"micro-todoList-k8s/app/gateway/middleware"
)

func NewRouter(tracer opentracing.Tracer) *gin.Engine {
	ginRouter := gin.Default()
	//跨域
	ginRouter.Use(middleware.Cors(),
		middleware.TracingMiddleware(tracer),
		middleware.PrometheusMiddleware())

	store := cookie.NewStore([]byte("something-very-secret"))
	ginRouter.Use(sessions.Sessions("mysession", store))
	//http://127.0.0.1:4000/swagger/index.html
	ginRouter.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	v1 := ginRouter.Group("/api/v1")
	{
		v1.GET("ping", func(context *gin.Context) {
			context.JSON(200, "success")
		})
		// 用户服务
		v1.POST("/user/register", http.UserRegisterHandler)
		v1.POST("/user/login", http.UserLoginHandler)

		// 需要登录鉴权保护
		authed := v1.Group("/")
		authed.Use(middleware.JWT())
		{
			authed.GET("tasks", http.ListTaskHandler)
			authed.POST("task", http.CreateTaskHandler)
			authed.GET("task/:id", http.GetTaskHandler)       // task_id
			authed.PUT("task/:id", http.UpdateTaskHandler)    // task_id
			authed.DELETE("task/:id", http.DeleteTaskHandler) // task_id
		}
	}
	return ginRouter
}

// 生成swagger文档
// swag init -g .\app\gateway\cmd\main.go --dir . -o .\app\docs
