package requests

type Language string

const (
	EN Language = "en"
	ID Language = "id"
)

type UpdateUserRecordCredential struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Avatar   string `json:"avatar" binding:"required"`
}
