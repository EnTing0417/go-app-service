package auth

import (
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/gin-gonic/gin"
)

type HandlerFunc func(c *gin.Context,client *mongo.Client)