package api

import (
	"golang.org/x/oauth2"
	"net/http"
	"github.com/EnTing0417/go-lib/model"
	"github.com/gin-gonic/gin"
    "log"
)


// @BasePath /
// PingExample godoc
// @Summary google login 
// @Schemes
// @Description google login
// @Tags authentication
// @Accept json
// @Produce json
// @Success 200
// @Router /google/login [get]
var (
    oauthConfig *oauth2.Config
	config *model.Config
)

func InitGoogleAuth() (){
	config = model.ReadConfig()

    oauthConfig = &oauth2.Config{
        ClientID:     config.GoogleOAuth.ClientID,
        ClientSecret: config.GoogleOAuth.ClientSecret,
        RedirectURL:  config.GoogleOAuth.RedirectUrl, 
        Scopes: []string{
            config.GoogleOAuth.Scopes[0],
            config.GoogleOAuth.Scopes[1],
        },
        Endpoint: oauth2.Endpoint{
            AuthURL:  config.GoogleOAuth.AuthURL,
            TokenURL: config.GoogleOAuth.TokenURL,
        },
    }
}

func GoogleLogin(c *gin.Context) {
	InitGoogleAuth()

    oauthState, err := model.GenerateAntiForgeryStateToken()

    if err != nil {
        log.Printf("Failed to generate state token: %v", err)
    }
    url := oauthConfig.AuthCodeURL(oauthState, oauth2.AccessTypeOffline)
    c.Redirect(http.StatusTemporaryRedirect,url)
}
