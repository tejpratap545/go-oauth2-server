package middleware

import (
	"IdentityServer/config"

	"github.com/kataras/iris/v12"
)

func AddMongoToContext(ctx iris.Context) {
	db := config.DB()
	ctx.Values().Set("mongoDB", db)
	ctx.Next()
}
