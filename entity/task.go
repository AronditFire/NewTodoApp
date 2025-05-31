package entity

type Task struct {
	ID          int    `gorm:"primaryKey" json:"id" redis:"id"`
	Description string `json:"description" binding:"required" redis:"description"`
	UserID      int    `gorm:"not null" redis:"user_id"`
}

type TaskRequest struct {
	Description string `json:"description" binding:"required"`
}
