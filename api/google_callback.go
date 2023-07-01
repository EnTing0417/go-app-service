package api

import (
	"net/http"
	"log"
	"github.com/EnTing0417/go-lib/model"
	"encoding/json"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/EnTing0417/go-lib/mongodb"
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"fmt"
	"github.com/go-openapi/strfmt"
)

func GoogleCallback(c *gin.Context, client *mongo.Client) {
    code := c.Query("code")

    OAuthToken, err := oauthConfig.Exchange(context.Background(), code)
    if err != nil {
        log.Printf("Error exchanging token: %v", err)
		c.JSON(http.StatusInternalServerError,gin.H{
			"error": "Internal Server Error",
		})
        return
    }

    OAuthclient := oauthConfig.Client(context.Background(), OAuthToken)
    
	resp, err := OAuthclient.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		log.Printf("Error getting user information: %v", err)
		c.JSON(http.StatusInternalServerError,gin.H{
			"error": "Internal Server Error",
		})
		return
	}
	defer resp.Body.Close()

	var userInfo model.User

	err = json.NewDecoder(resp.Body).Decode(&userInfo)
	if err != nil {
		log.Printf("Error decoding user information: %v", err)
		c.JSON(http.StatusInternalServerError,gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	criteria := map[string]interface{}{
		"email": userInfo.Email,
		"deleted_at": nil,
	}

	user, err := mongodb.FindBy(client,mongodb.COLLECTION_USER,criteria)

	if user== nil {
		//Create one new user
		_user := &model.User{
			ID : primitive.NewObjectID(),
			Username: userInfo.Email,
			Email: userInfo.Email,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(), 
		}
		mongodb.CreateOne(client,mongodb.COLLECTION_USER,_user)
		user = _user
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
			"username" : userInfo.Email,
			"email": userInfo.Email,
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
		_claims["exp"] = time.Now().Add(time.Hour * 1).Unix()

		refreshTokenString, err := model.GenerateToken(_claims,config.Auth.RefreshTokenSecretKey)
		if err != nil {
			log.Printf("Failed to generate refresh token: %v", err)
			c.JSON(http.StatusInternalServerError,gin.H{
				"error": "Internal Server Error",
			})
			return
		}
		refreshExpireAt := time.Unix(_claims["exp"].(int64), 0)
		
		sess := &model.UserSession{
			ID : primitive.NewObjectID(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(), 
			Token: &model.Token{
				ID : tokenString, 
				ExpireAt: expireAt,
			},
			RefreshToken : &model.Token{
				ID : refreshTokenString, 
				ExpireAt: refreshExpireAt,
			},
			UserID: _user.ID,
			Status: "active",
		}
		err = mongodb.CreateOne(client,mongodb.COLLECTION_USER_SESSION,sess)

		if err != nil {
			log.Fatal(err)
			return
		}

		response := map[string]interface{}{
			"token": tokenString,
			"refresh_token": refreshTokenString,
			"expire_at": fmt.Sprintf("%v",strfmt.DateTime(expireAt)),
		}

		c.JSON(http.StatusOK,response)
		return
	}
	c.JSON(http.StatusInternalServerError,gin.H{"err":"Internal Server Error"})
}
