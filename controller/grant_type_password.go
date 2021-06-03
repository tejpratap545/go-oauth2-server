package controller

import (
	"IdentityServer/models"
	"IdentityServer/utils"
	"context"
	"strings"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris/v12"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PasswordBody struct {
	GrantType    string `json:"grant_type,omitempty" url:"grant_type" xml:"grant_type" form:"grant_type"`
	ClientId     string `json:"client_id,omitempty" url:"client_id" xml:"client_id" form:"client_id"`
	ClientSecert string `json:"client_secert,omitempty" url:"client_secert" xml:"client_secert" form:"client_secert"`
	Username     string `json:"username,omitempty" url:"username" xml:"username" form:"username"`
	Password     string `json:"password,omitempty" url:"password" xml:"password" form:"password"`
	// Scope        []string `json:"scope,omitempty" xml:"scope" form:"scope"`
}

func PasswordGrant(ctx iris.Context) {

	var body PasswordBody

	if err := ctx.ReadBody(&body); err != nil {
		utils.InvalidGrantResponse(ctx)
		return
	}
	c, _ := context.WithTimeout(ctx.Request().Context(), 10*time.Second)

	db := ctx.Values().Get("mongoDB").(*mongo.Database)
	defer db.Client().Disconnect(c)

	clientCollection := models.OauthClientCollection(db)

	var client models.OauthClient
	clientQuery := bson.M{
		"key":    body.ClientId,
		"secret": body.ClientSecert}

	if err := clientCollection.FindOne(c, clientQuery).Decode(&client); err != nil {
		utils.InvalidGrantResponse(ctx)
		return
	}

	user, err := models.GetUserByEmailOrContactNumber(db, body.Username, c)
	if err != nil {
		utils.InvalidGrantResponse(ctx)
		return
	}

	if err := utils.VerifyPassword(user.Password, body.Password); err != nil {
		utils.InvalidGrantResponse(ctx)
		return
	}
	userAgent := ctx.GetHeader("User-Agent")
	currentTime := time.Now()

	accessToken := &models.OauthAccessToken{
		ClientID:  client.Id,
		Id:        primitive.NewObjectIDFromTimestamp(currentTime),
		ExpiresAt: currentTime.AddDate(0, 0, 2),
		CreatedAt: currentTime,

		UserAgent: strings.Split(userAgent, "/")[0],
		Ip:        ctx.GetHeader("X-Forwarded-For"),
		UserID:    user.Id,
	}

	go accessToken.Create(c, db)
	refreshToken := &models.OauthRefreshToken{
		Id:            primitive.NewObjectIDFromTimestamp(currentTime),
		ClientID:      client.Id,
		UserID:        user.Id,
		AccessTokenId: accessToken.Id,
		CreatedAt:     time.Now(),
		ExpiresAt:     currentTime.AddDate(0, 6, 0),
	}

	go refreshToken.Create(c, db)
	accessTokenClaim := &utils.JwtClaim{
		StandardClaims: &jwt.StandardClaims{
			Issuer:    client.Id.Hex(),
			Subject:   user.Id.Hex(),
			ExpiresAt: currentTime.Unix() + int64(2*24*time.Hour),
			IssuedAt:  currentTime.Unix(),
		}, TokenType: "access_token", ToeknId: accessToken.Id.Hex()}

	refreshTokenclaim := &utils.JwtClaim{
		StandardClaims: &jwt.StandardClaims{
			Issuer:    client.Id.Hex(),
			Subject:   user.Id.Hex(),
			ExpiresAt: currentTime.Unix() + int64(6*30*24*time.Hour),
			IssuedAt:  currentTime.Unix(),
		},
		TokenType: "refresh_token",
		ToeknId:   refreshToken.Id.Hex()}

	var wg sync.WaitGroup
	wg.Add(2)
	var accesstokenString, refreshTokenString string

	go func() {
		accesstokenString = accessTokenClaim.EncodeJwt()
		wg.Done()
	}()
	go func() {
		refreshTokenString = refreshTokenclaim.EncodeJwt()
		wg.Done()
	}()

	wg.Wait()

	ctx.JSON(
		iris.Map{
			"access_token":  accesstokenString,
			"refresh_token": refreshTokenString,
			"expires_in":    2 * time.Hour.Seconds(),
			"user_id":       user.Id.Hex(),
		},
	)

}
