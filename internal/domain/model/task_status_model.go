package model

type TaskStatus struct {
	AuditingModel
	ID    UUID   `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Title string `gorm:"type:varchar(100);not null" json:"title"`
}

func (TaskStatus) TableName() string {
	return "task_status_tbl"
}
