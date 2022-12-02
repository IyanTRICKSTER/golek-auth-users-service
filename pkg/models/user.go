package model

import (
	"errors"
	"golek-auth-user-service/pkg/database"
	"golek-auth-user-service/pkg/http/requests"
	bcryptUtils "golek-auth-user-service/pkg/utils/bcrypt"
	tokenUtils "golek-auth-user-service/pkg/utils/jwt"
	"html"
	"log"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID         uint           `gorm:"primary_key" json:"id,omitempty"`
	RoleID     uint           `json:"-"`
	Role       Role           `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	Username   string         `gorm:"size:255;not null;unique;constraint:OnUpdate:CASCADE;" json:"username,omitempty"`
	Password   string         `gorm:"size:255;not null;" json:"password,omitempty"`
	Email      string         `gorm:"email;not null;unique;constraint:OnUpdate:CASCADE" json:"email,omitempty"`
	Avatar     string         `json:"avatar,omitempty"`
	NIM        *string        `gorm:"constraint:OnUpdate:CASCADE;null;unique" json:"nim"`
	NIP        *string        `gorm:"constraint:OnUpdate:CASCADE;null;unique" json:"nip"`
	Major      string         `gorm:"not null" json:"major"`
	ResetToken string         `json:"-"`
	CreatedAt  time.Time      `json:"created_at,omitempty"`
	UpdatedAt  time.Time      `json:"updated_at,omitempty"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at,omitempty"`
}

func AllUser(limit uint, page uint) (*Pagination, error) {

	var users []*User

	//Paginate Result
	pagination := Pagination{
		Limit: int(limit),
		Page:  int(page),
	}

	database.GetConnection().Select(
		"id", "role_id", "username", "email", "major", "avatar", "nim", "n_ip", "created_at",
		"updated_at", "deleted_at",
	).Scopes(paginate(User{}, &pagination, database.GetConnection())).Find(&users)

	pagination.Rows = users

	return &pagination, nil
}

func (u *User) CreateUser() (*User, error) {

	err := database.GetConnection().Create(&u).Error
	if err != nil {
		log.Println("CREATE USER ERROR: ", err.Error())
		return &User{}, err
	}
	log.Println("CREATE USER: OK")
	return u, nil
}

func (u *User) BeforeSave(*gorm.DB) error {

	//turn password into hash
	hashedPassword, err := bcryptUtils.HashPassword(u.Password)
	if err != nil {
		log.Println("BEFORE SAVE USER > HASHING PASSWORD ERROR: ", err.Error())
		return err
	}
	u.Password = string(hashedPassword)
	log.Println("BEFORE SAVE USER > HASHIING PASSWORD: OK")

	//Set Default User Role
	u.Role = Role{ID: 2, Name: "member"}

	//remove spaces in username
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))

	if u.NIP != nil {
		if *u.NIP == "" {
			u.NIP = nil
		}
	}

	if u.NIM != nil {
		if *u.NIM == "" {
			u.NIM = nil
		}
	}

	log.Println("BEFORE SAVE USER: OK")
	return nil

}

func FindUser(uid uint) (User, error) {

	var u User

	if err := database.GetConnection().Preload("Role.Permissions", func(tx *gorm.DB) *gorm.DB {
		return tx.Select("id", "name")
	}).First(&u, uid).Error; err != nil {
		log.Println("FIND A USER ERROR: ", err.Error())
		return u, errors.New("user not found")
	}

	//Hide Dangerous Field
	u.ExcludeFields()

	log.Println("FIND A USER: OK")

	return u, nil

}

func IntrospectUser(uid uint, resourse string) (User, error) {

	var u User

	if err := database.GetConnection().Preload("Role.Permissions", func(tx *gorm.DB) *gorm.DB {
		return tx.Select("id", "name", "code", "resource").Where("resource = ?", resourse)
	}).First(&u, uid).Error; err != nil {
		log.Println("INTROSPECT USER ERROR: ", err.Error())
		return u, errors.New("user not found")
	}

	//Hide Dangerous Field
	u.ExcludeFields()

	log.Println("INTROSPECT USER: OK")

	return u, nil

}

func FindUserByEmail(email string) (User, error) {

	var user User
	if err := database.GetConnection().Where("email = ?", email).First(&user); errors.Is(err.Error, gorm.ErrRecordNotFound) {
		log.Println("FIND A USER BY EMAIL ERROR: ", err)
		return user, err.Error
	}

	//Hide Dangerous Field
	user.ExcludeFields()

	log.Println("FIND A USER BY EMAIL: OK")
	return user, nil
}

func (u *User) UpdateUser(data requests.UpdateUserRecordCredential) error {

	if err := database.GetConnection().Model(&u).Find(&u).Updates(User{
		Username: data.Username,
		Email:    data.Email,
		Avatar:   data.Avatar,
		NIP:      &data.NIP,
		NIM:      &data.NIM,
		Major:    data.Major,
	}).Error; err != nil {
		log.Println("UPDATE A USER ERROR: ", err.Error())
		return err
	}
	log.Println("UPDATE A USER: OK")
	return nil
}

func (u *User) UpdateUserPassword(password string) error {

	hashedPassword, err := bcryptUtils.HashPassword(password)
	if err != nil {
		log.Println("UPDATE USER PASSWORD ERROR: ", err.Error())
		return err
	}

	database.GetConnection().Model(&u).Find(&u).Updates(User{Password: string(hashedPassword)})

	log.Println("UPDATE USER PASSWORD: OK")
	return nil

}

func (u *User) DeleteUser() error {
	if err := database.GetConnection().Delete(&u).Error; err != nil {
		log.Println("DELETE USER ERROR: ", err.Error())
		return err
	}
	log.Println("DELETE USER: OK")
	return nil
}

func (u *User) IssueResetTokenUser() {

	resetToken, _ := tokenUtils.GenerateResetToken(u.ID)
	database.GetConnection().Model(&u).Updates(User{ResetToken: resetToken})
	log.Println("USER ISSUE RESET TOKEN: OK")
}

func (u *User) RemoveResetTokenUser() {

	database.GetConnection().Model(&u).Update("reset_token", nil)
	log.Println("USER DELETE RESET TOKEN: OK")

}

func VerifyUserPassword(password, hashedPassword string) error {

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		log.Println("USER VERIFY PASSWORD ERROR:", err.Error())
		return err
	}
	log.Println("USER VERIFY PASSWORD: OK")
	return nil
}

func AuthenticateUser(email string, password string) (interface{}, error) {

	u := User{}

	err := database.GetConnection().Model(User{}).Where("email = ?", email).Take(&u).Error

	if err != nil {
		log.Println("AUTHENTICATE USER ERROR:", err.Error())
		return "", err
	}

	err = VerifyUserPassword(password, u.Password)

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		log.Println("AUTHENTICATE USER ERROR:", err.Error())
		return "", err
	}

	token, err := tokenUtils.GenerateAccessToken(u.ID)
	refreshToken, err := tokenUtils.GenerateRefreshToken(u.ID)

	if err != nil {
		log.Println("AUTHENTICATE USER ERROR:", err.Error())
		return "", err
	}

	pairToken := map[string]interface{}{
		"access_token":  token,
		"refresh_token": refreshToken,
	}

	log.Println("AUTHENTICATE USER: OK")
	return pairToken, nil

}

func (u *User) ExcludeFields() {
	u.Password = ""
}
