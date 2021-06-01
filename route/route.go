package route

import (
	"github.com/kataras/iris/v12"
)

func Route(route iris.Party) {

	route.HandleDir("/static", iris.Dir("./static"))
	route.Get("/home", func(ctx iris.Context) {
		ctx.ViewData("title", "Home page")
		ctx.View("home.html")
	})

}
