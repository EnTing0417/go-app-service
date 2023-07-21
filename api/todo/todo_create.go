package todo

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/EnTing0417/go-lib/mongodb"
	"github.com/EnTing0417/go-lib/model"
	viewModel "github.com/go-app-service/model"
	"time"
	"github.com/go-openapi/strfmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
 )

 // @BasePath /

// PingExample godoc
// @Summary create a new to-do item
// @Schemes
// @Description create a new to-do item
// @Tags to-do
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token" requestBody body model.ToDoRequestBody true "Request body"
// @Success 200 {object} model.ToDoViewModel
// @Router /api/v1/to-do [post]
// @securityDefinitions.api_key Bearer:<TOKEN>
// @in header
// @name Authorization
func ToDoCreate(c *gin.Context, client *mongo.Client)  {

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

	var requestBody viewModel.ToDoRequestBody

	if err := c.ShouldBindJSON(&requestBody); err!= nil {
		c.JSON(http.StatusBadRequest,gin.H{
			"error": "Request body is required",
		})
		return
	}

	fmt.Printf("email: %v",claimMap["email"])

	criteria := map[string]interface{}{
		"email": claimMap["email"],
	}

	user, _ := mongodb.FindBy(client,mongodb.COLLECTION_USER,criteria)

	if user== nil {
		c.JSON(http.StatusNotFound,gin.H{
			"error": "User is not found",
		})
		return
	}

	if u,ok := user.(primitive.D); ok {

		var _user *model.User

		bsonData, err := bson.Marshal(u)

		if err != nil {
			fmt.Printf("Failed to marshal: %v", err)
		}

		err = bson.Unmarshal(bsonData,&_user)

		if err != nil {
			fmt.Printf("Failed to unmarshal: %v", err)
		}

	toDo := &model.Todo{
		ID : primitive.NewObjectID(),
		UserID: _user.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(), 
		Title: requestBody.Title,
		Description: requestBody.Description,
		Completed: false,
	}
	mongodb.CreateOne(client,mongodb.COLLECTION_TODO,toDo)
	

	response_body := viewModel.ToDoViewModel{
		ID : toDo.ID.Hex(),
		UserID: _user.ID.Hex(),
		CreatedAt: fmt.Sprintf("%v",strfmt.DateTime(toDo.CreatedAt)),
		UpdatedAt:fmt.Sprintf("%v",strfmt.DateTime(toDo.UpdatedAt)),
		Title: toDo.Title,
		Description: toDo.Description,
		Completed: toDo.Completed,
	}
	c.JSON(http.StatusOK,response_body)
	return
}
c.JSON(http.StatusInternalServerError,gin.H{"err":"Internal Server Error"})

 }

