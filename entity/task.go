package entity

type Task struct {
	ID          int    `gorm:"primaryKey" json:"id"`
	Description string `json:"description" binding:"required"`
	UserID      int    `gorm:"not null"`
}

type TaskRequest struct {
	Description string `json:"description" binding:"required"`
}
