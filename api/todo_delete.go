package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/EnTing0417/go-lib/mongodb"
	"github.com/EnTing0417/go-lib/model"
	viewModel "github.com/go-app-service/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/bson"
	"time"
	"go.mongodb.org/mongo-driver/mongo"
	"fmt"
 )

 // @BasePath /api/v1

// PingExample godoc
// @Summary delete single/multiple to-do items
// @Schemes
// @Description delete single/multiple to-do items
// @Tags to-do
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token" requestBody body model.ToDoDeleteRequestBody true "Request body"
// @Success 204 
// @Router /api/v1/to-do/delete [post]
// @securityDefinitions.api_key Bearer:<TOKEN>
// @in header
// @name Authorization
func ToDoDelete(c *gin.Context, client *mongo.Client)   {

	claims, exists := c.Get("claims")
    if !exists {
        c.AbortWithStatus(http.StatusUnauthorized)
        return
    }

	claimMap, ok := claims.(map[string]interface{})
    if !ok {
        c.AbortWithStatus(http.StatusInternalServerError)
        return
    }

	var requestBody viewModel.ToDoDeleteRequestBody

	if err := c.ShouldBindJSON(&requestBody); err!= nil {
		c.JSON(http.StatusBadRequest,gin.H{
			"error": "Request body is required",
		})
		return
	}

	if len(requestBody.ID) == 0 {
		c.JSON(http.StatusBadRequest,gin.H{
			"error": "ID is required",
		})
		return
	}

	criteria := map[string]interface{}{
		"email": claimMap["email"],
		"deleted_at": nil,
	}

	user, err := mongodb.FindBy(client,mongodb.COLLECTION_USER,criteria)

	if err!= nil {
		c.JSON(http.StatusNotFound,gin.H{
			"error": "User is not found",
		})
		return
	}

	u,ok := user.(primitive.D)
	
	if !ok {
		c.JSON(http.StatusNotFound,gin.H{
			"error": "User is not found",
		})
		return
	}

	var _user *model.User

		bsonData, err := bson.Marshal(u)

		if err != nil {
			fmt.Printf("Failed to marshal: %v", err)
		}

		err = bson.Unmarshal(bsonData,&_user)

		if err != nil {
			fmt.Printf("Failed to unmarshal: %v", err)
		}

	var id_list []primitive.ObjectID

	ids := requestBody.ID 

	for _, id := range ids {
		oid, err := primitive.ObjectIDFromHex(id)

		if err != nil {
			c.JSON(http.StatusBadRequest,gin.H{
				"error": "Invalid ID",
			})
			return
		}
		id_list = append(id_list,oid)
	}

	criteria = bson.M{
		"_id": bson.M{
			"$in": id_list,
		},
		"user_id": _user.ID,
	}

	set := bson.M{
		"$set": bson.M{
			"deleted_at": time.Now(), 
		},
	}

	mongodb.DeleteBy(client,mongodb.COLLECTION_TODO,criteria, set)
	c.JSON(http.StatusNoContent,gin.H{})
 }

