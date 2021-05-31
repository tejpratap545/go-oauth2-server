package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OauthClient struct {
	Key         string    `json:"key,omitempty" bson:"key" xml:"key"`
	Secret      string    `json:"secret,omitempty" bson:"secret" xml:"secret"`
	RedirectURI string    `json:"redirectURI,omitempty" bson:"redirect_url" xml:"redirectURI"`
	Name        string    `json:"name,omitempty" bson:"name,omitempty" xml:"name"`
	Description string    `json:"description,omitempty" bson:"description,omitempty" xml:"description"`
	Logo        string    `json:"logo,omitempty" bson:"logo,omitempty" xml:"logo"`
	CreatedAt   time.Time `json:"createdAt,omitempty" bson:"created_at" xml:"createdAt"`
	DeletedAt   time.Time `json:"deletedAt,omitempty" bson:"deleted_at" xml:"deletedAt"`
	IsActive    bool      `json:"isActive,omitempty" bson:"is_active" xml:"isActive"`
}

type OauthAccessToken struct {
	ClientID  primitive.ObjectID `json:"clientID,omitempty" bson:"client_id" xml:"clientID"`
	Client    OauthClient        `json:"client,omitempty" bson:"_" xml:"client"`
	UserID    primitive.ObjectID `json:"userID,omitempty" bson:"user_id" xml:"userID"`
	User      User               `json:"user,omitempty" bson:"user" xml:"user"`
	Token     string             `json:"token,omitempty" bson:"token" xml:"token"`
	ExpiresAt time.Time          `json:"expiresAt,omitempty" bson:"expires_at" xml:"expiresAt"`
	Scope     string             `json:"scope,omitempty" bson:"scope" xml:"scope"`
}

type OauthRefreshToken struct {
	ClientID  primitive.ObjectID `json:"clientID,omitempty" bson:"client_id" xml:"clientID"`
	Client    OauthClient        `json:"client,omitempty" bson:"client" xml:"client"`
	UserID    primitive.ObjectID `json:"userID,omitempty" bson:"user_id" xml:"userID"`
	User      User               `json:"user,omitempty" bson:"user" xml:"user"`
	Token     string             `json:"token,omitempty" bson:"token" xml:"token"`
	ExpiresAt time.Time          `json:"expiresAt,omitempty" bson:"expires_at" xml:"expiresAt"`
	Scope     string             `json:"scope,omitempty" bson:"scope" xml:"scope"`
}

type OauthAuthorizationCode struct {
	ClientID    primitive.ObjectID `json:"clientID,omitempty" bson:"client_id" xml:"clientID"`
	UserID      primitive.ObjectID `json:"userID,omitempty" bson:"user_id" xml:"userID"`
	Client      OauthClient        `json:"client,omitempty" bson:"client" xml:"client"`
	User        User               `json:"user,omitempty" bson:"user" xml:"user"`
	Code        string             `json:"code,omitempty" bson:"code" xml:"code"`
	RedirectURI string             `json:"redirectURI,omitempty" bson:"redirect_url" xml:"redirectURI"`
	ExpiresAt   time.Time          `json:"expiresAt,omitempty" bson:"expires_at" xml:"expiresAt"`
	Scope       string             `json:"scope,omitempty" bson:"scope" xml:"scope"`
}
