package controller

import (
	"IdentityServer/models"
	"IdentityServer/utils"
	"context"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris/v12"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthorizationCodeRequestBody struct {
	GrantType    string `json:"grant_type,omitempty" url:"grant_type" xml:"grant_type" form:"grant_type"`
	Code         string `json:"code,omitempty" xml:"code" form:"code"`
	ClientId     string `json:"client_id,omitempty" xml:"client_id" form:"client_id"`
	ClientSecert string `json:"client_secert,omitempty" xml:"client_secert" form:"client_secert"`
	RedirectUrl  string `json:"redirect_url,omitempty" xml:"redirect_url" form:"redirect_url"`
}

func AuthorizationCodeGrant(ctx iris.Context) {
	var body AuthorizationCodeRequestBody
	if err := ctx.ReadBody(&body); err != nil {
		utils.InvalidGrantResponse(ctx)
		return

	}

	db := ctx.Values().Get("mongoDB").(*mongo.Database)
	clientCollection := models.OauthClientCollection(db)
	c, _ := context.WithTimeout(ctx.Request().Context(), 10*time.Second)
	var client models.OauthClient
	clientQuery := bson.M{
		"key":    body.ClientId,
		"secret": body.ClientSecert}

	if err := clientCollection.FindOne(c, clientQuery).Decode(&client); err != nil {
		utils.InvalidGrantResponse(ctx)
		return
	}

	var code models.OauthAuthorizationCode
	codeId, err := primitive.ObjectIDFromHex(body.Code)
	userAgent := strings.Split(ctx.GetHeader("User-Agent"), "/")[0]
	ip := ctx.GetHeader("X-Forwarded-For")
	if err != nil {
		utils.InvalidGrantResponse(ctx)
		return
	}
	codeCollection := models.OauthAuthorizationCodeCollection(db)
	codeQuery := bson.M{
		"_id":         codeId,
		"redirectURI": body.RedirectUrl,
		"clientID":    client.Id,
		"userAgent":   userAgent,
		"ip":          ip,
	}

	if err := codeCollection.FindOne(c, codeQuery).Decode(&code); err != nil {
		log.Println("Not Valid Code")
		utils.InvalidGrantResponse(ctx)
		return
	}

	if !code.IsValid(c, db) {
		log.Println("Expired Code")
		utils.InvalidGrantResponse(ctx)
		return
	}

	currentTime := time.Now()
	userId := code.UserID

	accessToken := &models.OauthAccessToken{
		ClientID:  client.Id,
		Id:        primitive.NewObjectIDFromTimestamp(currentTime),
		ExpiresAt: currentTime.AddDate(0, 0, 2),
		CreatedAt: currentTime,

		UserAgent: userAgent,
		Ip:        ip,
		UserID:    userId,
	}

	go accessToken.Create(c, db)
	refreshToken := &models.OauthRefreshToken{
		Id:            primitive.NewObjectIDFromTimestamp(currentTime),
		ClientID:      client.Id,
		UserID:        userId,
		AccessTokenId: accessToken.Id,
		CreatedAt:     time.Now(),
		ExpiresAt:     currentTime.AddDate(0, 6, 0),
	}

	go refreshToken.Create(c, db)
	accessTokenClaim := &utils.JwtClaim{
		StandardClaims: &jwt.StandardClaims{
			Issuer:    client.Id.Hex(),
			Subject:   userId.Hex(),
			ExpiresAt: currentTime.Unix() + int64(2*24*time.Hour),
			IssuedAt:  currentTime.Unix(),
		}, TokenType: "access_token", ToeknId: accessToken.Id.Hex()}

	refreshTokenclaim := &utils.JwtClaim{
		StandardClaims: &jwt.StandardClaims{
			Issuer:    client.Id.Hex(),
			Subject:   userId.Hex(),
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
			"user_id":       userId.Hex(),
		},
	)

}
