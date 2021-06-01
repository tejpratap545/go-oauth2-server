package models

type NewPassword struct {
	Password1 string `json:"password1" xml:"password1" form:"password1"`
	Password2 string `json:"password2" xml:"password2" form:"password2"`
}

type ChangePassword struct {
	OldPassword string `json:"oldPassword" xml:"oldPassword" form:"oldPassword"`
	Password1   string `json:"password1" xml:"password1" form:"password1"`
	Password2   string `json:"password2" xml:"password2" form:"password2"`
}
