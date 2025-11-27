package model

type Product struct {
	AuditingModel
	ID UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`

	// TODO: Add other fields
}

func (Product) TableName() string {
	return "product_tbl"
}
