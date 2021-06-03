package controller

import (
	"IdentityServer/utils"

	"github.com/kataras/iris/v12"
)

type GrantType struct {
	GrantType string `json:"grant_type,omitempty" url:"grant_type" xml:"grant_type" form:"grant_type"`
}

func OauthTokensHandler(ctx iris.Context) {

	// read grant type
	var grantType GrantType
	ctx.ReadBody(&grantType)

	grantTypes := map[string]func(ctx iris.Context){
		"authorization_code": AuthorizationCodeGrant,
		"password":           PasswordGrant,
		"client_credentials": ClientCredentialsGrant,
		"refresh_token":      RefreshTokenGrant,
	}

	// Check the grant type
	grantHandler, ok := grantTypes[grantType.GrantType]
	if !ok {
		utils.InvalidGrantResponse(ctx)
		return

	}

	grantHandler(ctx)

}
