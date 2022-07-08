package requests

import "errors"

type LoginCredential struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterCredential struct {
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required"`
	Email       string `json:"email" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
}

type ChangePasswordCredential struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordCredential struct {
	Password   string `json:"password" binding:"required,min=8"`
	C_Password string `json:"c_password" binding:"required,min=8"`
}

func (credential *ResetPasswordCredential) ValidateResetPasswordCredential() error {
	if credential.Password != credential.C_Password {
		return errors.New("Provided password is not match")
	}
	return nil
}
