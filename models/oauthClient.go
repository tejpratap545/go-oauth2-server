package models

import (
	"context"
	"time"

	"IdentityServer/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ClientType int

const (
	AuthorizationCode                ClientType = iota //0
	Implicit                                           //1
	ResourceOwnerPasswordCredentials                   //2
	ClientCredentials                                  //3
)

type OauthClient struct {
	Id                       primitive.ObjectID `json:"id,omitempty" bson:"_id" xml:"id" form:"id"`
	Key                      string             `json:"key,omitempty" bson:"key" xml:"key" form:"key"`
	Secret                   string             `json:"secret,omitempty" bson:"secret" xml:"secret" form:"secret"`
	RedirectURI              []string           `json:"redirectURI,omitempty" bson:"redirectURI" xml:"redirectURI" form:"redirectURI"`
	ClientType               ClientType         `json:"clientType,omitempty" bson:"clientType" xml:"clientType" form:"clientType"`
	Name                     string             `json:"name,omitempty" bson:"name" xml:"name" form:"name"`
	Description              string             `json:"description,omitempty" bson:"description" xml:"description" form:"description"`
	Logo                     string             `json:"logo,omitempty" bson:"logo" xml:"logo" form:"logo"`
	CreatedAt                time.Time          `json:"createdAt,omitempty" bson:"createdAt" xml:"createdAt" form:"createdAt"`
	DeletedAt                time.Time          `json:"deletedAt,omitempty" bson:"deletedAt" xml:"deletedAt" form:"deletedAt"`
	IsActive                 bool               `json:"isActive,omitempty" bson:"isActive" xml:"isActive" form:"isActive"`
	IsCoreApplication        bool               `json:"isCoreApplication,omitempty" bson:"isCoreApplication" xml:"isCoreApplication" form:"isCoreApplication"`
	ApplicationPrivacyPolicy string             `json:"applicationPrivacyPolicy,omitempty" bson:"applicationPrivacyPolicy" xml:"applicationPrivacyPolicy" form:"applicationPrivacyPolicy"`
	ApplicationTAC           string             `json:"applicationTAC,omitempty" bson:"applicationTAC" xml:"applicationTAC" form:"applicationTAC"`
	ApplicationHomePage      string             `json:"applicationHomePage,omitempty" bson:"applicationHomePage" xml:"applicationHomePage" form:"applicationHomePage"`
	ContactEmail             []string           `json:"contactEmail,omitempty" bson:"contactEmail" xml:"contactEmail" form:"contactEmail"`
}

func OauthClientCollection(db *mongo.Database) *mongo.Collection {
	return db.Collection("oAuthClient")
}

func (client *OauthClient) Create(ctx context.Context, db *mongo.Database) (*mongo.InsertOneResult, error) {
	client.Id = primitive.NewObjectID()
	client.Key = utils.RandomBase16String(12) + "_" + utils.RandomBase16String(16) + ".apps.feblicusercontent.com"
	client.Secret = utils.GenerateRandomSecret(32)
	oauthClientCollection := OauthClientCollection(db)
	client.CreatedAt = time.Now()
	client.IsActive = true

	return oauthClientCollection.InsertOne(ctx, client)

}
