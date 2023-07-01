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
// @Summary update to-do item
// @Schemes
// @Description update to-do item
// @Tags to-do
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token" requestBody body model.ToDoPatchRequestBody true "Request body" id path string true "ToDoItem ID"
// @Success 200 {object} model.ToDoViewModel
// @Router /api/v1/to-do/{id} [put]
// @securityDefinitions.api_key Bearer:<TOKEN>
// @in header
// @name Authorization
func ToDoUpdate(c *gin.Context, client *mongo.Client)   {

	var requestBody viewModel.ToDoPatchRequestBody
	id := c.Param("id")

	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		c.JSON(http.StatusBadRequest,gin.H{
			"error": "Invalid ID",
		})
		return
	}

	if err := c.ShouldBindJSON(&requestBody); err!= nil {
		c.JSON(http.StatusBadRequest,gin.H{
			"error": "Request body is required",
		})
		return
	}

	criteria := map[string]interface{}{
		"_id": oid,
		"deleted_at": nil,
	}

	toDo, err := mongodb.FindBy(client,mongodb.COLLECTION_TODO,criteria)

	if toDo== nil || err != nil {
		c.JSON(http.StatusNotFound,gin.H{
			"error": "Record not found",
		})
		return
	}

	set :=make(map[string]interface{})

	if requestBody.Completed != nil {
		set["completed"] = *requestBody.Completed
	}

	newtoDo, err := mongodb.UpdateBy(client,mongodb.COLLECTION_TODO,criteria,set)

	if newtoDo== nil || err != nil {
		c.JSON(http.StatusBadRequest,gin.H{
			"error": "Failed to update",
		})
		return
	}


	if _newTodo,ok := newtoDo.(primitive.D); ok {

		var _todo *model.Todo

		bsonData, err := bson.Marshal(_newTodo)

		if err != nil {
			fmt.Printf("Failed to marshal: %v", err)
		}

		err = bson.Unmarshal(bsonData,&_todo)

		if err != nil {
			fmt.Printf("Failed to unmarshal: %v", err)
		}
			
	response_body := viewModel.ToDoViewModel{
		ID : _todo.ID.Hex(),
		UserID: _todo.UserID.Hex(),
		CreatedAt: fmt.Sprintf("%v",strfmt.DateTime(_todo.CreatedAt)),
		UpdatedAt:fmt.Sprintf("%v",strfmt.DateTime(_todo.UpdatedAt)),
		Title: _todo.Title,
		Description: _todo.Description,
		Completed: _todo.Completed,
	}
	c.JSON(http.StatusOK,response_body)
	return
}
c.JSON(http.StatusInternalServerError,gin.H{"err":"Internal Server Error"})

 }

