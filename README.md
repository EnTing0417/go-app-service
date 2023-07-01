# go-app-service

# Development Enviroment
- OS : Debian 11
- IDE : Visual Studio Code
- Language : Go

# To run the program
- Make sure `docker compose` is installed in your machine
- Execute `docker compose up`

# To test the app
1. Obtain bearer token from the callback of GET http://localhost:8080/google/login after the server is up

2. Use the bearer token to access the following api:
- POST http://localhost:8080/api/v1/todo
- PUT http://localhost:8080/api/v1/todo/:id
- POST http://localhost:8080/api/v1/to-do/delete
- GET http://localhost:8080/api/v1/to-do-list/:user_id

3. To refresh the token for specific session: 
- POST http://localhost:8080/api/v1/token/refresh

# To build the app
- Install `make` command and add to path variable

For Linux:
`sudo apt-get update`
`sudo apt-get -y install make`

Exec `make`

# To add new api
1. Create a new file and named as `xxx_api.go`  
2. Add the following new code snippet in the file.Example:
`package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
 )
 
``// @BasePath /api/v1
// PingExample godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} Helloworld
// @Router /helloworld [get]
func Helloworld(c *gin.Context)  {
	c.JSON(http.StatusOK,gin.H{
		"message": "hello world",
	})
 }`
`
3. `go get -u github.com/swaggo/swag/cmd/swag`
4. `export PATH=$(go env GOPATH)/bin:$PATH` OR update .bashrc for path variables
5. `swag init`
