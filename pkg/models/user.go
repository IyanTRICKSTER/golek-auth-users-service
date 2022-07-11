package model

import (
	"acourse-auth-user-service/pkg/database"
	"acourse-auth-user-service/pkg/http/requests"
	bcryptUtils "acourse-auth-user-service/pkg/utils/bcrypt"
	tokenUtils "acourse-auth-user-service/pkg/utils/jwt"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"html"
	"log"
	"strings"
	"time"
)

type Language string

const (
	EN Language = "en"
	ID Language = "id"
)

type User struct {
	ID             uint           `gorm:"primary_key" json:"id,omitempty"`
	RoleID         uint           `json:"role_id,omitempty"`
	Role           Role           `gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE" json:"-"`
	Username       string         `gorm:"size:255;not null;unique;constraint:OnUpdate:CASCADE;" json:"username,omitempty"`
	Password       string         `gorm:"size:255;not null;" json:"password,omitempty"`
	Email          string         `gorm:"email;not null;unique;constraint:OnUpdate:CASCADE" json:"email,omitempty"`
	Avatar         string         `json:"avatar,omitempty"`
	PhoneNumber    string         `gorm:"constraint:OnUpdate:CASCADE;not null;unique" json:"phone_number,omitempty"`
	Lang           Language       `json:"language,omitempty"`
	RecommendClass bool           `gorm:"default:true" json:"recommend_class,omitempty"`
	Promotion      bool           `gorm:"default:true" json:"promotion,omitempty"`
	Notification   bool           `gorm:"default:true" json:"notification,omitempty"`
	LatestNews     bool           `gorm:"default:true" json:"latest_news,omitempty"`
	ResetToken     string         `json:"reset_token,omitempty"`
	CreatedAt      time.Time      `json:"created_at,omitempty"`
	UpdatedAt      time.Time      `json:"updated_at,omitempty"`
	DeletedAt      gorm.DeletedAt `json:"deleted_at,omitempty"`
}

func AllUser(limit uint, page uint) (*Pagination, error) {

	var users []*User

	//Paginate Result
	pagination := Pagination{
		Limit: int(limit),
		Page:  int(page),
	}

	database.GetConnection().Select(
		"id", "role_id", "username", "email", "phone_number", "avatar", "lang",
		"recommend_class", "promotion", "notification", "latest_news", "created_at",
		"updated_at", "deleted_at",
	).Scopes(paginate(User{}, &pagination, database.GetConnection())).Find(&users)

	pagination.Rows = users

	return &pagination, nil
}

func (u *User) CreateUser() (*User, error) {

	err := database.GetConnection().Create(&u).Error
	if err != nil {
		log.Println("CREATE USER: ", err.Error())
		return &User{}, err
	}
	log.Println("CREATE USER: OK")
	return u, nil
}

func (u *User) BeforeSave(*gorm.DB) error {

	//turn password into hash
	hashedPassword, err := bcryptUtils.HashPassword(u.Password)
	if err != nil {
		log.Println("BEFORE SAVE USER > HASHING PASSWORD: ", err.Error())
		return err
	}
	u.Password = string(hashedPassword)
	log.Println("BEFORE SAVE USER > HASHIING PASSWORD: OK")

	//Set Default Locale
	u.Lang = "id"

	//Set Default User Role
	u.Role = Role{ID: 2}

	//remove spaces in username
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))

	log.Println("BEFORE SAVE USER: OK")
	return nil

}

func FindUser(uid uint) (User, error) {

	var u User

	if err := database.GetConnection().First(&u, uid).Error; err != nil {
		log.Println("FIND A USER: ", err.Error())
		return u, errors.New("User not found!")
	}

	//Hide Dangerous Field
	u.ExcludeFields()

	log.Println("FIND A USER: OK")
	return u, nil

}

func FindUserByEmail(email string) (User, error) {

	var user User
	if err := database.GetConnection().Where("email = ?", email).First(&user); errors.Is(err.Error, gorm.ErrRecordNotFound) {
		log.Println("FIND A USER BY EMAIL: ", err)
		return user, err.Error
	}

	//Hide Dangerous Field
	user.ExcludeFields()

	log.Println("FIND A USER BY EMAIL: OK")
	return user, nil
}

func (user *User) UpdateUser(data requests.UpdateUserRecordCredential) error {

	if err := database.GetConnection().Model(&user).Find(&user).Updates(User{
		Username:       data.Username,
		Email:          data.Email,
		Avatar:         data.Avatar,
		PhoneNumber:    data.PhoneNumber,
		Lang:           Language(data.Lang),
		RecommendClass: data.RecommendClass,
		Promotion:      data.Promotion,
		Notification:   data.Notification,
		LatestNews:     data.LatestNews,
	}); err != nil {
		log.Println("UPDATE A USER: ", err)
		return err.Error
	}
	log.Println("UPDATE A USER: OK")
	return nil
}

func (user *User) UpdateUserPassword(password string) error {

	hashedPassword, err := bcryptUtils.HashPassword(password)
	if err != nil {
		log.Println("UPDATE USER PASSWORD: ", err.Error())
		return err
	}

	database.GetConnection().Model(&user).Find(&user).Updates(User{Password: string(hashedPassword)})

	log.Println("UPDATE USER PASSWORD: OK")
	return nil

}

func (user *User) DeleteUser() error {
	if err := database.GetConnection().Delete(&user); err != nil {
		log.Println("DELETE USER: ", err)
		return err.Error
	}
	log.Println("DELETE USER: OK")
	return nil
}

func (user *User) IssueResetTokenUser() {

	resetToken, _ := tokenUtils.GenerateResetToken(user.ID)
	database.GetConnection().Model(&user).Updates(User{ResetToken: resetToken})
	log.Println("USER ISSUE RESET TOKEN: OK")
}

func (user *User) RemoveResetTokenUser() {

	database.GetConnection().Model(&user).Update("reset_token", nil)
	log.Println("USER DELETE RESET TOKEN: OK")

}

func VerifyUserPassword(password, hashedPassword string) error {

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		log.Println("USER VERIFY PASSWORD:", err.Error())
		return err
	}
	log.Println("USER VERIFY PASSWORD: OK")
	return nil
}

func AuthenticateUser(email string, password string) (interface{}, error) {

	u := User{}

	err := database.GetConnection().Model(User{}).Where("email = ?", email).Take(&u).Error

	if err != nil {
		log.Println("AUTHENTICATE USER:", err.Error())
		return "", err
	}

	err = VerifyUserPassword(password, u.Password)

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		log.Println("AUTHENTICATE USER:", err.Error())
		return "", err
	}

	token, err := tokenUtils.GenerateAccessToken(u.ID)
	refreshToken, err := tokenUtils.GenerateRefershToken(u.ID)

	if err != nil {
		log.Println("AUTHENTICATE USER:", err.Error())
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
