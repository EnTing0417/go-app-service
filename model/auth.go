package model

type RefreshTokenRequestBody struct {
	RefreshToken       string    `json:"refresh_token" bson:"refresh_token"`
}

type RefreshTokenResponseBoby struct {
	Token string  `json:"token" bson:"token"`
	RefreshToken string  `json:"refresh_token" bson:"refresh_token"`
	ExpireAt string	`json:"expire_at" bson:"expire_at"`
}

