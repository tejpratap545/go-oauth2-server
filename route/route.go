package route

import (
	"IdentityServer/controller"

	"github.com/kataras/iris/v12"
)

func Route(route iris.Party) {

	route.HandleDir("/static", iris.Dir("./static"))

	v1 := route.Party("/v1")
	oauth := v1.Party("/o/oauth")
	oauth.Get("/authorize", controller.Authorize)
	oauth.Post("/siginin", controller.SiginInHandler)

	oauth.Post("/tokens", controller.OauthTokensHandler)
}
