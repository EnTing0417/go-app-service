package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/EnTing0417/go-lib/mongodb"
	"github.com/EnTing0417/go-lib/model"
	viewModel "github.com/go-app-service/model"
	"github.com/go-openapi/strfmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
 )

 // @BasePath /api/v1

// PingExample godoc
// @Summary list to-do items by user
// @Schemes
// @Description list to-do items by user
// @Tags to-do
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token" user_id path string true "User ID"
// @Success 200 {object} []model.ToDoViewModel
// @Router /api/v1/to-do/list/{user_id} [get]
// @securityDefinitions.api_key Bearer:<TOKEN>
// @in header
// @name Authorization
func ToDoList(c *gin.Context, client *mongo.Client)  {
	
	user_id := c.Param("user_id")

	userId, err := primitive.ObjectIDFromHex(user_id)

	if err != nil {
		fmt.Printf("Err : %v", err)
		c.JSON(http.StatusBadRequest,gin.H{
			"error": "Invalid ID",
		})
		return
	}
	criteria := map[string]interface{}{
		"user_id": userId,
		"deleted_at": nil,
	}

	sort := map[string]interface{}{
		"created_at": -1,
	}

	var todoListResponse []viewModel.ToDoViewModel

	todoList ,err := mongodb.ListBy(client,mongodb.COLLECTION_TODO,criteria,sort)
	if todoList == nil || err != nil {
		c.JSON(http.StatusNotFound,gin.H{
			"error": "Record Not Found",
		})
		return
	}

	if len(todoList)> 0 {
		for _, toDo := range todoList {

			if td, ok :=  toDo.(primitive.D); ok {


			var _todo *model.Todo

			bsonData, err := bson.Marshal(td)

			if err != nil {
				fmt.Printf("Failed to marshal: %v", err)
			}

			err = bson.Unmarshal(bsonData,&_todo)

			if err != nil {
				fmt.Printf("Failed to unmarshal: %v", err)
			}
				
			item := viewModel.ToDoViewModel{
				ID : _todo.ID.Hex(),
				CreatedAt: fmt.Sprintf("%v",strfmt.DateTime(_todo.CreatedAt)),
				UpdatedAt:fmt.Sprintf("%v",strfmt.DateTime(_todo.UpdatedAt)),
				Title: _todo.Title,
				Description: _todo.Description,
				Completed: _todo.Completed,
				UserID : _todo.UserID.Hex(),
			}
			todoListResponse = append(todoListResponse,item)

			}
		}
	}
	c.JSON(http.StatusOK, gin.H{"data": todoListResponse})
 }

