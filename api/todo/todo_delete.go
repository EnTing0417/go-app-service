package todo

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/EnTing0417/go-lib/model"
	"github.com/EnTing0417/go-lib/mongodb"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// @BasePath /api/v1

// PingExample godoc
// @Summary delete single/multiple to-do items
// @Schemes
// @Description delete single/multiple to-do items
// @Tags to-do
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param _ids query []string true  "ids collection" collectionFormat(csv)
// @Success 204
// @Router /api/v1/to-do [delete]
// @securityDefinitions.api_key Bearer:<TOKEN>
// @in header
// @name Authorization
func ToDoDelete(c *gin.Context, client *mongo.Client) {

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

	_ids := c.Query("_ids")

	if _ids == "" {
		c.JSON(400, gin.H{"error": "_ids parameter is required"})
		return
	}

	ids := strings.Split(_ids, ",")

	if len(ids) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID is required",
		})
		return
	}

	criteria := map[string]interface{}{
		"email": claimMap["email"],
	}

	user, err := mongodb.FindBy(client, mongodb.COLLECTION_USER, criteria)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User is not found",
		})
		return
	}

	u, ok := user.(primitive.D)

	if !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User is not found",
		})
		return
	}

	var _user *model.User

	bsonData, err := bson.Marshal(u)

	if err != nil {
		fmt.Printf("Failed to marshal: %v", err)
	}

	err = bson.Unmarshal(bsonData, &_user)

	if err != nil {
		fmt.Printf("Failed to unmarshal: %v", err)
	}

	var id_list []primitive.ObjectID

	for _, id := range ids {
		oid, err := primitive.ObjectIDFromHex(id)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid ID",
			})
			return
		}
		id_list = append(id_list, oid)
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

	mongodb.DeleteBy(client, mongodb.COLLECTION_TODO, criteria, set)
	c.JSON(http.StatusNoContent, gin.H{})
}
