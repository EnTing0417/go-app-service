package main

import (
   "github.com/gin-gonic/gin"
   docs "github.com/go-app-service/docs"
   api "github.com/go-app-service/api"
   swaggerfiles "github.com/swaggo/files"
   ginSwagger "github.com/swaggo/gin-swagger"
   model "github.com/go-app-service/model"
   "github.com/EnTing0417/go-lib/mongodb"
)


func main()  {
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

   ToDoCreateHandler := func(c *gin.Context){
      api.ToDoCreate(c,client)
   }
   ToDoUpdateHandler := func(c *gin.Context){
      api.ToDoUpdate(c,client)
   }
   ToDoDeleteHandler := func(c *gin.Context){
      api.ToDoDelete(c,client)
   }

   ToDoListHandler := func(c *gin.Context){
      api.ToDoList(c,client)
   }

   GoogleCallbackHandler := func(c *gin.Context){
      api.GoogleCallback(c,client)
   }

   RefreshTokenHandler := func(c *gin.Context){
      api.AuthTokenRefresh(c,client)
   }

   v1 := r.Group("/api/v1")
   {
      v1.POST("/to-do", ToDoCreateHandler)
      v1.PUT("/to-do/:id", ToDoUpdateHandler)
      v1.POST("/to-do/delete",ToDoDeleteHandler)
      v1.GET("/to-do/list/:user_id", ToDoListHandler)
      v1.POST("/token/refresh", RefreshTokenHandler)
   }

   r.GET("/google/login", api.GoogleLogin)
   r.GET("/google/callback", GoogleCallbackHandler)

   r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
   r.Run(":8080")

}