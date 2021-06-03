package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type OauthRefreshToken struct {
	Id            primitive.ObjectID `json:"id,omitempty" bson:"_id" xml:"id" form:"id"`
	ClientID      primitive.ObjectID `json:"clientID,omitempty" bson:"clientID" xml:"clientID" form:"clientID"`
	Client        *OauthClient       `json:"client,omitempty" bson:"client" xml:"client" form:"client"`
	UserID        primitive.ObjectID `json:"userID,omitempty" bson:"userID" xml:"userID" form:"userID"`
	User          *User              `json:"user,omitempty" bson:"user" xml:"user" form:"user"`
	AccessTokenId primitive.ObjectID `json:"accessTokenId,omitempty" bson:"accessTokenId" xml:"accessTokenId" form:"accessTokenId"`
	AccessToken   *OauthAccessToken  `json:"accessToken,omitempty" bson:"accessToken" xml:"accessToken" form:"accessToken"`
	ExpiresAt     time.Time          `json:"expiresAt,omitempty" bson:"expiresAt" xml:"expiresAt" form:"expiresAt"`
	CreatedAt     time.Time          `json:"createdAt,omitempty" bson:"createdAt" xml:"createdAt" form:"createdAt"`
	IsRevoked     bool               `json:"isRevoked,omitempty" bson:"isRevoked" xml:"isRevoked" form:"isRevoked"`
}

func OauthRefreshTokenCollection(db *mongo.Database) *mongo.Collection {
	return db.Collection("OauthRefreshToken")
}

func (refreshToken *OauthRefreshToken) IsValid() bool {
	return !refreshToken.IsRevoked && time.Now().Unix() < refreshToken.ExpiresAt.Unix()
}

func (refreshToken *OauthRefreshToken) Create(ctx context.Context, db *mongo.Database) (*mongo.InsertOneResult, error) {
	oauthRefreshTokenCollection := OauthRefreshTokenCollection(db)
	refreshToken.Client = nil
	refreshToken.AccessToken = nil
	refreshToken.AccessToken = nil

	return oauthRefreshTokenCollection.InsertOne(ctx, refreshToken)
}

func (refreshToken *OauthRefreshToken) Revoked(ctx context.Context, db *mongo.Database) {
	oauthRefreshTokenCollection := OauthRefreshTokenCollection(db)
	oauthRefreshTokenCollection.UpdateByID(ctx, refreshToken.Id, bson.M{"$set": bson.M{"isRevoked": true}})

	oAuthAccessTokenCollection := OAuthAccessTokenCollection(db)
	oAuthAccessTokenCollection.UpdateByID(ctx, refreshToken.AccessTokenId, bson.M{"$set": bson.M{"isRevoked": true}})
}
