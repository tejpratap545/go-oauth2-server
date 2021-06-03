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

type request struct {
	GrantType    string `json:"grant_type,omitempty" url:"grant_type" xml:"grant_type" form:"grant_type"`
	RefreshToken string `json:"refresh_token,omitempty" url:"refresh_token" xml:"refresh_token" form:"refresh_token"`
	ClientId     string `json:"client_id,omitempty" url:"client_id" xml:"client_id" form:"client_id"`
	ClientSecert string `json:"client_secert,omitempty" url:"client_secert" xml:"client_secert" form:"client_secert"`
}

func RefreshTokenGrant(ctx iris.Context) {
	var request request

	if err := ctx.ReadBody(&request); err != nil {
		utils.InvalidGrantResponse(ctx)
		return
	}

	db := ctx.Values().Get("mongoDB").(*mongo.Database)
	clientCollection := models.OauthClientCollection(db)
	c, _ := context.WithTimeout(ctx.Request().Context(), 10*time.Second)
	var client models.OauthClient
	clientQuery := bson.M{
		"key":    request.ClientId,
		"secret": request.ClientSecert}

	if err := clientCollection.FindOne(c, clientQuery).Decode(&client); err != nil {
		utils.InvalidGrantResponse(ctx)
		return
	}

	claims, err := utils.DecodeJwt(request.RefreshToken)
	if (claims["TokenType"] != "refresh_token" || claims["sub"] != client.Id.Hex()) && err != nil {
		utils.InvalidGrantResponse(ctx)
		return
	}

	refreshTokenCollection := models.OauthRefreshTokenCollection(db)
	refreshTokenId, err := primitive.ObjectIDFromHex(claims["ToeknId"].(string))
	if err != nil {
		log.Println("Invalid Id")
	}

	matchStage := bson.D{{Key: "$match", Value: bson.D{
		{Key: "_id", Value: refreshTokenId}}},
	}

	lookupStage := bson.D{{
		Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "oAuthAccessToken"},
			{Key: "localField", Value: "accessTokenId"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "accessToken"}}},
	}
	unwindStage := bson.D{{
		Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$accessToken"},
			{Key: "preserveNullAndEmptyArrays", Value: false},
		}},
	}

	cursor, err := refreshTokenCollection.Aggregate(c, mongo.Pipeline{matchStage, lookupStage, unwindStage})
	if err != nil {
		log.Fatal(err)
	}
	var refreshTokens []models.OauthRefreshToken
	if err = cursor.All(c, &refreshTokens); err != nil {
		log.Fatal(err)
	}

	refreshToken := refreshTokens[0]

	if !refreshToken.IsValid() {
		utils.InvalidGrantResponse(ctx)
		return
	}

	userAgent := strings.Split(ctx.GetHeader("User-Agent"), "/")[0]
	ip := ctx.GetHeader("X-Forwarded-For")

	if refreshToken.AccessToken.ClientID != client.Id || refreshToken.AccessToken.Ip != ip || refreshToken.AccessToken.UserAgent != userAgent {
		utils.InvalidGrantResponse(ctx)
		return
	}

	go refreshToken.Revoked(c, db)

	currentTime := time.Now()
	userId := refreshToken.UserID

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
	newRefreshToken := &models.OauthRefreshToken{
		Id:            primitive.NewObjectIDFromTimestamp(currentTime),
		ClientID:      client.Id,
		UserID:        userId,
		AccessTokenId: accessToken.Id,
		CreatedAt:     time.Now(),
		ExpiresAt:     currentTime.AddDate(0, 6, 0),
	}

	go newRefreshToken.Create(c, db)
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
		ToeknId:   newRefreshToken.Id.Hex()}

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
