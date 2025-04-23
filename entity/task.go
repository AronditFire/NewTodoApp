package entity

type Task struct {
	ID          int    `gorm:"primaryKey" json:"id"`
	Description string `json:"description"`
}

type TaskRequest struct {
	Description string `json:"description"`
}
