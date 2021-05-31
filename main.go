package main

import "github.com/kataras/iris/v12"

func main() {
	app := newApp()

	app.Listen("8080")
}

func newApp() *iris.Application {
	app := iris.New()
	app.Get("/", index)

	return app
}

func index(ctx iris.Context) {
	ctx.HTML("<h1>Index Page</h1>")
}
