package requests

import "errors"

type LoginCredentialRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterCredentialRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	NIM      string `json:"nim" binding:""`
	NIP      string `json:"nip" binding:""`
	Major    string `json:"major" binding:"required"`
}

type ChangePasswordCredentialRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordCredentialRequest struct {
	Password   string `json:"password" binding:"required,min=8"`
	C_Password string `json:"c_password" binding:"required,min=8"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (credential *ResetPasswordCredentialRequest) ValidateResetPasswordCredential() error {
	if credential.Password != credential.C_Password {
		return errors.New("Provided password is not match")
	}
	return nil
}
