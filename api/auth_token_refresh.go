package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/EnTing0417/go-lib/mongodb"
	"github.com/EnTing0417/go-lib/model"
	viewModel "github.com/go-app-service/model"
	"time"
	"github.com/go-openapi/strfmt"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"go.mongodb.org/mongo-driver/bson/primitive"
 )

 // @BasePath /

// PingExample godoc
// @Summary refresh auth token
// @Schemes
// @Description refresh auth token
// @Tags authentication
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token" requestBody body model.RefreshTokenRequestBody true "Request body"
// @Success 200 {object} model.RefreshTokenResponseBoby
// @Router /api/v1/token/refresh [post]
// @securityDefinitions.api_key Bearer:<TOKEN>
// @in header
// @name Authorization
func AuthTokenRefresh(c *gin.Context, client *mongo.Client)  {

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

	var requestBody viewModel.RefreshTokenRequestBody

	if err := c.ShouldBindJSON(&requestBody); err!= nil {
		c.JSON(http.StatusBadRequest,gin.H{
			"error": "Request body is required",
		})
		return
	}

	oldtoken := c.GetHeader("Authorization")
	_oldtoken := oldtoken[7:len(oldtoken)]

	config := model.ReadConfig()

	//Validate refresh token
	_, isValidRefreshToken := model.IsTokenValid(requestBody.RefreshToken, config.Auth.RefreshTokenSecretKey)
	if !isValidRefreshToken {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
        c.Abort()
        return
	}

	criteria := map[string]interface{}{
		"email": claimMap["email"],
		"deleted_at": nil,
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
			log.Printf("Failed to marshal: %v", err)
		}

		err = bson.Unmarshal(bsonData,&_user)

		if err != nil {
			log.Printf("Failed to unmarshal: %v", err)
		}


		config := model.ReadConfig()

		_claims := map[string]interface{}{
			"username" : _user.Email,
			"email": _user.Email,
			"exp": time.Now().Add(time.Minute * 15).Unix(), 
		}
		expireAt := time.Unix(_claims["exp"].(int64), 0)
		tokenString, err := model.GenerateToken(_claims, config.Auth.SecretKey)

		if err != nil {
			log.Printf("Failed to generate access token: %v", err)
			c.JSON(http.StatusInternalServerError,gin.H{
				"error": "Internal Server Error",
			})
			return
		}

		criteria = bson.M{
			"user_id": _user.ID,
			"token.id":  _oldtoken,
			"refresh_token.id": requestBody.RefreshToken,
			"deleted_at": nil,
		}

		userSession, err := mongodb.FindBy(client,mongodb.COLLECTION_USER_SESSION,criteria)

		if userSession == nil {
			log.Printf("Failed to find user session: %v", err)
			c.JSON(http.StatusNotFound,gin.H{
				"error": "User Session Not Found",
			})
			return
		}

		uSess,ok := userSession.(primitive.D)
		
		if ok {

			var _userSess *model.UserSession
	
			bsonData, err := bson.Marshal(uSess)
	
			if err != nil {
				fmt.Printf("Failed to marshal: %v", err)
			}
	
			err = bson.Unmarshal(bsonData,&_userSess)
	
			if err != nil {
				fmt.Printf("Failed to unmarshal: %v", err)
			}

			criteria = map[string]interface{}{
				"_id": _userSess.ID,
			}
		}

		set := make(map[string]interface{})

		set["token.id"] = tokenString
		set["updated_at"] = time.Now()
		set["token.expire_at"] = expireAt
	
		_, err = mongodb.UpdateBy(client,mongodb.COLLECTION_USER_SESSION,criteria,set)

		if err != nil {
			log.Printf("Failed to update user session: %v", err)
			c.JSON(http.StatusInternalServerError,gin.H{
				"error": "Internal Server Error",
			})
			return
		}

		response := map[string]interface{}{
			"token": tokenString,
			"refresh_token": requestBody.RefreshToken,
			"expire_at": fmt.Sprintf("%v",strfmt.DateTime(expireAt)),
		}
		c.JSON(http.StatusOK,response)
		return
	}

	c.JSON(http.StatusInternalServerError,gin.H{"err":"Internal Server Error"})
 }

