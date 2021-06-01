package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type OauthAuthorizationCode struct {
	Id          primitive.ObjectID `json:"id,omitempty" bson:"_id" xml:"id" form:"id"`
	ClientID    primitive.ObjectID `json:"clientID,omitempty" bson:"clientID" xml:"clientID" form:"clientID"`
	UserID      primitive.ObjectID `json:"userID,omitempty" bson:"userID" xml:"userID" form:"userID"`
	Client      OauthClient        `json:"client,omitempty" bson:"client" xml:"client" form:"client"`
	User        User               `json:"user,omitempty" bson:"user" xml:"user" form:"user"`
	CreatedAt   time.Time          `json:"createdAt,omitempty" bson:"createdAt" xml:"createdAt" form:"createdAt"`
	UserAgent   string             `json:"userAgent,omitempty" bson:"userAgent" xml:"userAgent" form:"userAgent"`
	ip          string             `json:"ip,omitempty" bson:"ip" xml:"ip" form:"ip"`
	RedirectURI string             `json:"redirectURI,omitempty" bson:"redirectURI" xml:"redirectURI" form:"redirectURI"`
	ExpiresAt   time.Time          `json:"expiresAt,omitempty" bson:"expiresAt" xml:"expiresAt" form:"expiresAt"`
	Scope       []string           `json:"scope,omitempty" bson:"scope" xml:"scope" form:"scope"`
}

func OauthAuthorizationCodeCollection(db *mongo.Database) *mongo.Collection {
	return db.Collection("OauthAuthorizationCode")
}

func (code *OauthAuthorizationCode) Create(ctx context.Context, db *mongo.Database) (*mongo.InsertOneResult, error) {
	code.CreatedAt = time.Now()
	oauthAuthorizationCodeCollection := OauthAuthorizationCodeCollection(db)
	return oauthAuthorizationCodeCollection.InsertOne(ctx, code)
}
