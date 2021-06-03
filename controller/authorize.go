package controller

import (
	"IdentityServer/models"
	"IdentityServer/service"
	"context"
	"log"
	"net/url"
	"time"

	"github.com/kataras/iris/v12"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// flow_entry = onof --->  add_session,continue, force_password

type AuthorizeBody struct {
	ClientId     string   `url:"client_id"`
	FlowEntry    string   `url:"flow_entry,omitempty"`
	RedirectUrl  string   `url:"redirect_url"`
	ResponseType string   `url:"response_type,omitempty"`
	State        string   `url:"state,omitempty"`
	Scope        []string `url:"scope,omitempty"`
}

func Authorize(ctx iris.Context) {

	// read url query
	var body AuthorizeBody
	ctx.ReadQuery(&body)

	c, _ := context.WithTimeout(ctx.Request().Context(), 10*time.Second)

	db := ctx.Values().Get("mongoDB").(*mongo.Database)
	defer db.Client().Disconnect(c)

	// basic client authentocation

	clientCollection := models.OauthClientCollection(db)

	var client models.OauthClient
	clientQuery := bson.M{
		"key":         body.ClientId,
		"redirectURI": body.RedirectUrl,
	}

	if err := clientCollection.FindOne(c, clientQuery).Decode(&client); err != nil {

		redirectUrl, err := url.Parse(body.RedirectUrl)
		if err != nil {
			log.Println("Invalid Url ")
		}

		q := redirectUrl.Query()
		q.Set("error", "access_denied")

		if body.State != "" {
			q.Set("state", body.State)
		}
		redirectUrl.RawQuery = q.Encode()

		ctx.Redirect(redirectUrl.String())
		return

	}

	if body.FlowEntry == "add_session" {
		ctx.ViewData("Title", "siginin")
		ctx.ViewData("Query", ctx.Request().URL.RawQuery)
		ctx.ViewData("SiginInUrl", "/v1/o/oauth/siginin?"+ctx.Request().URL.RawQuery)
		ctx.View("signin.html")

	} else if body.FlowEntry == "choose_account" {
		ctx.View("choose_account.html")
		ctx.ViewData("query", ctx.Request().URL.RawQuery)
	} else {

		if service.IsAuthenticate(ctx) {
			user := service.GetCurrentUser(ctx)
			if body.ResponseType == "token" {
				ResponseWithToken(ctx, client, body, user, db, c)
				return
			} else {
				ResponseWithCode(ctx, client, body, user, db, c)
				return

			}

		}

		ctx.Redirect("/v1/o/oauth/authorize?flow_entry=add_session&" + ctx.Request().URL.RawQuery)
	}

}
