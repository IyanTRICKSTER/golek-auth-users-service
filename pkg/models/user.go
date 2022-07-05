package UserModel

import (
	"acourse-auth-user-service/pkg/database"
	bcryptUtils "acourse-auth-user-service/pkg/utils/bcrypt"
	tokenUtils "acourse-auth-user-service/pkg/utils/jwt"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"html"
	"strings"
	"time"
)

type Language string

const (
	EN Language = "en"
	ID Language = "id"
)

type User struct {
	ID             uint       `gorm:"primary_key" json:"id"`
	RoleID         uint       `json:"role_id,omitempty"`
	Role           Role       `gorm:"foreignKey:RoleID" json:"-"`
	Username       string     `gorm:"size:255;not null;unique" json:"username"`
	Password       string     `gorm:"size:255;not null;" json:"password,omitempty"`
	Email          string     `gorm:"email;not null;unique" json:"email"`
	Avatar         string     `json:"avatar"`
	PhoneNumber    string     `gorm:"not null; unique" json:"phone_number"`
	Lang           Language   `json:"language"`
	RecommendClass bool       `gorm:"default:true" json:"recommend_class"`
	Promotion      bool       `gorm:"default:true" json:"promotion"`
	Notification   bool       `gorm:"default:true" json:"notification"`
	LatestNews     bool       `gorm:"default:true" json:"latest_news"`
	ResetToken     string     `json:"reset_token,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `sql:"index" json:"deleted_at"`
}

func init() {
	database.Connect()
	dbConnection := database.GetConnection()
	dbConnection.AutoMigrate(&User{})
}

func (u *User) Create() (*User, error) {

	err := database.GetConnection().Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) BeforeSave() error {

	//turn password into hash
	hashedPassword, err := bcryptUtils.HashPassword(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)

	//remove spaces in username
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))

	return nil

}

func Find(uid uint) (User, error) {

	var u User

	if err := database.GetConnection().First(&u, uid).Error; err != nil {
		return u, errors.New("User not found!")
	}

	u.PrepareGive()

	return u, nil

}

func FindByEmail(email string) (User, error) {

	var user User
	if err := database.GetConnection().Where("email = ?", email).First(&user); errors.Is(err.Error, gorm.ErrRecordNotFound) {
		return user, err.Error
	}

	user.PrepareGive()

	return user, nil
}

func (user *User) UpdatePassword(password string) error {

	hashedPassword, err := bcryptUtils.HashPassword(password)
	if err != nil {
		return err
	}

	database.GetConnection().Model(&user).Updates(User{Password: string(hashedPassword)})
	return nil

}

func (user *User) IssueResetToken() {

	resetToken, _ := tokenUtils.GenerateResetToken(user.ID)
	database.GetConnection().Model(&user).Updates(User{ResetToken: resetToken})
	fmt.Println(resetToken)
}

func (user *User) RemoveResetToken() {

	database.GetConnection().Model(&user).Update("reset_token", nil)

}

func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func Authenticate(username string, password string) (interface{}, error) {

	u := User{}

	err := database.GetConnection().Model(User{}).Where("username = ?", username).Take(&u).Error

	if err != nil {
		return "", err
	}

	err = VerifyPassword(password, u.Password)

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}

	token, err := tokenUtils.GenerateAccessToken(u.ID)
	refreshToken, err := tokenUtils.GenerateRefershToken(u.ID)

	if err != nil {
		return "", err
	}

	pairToken := map[string]interface{}{
		"token":         token,
		"refresh_token": refreshToken,
	}

	return pairToken, nil

}

func (u *User) PrepareGive() {
	u.Password = ""
}
