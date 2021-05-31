package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OauthClient struct {
	Key         string `json:"key,omitempty" bjson:"key" xml:"key"`
	Secret      string `json:"secret,omitempty" bjson:"secret" xml:"secret"`
	RedirectURI string `json:"redirectURI,omitempty" bjson:"redirectURI" xml:"redirectURI"`
}

type OauthAccessToken struct {
	ClientID  primitive.ObjectID `json:"clientID,omitempty" bjson:"clientID" xml:"clientID"`
	Client    OauthClient        `json:"client,omitempty" bjson:"_" xml:"client"`
	UserID    primitive.ObjectID `json:"userID,omitempty" bjson:"userID" xml:"userID"`
	User      User               `json:"user,omitempty" bjson:"user" xml:"user"`
	Token     string             `json:"token,omitempty" bjson:"token" xml:"token"`
	ExpiresAt time.Time          `json:"expiresAt,omitempty" bjson:"expiresAt" xml:"expiresAt"`
	Scope     string             `json:"scope,omitempty" bjson:"scope" xml:"scope"`
}

type OauthRefreshToken struct {
	ClientID  primitive.ObjectID `json:"clientID,omitempty" bjson:"clientID" xml:"clientID"`
	Client    OauthClient        `json:"client,omitempty" bjson:"client" xml:"client"`
	UserID    primitive.ObjectID `json:"userID,omitempty" bjson:"userID" xml:"userID"`
	User      User               `json:"user,omitempty" bjson:"user" xml:"user"`
	Token     string             `json:"token,omitempty" bjson:"token" xml:"token"`
	ExpiresAt time.Time          `json:"expiresAt,omitempty" bjson:"expiresAt" xml:"expiresAt"`
	Scope     string             `json:"scope,omitempty" bjson:"scope" xml:"scope"`
}

type OauthAuthorizationCode struct {
	ClientID    primitive.ObjectID `json:"clientID,omitempty" bjson:"clientID" xml:"clientID"`
	UserID      primitive.ObjectID `json:"userID,omitempty" bjson:"userID" xml:"userID"`
	Client      OauthClient        `json:"client,omitempty" bjson:"client" xml:"client"`
	User        User               `json:"user,omitempty" bjson:"user" xml:"user"`
	Code        string             `json:"code,omitempty" bjson:"code" xml:"code"`
	RedirectURI string             `json:"redirectURI,omitempty" bjson:"redirectURI" xml:"redirectURI"`
	ExpiresAt   time.Time          `json:"expiresAt,omitempty" bjson:"expiresAt" xml:"expiresAt"`
	Scope       string             `json:"scope,omitempty" bjson:"scope" xml:"scope"`
}
