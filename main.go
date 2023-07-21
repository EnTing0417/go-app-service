package main

import (
	"github.com/EnTing0417/go-lib/mongodb"
	"github.com/gin-gonic/gin"
	auth "github.com/go-app-service/api/auth"
	todo "github.com/go-app-service/api/todo"
	docs "github.com/go-app-service/docs"
	model "github.com/go-app-service/model"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	client := mongodb.Init()
	mongodb.Connect(client)
	defer mongodb.Disconnect(client)

	r := gin.Default()

	protectedRoutes := []string{
		"/to-do",
		"/token/refresh",
	}

	r.Use(model.AuthMiddleware(protectedRoutes))
	docs.SwaggerInfo.BasePath = "/"

	ToDoCreateHandler := func(c *gin.Context) {
		todo.ToDoCreate(c, client)
	}
	ToDoUpdateHandler := func(c *gin.Context) {
		todo.ToDoUpdate(c, client)
	}
	ToDoDeleteHandler := func(c *gin.Context) {
		todo.ToDoDelete(c, client)
	}

	ToDoListHandler := func(c *gin.Context) {
		todo.ToDoList(c, client)
	}

	GoogleCallbackHandler := func(c *gin.Context) {
		auth.GoogleCallback(c, client)
	}

	RefreshTokenHandler := func(c *gin.Context) {
		auth.AuthTokenRefresh(c, client)
	}

	v1 := r.Group("/api/v1")
	{
		v1.POST("/to-do", ToDoCreateHandler)
		v1.PUT("/to-do/:id", ToDoUpdateHandler)
		v1.DELETE("/to-do", ToDoDeleteHandler)
		v1.GET("/to-do/list", ToDoListHandler)
		v1.POST("/token/refresh", RefreshTokenHandler)
	}

	r.GET("/google/login", auth.GoogleLogin)
	r.GET("/google/callback", GoogleCallbackHandler)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.Run(":8080")

}
