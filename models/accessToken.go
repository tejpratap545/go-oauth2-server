package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type OauthAccessToken struct {
	Id        primitive.ObjectID `json:"id,omitempty" bson:"_id" xml:"id" form:"id"`
	ClientID  primitive.ObjectID `json:"clientID,omitempty" bson:"clientID" xml:"clientID" form:"clientID"`
	Client    OauthClient        `json:"client,omitempty" bson:"client" xml:"client" form:"client"`
	UserID    primitive.ObjectID `json:"userID,omitempty" bson:"userID" xml:"userID" form:"userID"`
	User      User               `json:"user,omitempty" bson:"user" xml:"user" form:"user"`
	CreatedAt time.Time          `json:"createdAt,omitempty" bson:"createdAt" xml:"createdAt" form:"createdAt"`
	ExpiresAt time.Time          `json:"expiresAt,omitempty" bson:"expiresAt" xml:"expiresAt" form:"expiresAt"`
	IsRevoked bool               `json:"isRevoked,omitempty" bson:"isRevoked" xml:"isRevoked" form:"isRevoked"`
	Scope     []string           `json:"scope,omitempty" bson:"scope" xml:"scope" form:"scope"`
	UserAgent string             `json:"userAgent,omitempty" bson:"userAgent" xml:"userAgent" form:"userAgent"`
	Ip        string             `json:"ip,omitempty" bson:"ip" xml:"ip" form:"ip"`
}

func OAuthAccessTokenCollection(db *mongo.Database) *mongo.Collection {
	return db.Collection("oAuthAccessToken")
}

func (accessToken *OauthAccessToken) Create(ctx context.Context, db *mongo.Database) (*mongo.InsertOneResult, error) {
	oAuthAccessTokenCollection := OAuthAccessTokenCollection(db)
	accessToken.Id = primitive.NewObjectID()

	accessToken.CreatedAt = time.Now()
	accessToken.IsRevoked = false

	return oAuthAccessTokenCollection.InsertOne(ctx, accessToken)

}
