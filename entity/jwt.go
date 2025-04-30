package entity

type RefreshRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}
