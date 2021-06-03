package route

import (
	"IdentityServer/controller"

	"github.com/kataras/iris/v12"
)

func Route(route iris.Party) {

	route.HandleDir("/static", iris.Dir("./static"))
	route.Get("/home", func(ctx iris.Context) {
		ctx.ViewData("title", "Home page")
		ctx.View("home.html")
	})

	v1 := route.Party("/v1")

	v1.Post("/o/oauth/tokens", controller.OauthTokensHandler)

}
