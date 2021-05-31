package models

import (
	"time"
)

type Address struct {
	Address1 string `json:"address1,omitempty" bjson:"address1" xml:"address1"`
	Address2 string `json:"address2,omitempty" bjson:"address2" xml:"address2"`
	City     string `json:"city,omitempty" bjson:"city" xml:"city"`
	State    string `json:"state,omitempty" bjson:"state" xml:"state"`
	Country  string `json:"country,omitempty" bjson:"country" xml:"country"`
}

type CreditCard struct {
	CardNumner string    `json:"cardNumner,omitempty" bjson:"card_numner" xml:"cardNumner"`
	ExpiryDate time.Time `json:"expiryDate,omitempty" bjson:"expiry_date" xml:"expiryDate"`
	CreatedAt  time.Time `json:"createdAt,omitempty" bjson:"created_at" xml:"createdAt"`
	IsDefault  bool      `json:"isDefault,omitempty" bjson:"is_default" xml:"isDefault"`
}
type User struct {
	Name          string       `json:"name,omitempty" bjson:"name" xml:"name"`
	Email         string       `json:"email,omitempty" bjson:"email" xml:"email"`
	Password      string       `json:"_" bjson:"password" xml:"_"`
	ContactNumber string       `json:"contactNumber,omitempty" bjson:"contact_number" xml:"contactNumber"`
	Address       []Address    `json:"address,omitempty" bjson:"address" xml:"address"`
	Card          []CreditCard `json:"card,omitempty" bjson:"card" xml:"card"`
}
