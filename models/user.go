package models

import (
	"IdentityServer/config"
	"IdentityServer/utils"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Address struct {
	Address1 string `json:"address1,omitempty" bson:"address1" xml:"address1" form:"address1"`
	Address2 string `json:"address2,omitempty" bson:"address2" xml:"address2" form:"address2"`
	City     string `json:"city,omitempty" bson:"city" xml:"city" form:"city"`
	State    string `json:"state,omitempty" bson:"state" xml:"state" form:"state"`
	Country  string `json:"country,omitempty" bson:"country" xml:"country" form:"country"`
}
type CreditCard struct {
	CardNumner string    `json:"cardNumner,omitempty" bson:"cardNumner" xml:"cardNumner" form:"cardNumner"`
	ExpiryDate time.Time `json:"expiryDate,omitempty" bson:"expiryDate" xml:"expiryDate" form:"expiryDate"`
	CreatedAt  time.Time `json:"createdAt,omitempty" bson:"createdAt" xml:"createdAt" form:"createdAt"`
	IsDefault  bool      `json:"isDefault,omitempty" bson:"isDefault" xml:"isDefault" form:"isDefault"`
}

type TwoFactorApplcation struct {
	Secert    string    `json:"secert,omitempty" bson:"secert" xml:"secert" form:"secert"`
	CreatedAT time.Time `json:"createdAT,omitempty" bson:"createdAT" xml:"createdAT" form:"createdAT"`
}
type TwoFector struct {
	IsEnable      bool                  `json:"isEnable,omitempty" bson:"isEnable" xml:"isEnable" form:"isEnable"`
	Email         []string              `json:"email,omitempty" bson:"email" xml:"email" form:"email"`
	ContactNumber []string              `json:"contactNumber,omitempty" bson:"contactNumber" xml:"contactNumber" form:"contactNumber"`
	Application   []TwoFactorApplcation `json:"application,omitempty" bson:"application" xml:"application" form:"application"`
}
type User struct {
	Id            primitive.ObjectID `json:"id,omitempty" bson:"_id" xml:"id" form:"id"`
	FirstName     string             `json:"firstName,omitempty" bson:"firstName" xml:"firstName" form:"firstName"`
	LastName      string             `json:"lastName,omitempty" bson:"lastName" xml:"lastName" form:"lastName"`
	Email         string             `json:"email,omitempty" bson:"email" xml:"email" form:"email"`
	Password      string             `json:"password,omitempty" bson:"password" xml:"password" form:"password"`
	Avatar        string             `json:"avatar,omitempty" bson:"avatar" xml:"avatar" form:"avatar"`
	ContactNumber string             `json:"contactNumber,omitempty" bson:"contactNumber" xml:"contactNumber" form:"contactNumber"`
	Address       []Address          `json:"address,omitempty" bson:"address" xml:"address" form:"address"`
	Card          []CreditCard       `json:"card,omitempty" bson:"card" xml:"card" form:"card"`
	TwoFector     TwoFector          `json:"twoFector,omitempty" bson:"twoFector" xml:"twoFector" form:"twoFector"`
}

func UserCollection(db *mongo.Database) *mongo.Collection {
	return db.Collection("User")
}

func GetUserByEmailOrContactNumber(db *mongo.Database, username string, c context.Context) (*User, error) {
	var user User
	userCollection := UserCollection(db)

	query := bson.M{
		"$or": bson.A{
			bson.M{"email": username},
			bson.M{"contactNumber": username},
		},
	}
	if err := userCollection.FindOne(c, query).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil

}
func (user *User) Create(db *mongo.Database) (*mongo.InsertOneResult, error) {
	user.Id = primitive.NewObjectIDFromTimestamp(time.Now())

	passwordByte, _ := utils.HashPassword(user.Password)
	user.Password = string(passwordByte)
	userCollection := UserCollection(config.DB())
	return userCollection.InsertOne(context.Background(), user)

}
