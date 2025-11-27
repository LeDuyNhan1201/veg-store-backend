package model

type User struct {
	AuditingModel
	ID       UUID   `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name     string `gorm:"type:varchar(100);not null"`
	Age      int8   `gorm:"not null"`
	Sex      bool   `gorm:"not null"`
	Email    string `gorm:"type:varchar(150);not null;uniqueIndex"`
	Password string `gorm:"type:text;not null"`
}

func (User) TableName() string {
	return "user_tbl"
}
