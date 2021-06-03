package utils

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/dgrijalva/jwt-go"
)

type JwtClaim struct {
	*jwt.StandardClaims
	TokenType string
	ToeknId   string
}

func (claim JwtClaim) EncodeJwt() string {

	signBytes, err := ioutil.ReadFile("./keys/jwt-pri.key")
	if err != nil {
		log.Println("Can Not read private key ")
	}

	signKey, _ := jwt.ParseRSAPrivateKeyFromPEM(signBytes)

	token := jwt.New(jwt.GetSigningMethod("RS256"))
	token.Claims = claim

	result, _ := token.SignedString(signKey)
	return result

}

func DecodeJwt(jwtToken string) (jwt.MapClaims, error) {
	verifyBytes, err := ioutil.ReadFile("./keys/jwt-pub.key")
	if err != nil {
		log.Println("Can Not Open jwt public key ")
	}

	verifyKey, _ := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)

	token, err := jwt.Parse(jwtToken, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", jwtToken.Header["alg"])
		}

		return verifyKey, nil
	})

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok || !token.Valid && err != nil {
		return claims, fmt.Errorf("validate: invalid")
	}
	return claims, nil

}
