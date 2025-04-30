package entity

type Task struct {
	ID          int    `gorm:"primaryKey" json:"id"`
	Description string `json:"description"`
	UserID      int    `gorm:"not null"`
	User        User   `gorm:"foreignKey:UserID;references:ID"` // Foreign key relationship
}

type TaskRequest struct {
	Description string `json:"description"`
}
