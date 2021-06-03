package utils

import "github.com/kataras/iris/v12"

func InvalidGrantResponse(ctx iris.Context) {
	ctx.StopWithJSON(400, iris.Map{"msg": "Invalid Grant "})

}
