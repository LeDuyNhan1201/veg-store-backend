package model

type Order struct {
	AuditingModel
	ID UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`

	// TODO: Add other fields
}

func (Order) TableName() string {
	return "order_tbl"
}
