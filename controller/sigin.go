package controller

import (
	"IdentityServer/models"
	"IdentityServer/service"
	"IdentityServer/utils"
	"context"
	"log"
	"net/url"
	"time"

	"github.com/kataras/iris/v12"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type SiginInRequest struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

func SiginInHandler(ctx iris.Context) {
	var body SiginInRequest
	var oauthData AuthorizeBody

	ctx.ReadForm(&body)
	oauthData.ClientId = ctx.URLParam("client_id")
	oauthData.FlowEntry = ctx.URLParam("client_id")
	oauthData.RedirectUrl = ctx.URLParam("redirect_url")

	c, _ := context.WithTimeout(ctx.Request().Context(), 10*time.Second)

	db := ctx.Values().Get("mongoDB").(*mongo.Database)
	defer db.Client().Disconnect(c)

	// basic client authentication
	clientCollection := models.OauthClientCollection(db)

	var client models.OauthClient
	clientQuery := bson.M{
		"key":         oauthData.ClientId,
		"redirectURI": oauthData.RedirectUrl,
	}

	if err := clientCollection.FindOne(c, clientQuery).Decode(&client); err != nil {

		redirectUrl, err := url.Parse(oauthData.RedirectUrl)
		if err != nil {
			log.Println("Invalid Url ")
		}

		q := redirectUrl.Query()
		q.Set("error", "access_denied")

		if oauthData.State != "" {
			q.Set("state", oauthData.State)
		}
		redirectUrl.RawQuery = q.Encode()

		ctx.Redirect(redirectUrl.String())
		return

	}

	user, err := models.GetUserByEmailOrContactNumber(db, body.Username, c)
	if err != nil {
		log.Println("No such user . Check Email or password")
		return
	}

	if err := utils.VerifyPassword(user.Password, body.Password); err != nil {
		log.Println("Invalid Password")
		return
	}

	if user.TwoFector.IsEnable {
		log.Println("Handle Two Factor ")
		return
	} else {
		sessionUser := &service.User{
			Id:    user.Id,
			Email: user.Email,
		}
		go sessionUser.AddTosession(ctx)

		if oauthData.ResponseType == "token" {
			ResponseWithToken(ctx, client, oauthData, *sessionUser, db, c)
			return
		}

		ResponseWithCode(ctx, client, oauthData, *sessionUser, db, c)
		return

		// response with code

	}

}
