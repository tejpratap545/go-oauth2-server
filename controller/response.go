package controller

import (
	"IdentityServer/models"
	"IdentityServer/service"
	"IdentityServer/utils"
	"context"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris/v12"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func ResponseWithToken(ctx iris.Context, client models.OauthClient, body AuthorizeBody, user service.User, db *mongo.Database, c context.Context) {

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

	accessTokenClaim := &utils.JwtClaim{
		StandardClaims: &jwt.StandardClaims{
			Issuer:    client.Id.Hex(),
			Subject:   user.Id.Hex(),
			ExpiresAt: currentTime.Unix() + int64(2*24*time.Hour),
			IssuedAt:  currentTime.Unix(),
		}, TokenType: "access_token", ToeknId: accessToken.Id.Hex()}

	accesstokenString := accessTokenClaim.EncodeJwt()
	redirectUrl, err := url.Parse(body.RedirectUrl)
	if err != nil {
		log.Println("Invalid Url ")
		ctx.Redirect("/error")

	}

	q := redirectUrl.Query()
	q.Set("access_token", accesstokenString)

	if body.State != "" {
		q.Set("state", body.State)
	}
	redirectUrl.RawQuery = q.Encode()

	ctx.Redirect(redirectUrl.String())

}

func ResponseWithCode(ctx iris.Context, client models.OauthClient, body AuthorizeBody, user service.User, db *mongo.Database, c context.Context) {
	userAgent := strings.Split(ctx.GetHeader("User-Agent"), "/")[0]
	currentTime := time.Now()
	code := &models.OauthAuthorizationCode{
		Id:          primitive.NewObjectIDFromTimestamp(currentTime),
		UserID:      user.Id,
		ClientID:    client.Id,
		CreatedAt:   currentTime,
		UserAgent:   userAgent,
		Ip:          ctx.GetHeader("X-Forwarded-For"),
		RedirectURI: body.RedirectUrl,
		ExpiresAt:   time.Now().Add(5 * time.Minute),
	}

	code.Create(c, db)

	redirectUrl, err := url.Parse(body.RedirectUrl)
	if err != nil {
		log.Println("Invalid Url ")
		ctx.Redirect("/error")

	}

	q := redirectUrl.Query()
	q.Set("code", code.Id.Hex())

	if body.State != "" {
		q.Set("state", body.State)
	}
	redirectUrl.RawQuery = q.Encode()

	ctx.Redirect(redirectUrl.String())
}
