package model

type User struct {
	AuditingModel[string]        // Embed generic auditing model with string ID
	Name                  string `gorm:"type:varchar(100);not null"`
	Age                   int    `gorm:"not null"`
	Sex                   bool   `gorm:"not null"`
}

func (User) TableName() string {
	return "user_tbl"
}
