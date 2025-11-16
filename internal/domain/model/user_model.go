package model

type User struct {
	AuditingModel[UUID]        // Embed generic auditing model with string ID
	Name                string `gorm:"type:varchar(100);not null"`
	Age                 int8   `gorm:"not null"`
	Sex                 bool   `gorm:"not null"`
	Email               string `gorm:"type:varchar(150);not null;uniqueIndex"`
	Password            string `gorm:"type:text;not null"`
}

func (User) TableName() string {
	return "user_tbl"
}
