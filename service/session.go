package service

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/sessions"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id    primitive.ObjectID
	Email string
}

type Session struct {
	Users          []User
	IsAuthenticate bool
	CurrectUser    User
}

func IsAuthenticate(ctx iris.Context) bool {
	session := sessions.Get(ctx)
	return session.Get("IsAuthenticate") == true

}

func GetCurrentUser(ctx iris.Context) User {
	session := sessions.Get(ctx)
	s := session.Get("currectUser")
	myMap := s.(map[string]interface{})
	userId, _ := primitive.ObjectIDFromHex(myMap["Id"].(string))

	return User{Id: userId, Email: myMap["Email"].(string)}
}
func (user User) AddTosession(ctx iris.Context) {
	session := sessions.Get(ctx)

	if session.Get("IsAuthenticate") != true {
		session.Set("currectUser", user)
		session.Set("users", []User{user})
		session.Set("IsAuthenticate", true)
		return
	}

	// users := session.Get("users").([]User)
	// users = append(users, *user)
	// session.Set("users", users)
	// session.Set("currectUser", user)

}
