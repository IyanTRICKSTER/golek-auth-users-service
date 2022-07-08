package requests

type Language string

const (
	EN Language = "en"
	ID Language = "id"
)

type UpdateUserRecordCredential struct {
	Username       string   `json:"username" binding:"required"`
	Email          string   `json:"email" binding:"required,email"`
	PhoneNumber    string   `json:"phone_number" binding:"required"`
	Avatar         string   `json:"avatar" binding:"required"`
	Lang           Language `json:"language" binding:"required"`
	RecommendClass bool     `json:"recommend_class" binding:"required"`
	Promotion      bool     `json:"promotion" binding:"required"`
	Notification   bool     `json:"notification" binding:"required"`
	LatestNews     bool     `json:"latest_news" binding:"required"`
}
