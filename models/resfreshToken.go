package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type OauthRefreshToken struct {
	Id        primitive.ObjectID `json:"id,omitempty" bson:"_id" xml:"id" form:"id"`
	ClientID  primitive.ObjectID `json:"clientID,omitempty" bson:"clientID" xml:"clientID" form:"clientID"`
	Client    OauthClient        `json:"client,omitempty" bson:"client" xml:"client" form:"client"`
	UserID    primitive.ObjectID `json:"userID,omitempty" bson:"userID" xml:"userID" form:"userID"`
	User      User               `json:"user,omitempty" bson:"user" xml:"user" form:"user"`
	ExpiresAt time.Time          `json:"expiresAt,omitempty" bson:"expiresAt" xml:"expiresAt" form:"expiresAt"`
	CreatedAt time.Time          `json:"createdAt,omitempty" bson:"createdAt" xml:"createdAt" form:"createdAt"`
}

func OauthRefreshTokenCollection(db *mongo.Database) *mongo.Collection {
	return db.Collection("OauthRefreshToken")
}

func (refreshToken *OauthRefreshToken) Create(ctx context.Context, db *mongo.Database) (*mongo.InsertOneResult, error) {
	refreshToken.CreatedAt = time.Now()
	oauthRefreshTokenCollection := OAuthAccessTokenCollection(db)
	return oauthRefreshTokenCollection.InsertOne(ctx, refreshToken)
}
