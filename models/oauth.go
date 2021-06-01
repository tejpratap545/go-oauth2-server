package models

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log"
	"math"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type OauthClient struct {
	Key               string    `json:"key,omitempty" bson:"key" xml:"key" bjson:"key"`
	Secret            string    `json:"secret,omitempty" bson:"secret" xml:"secret" bjson:"secret"`
	RedirectURI       string    `json:"redirectURI,omitempty" bson:"redirect_url" xml:"redirectURI" bjson:"redirectURI"`
	Name              string    `json:"name,omitempty" bson:"name,omitempty" xml:"name" bjson:"name"`
	Description       string    `json:"description,omitempty" bson:"description,omitempty" xml:"description" bjson:"description"`
	Logo              string    `json:"logo,omitempty" bson:"logo,omitempty" xml:"logo" bjson:"logo"`
	CreatedAt         time.Time `json:"createdAt,omitempty" bson:"created_at" xml:"createdAt" bjson:"createdAt"`
	DeletedAt         time.Time `json:"deletedAt,omitempty" bson:"deleted_at" xml:"deletedAt" bjson:"deletedAt"`
	IsActive          bool      `json:"isActive,omitempty" bson:"is_active" xml:"isActive" bjson:"isActive"`
	IsCoreApplication bool      `json:"isCoreApplication,omitempty" bjson:"is_core_application" xml:"isCoreApplication"`
}

func RandomBase16String(l int) string {
	buff := make([]byte, int(math.Ceil(float64(l)/2)))
	rand.Read(buff)
	str := hex.EncodeToString(buff)
	return str[:l]
}

func (client *OauthClient) Create(db *mongo.Database) (*mongo.InsertOneResult, error) {
	oauthClientCollection := db.Collection("oauthclient")
	client.CreatedAt = time.Now()
	client.IsActive = true
	client.Key = RandomBase16String(16)
	client.Secret = RandomBase16String(32)
	ctx := context.Background()

	result, err := oauthClientCollection.InsertOne(ctx, client)
	if err != nil {
		log.Fatal(err)
		return result, err
	}
	return result, nil

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
