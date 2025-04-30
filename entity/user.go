package entity

type User struct {
	ID       int    `gorm:"primaryKey" json:"id"`
	Username string `gorm:"uniqueIndex" json:"username"`
	Password string `json:"-"`
	IsAdmin  bool   `gorm:"default:false" json:"-"`
}

type UserRegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"pass" binding:"required"`
}

type UserAuthRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"pass" binding:"required"`
}
