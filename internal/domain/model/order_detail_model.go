package model

type OrderDetail struct {
	AuditingModel
	OrderID   UUID    `gorm:"primaryKey;column:order_id"`
	ProductID UUID    `gorm:"primaryKey;column:product_id"`
	Order     Order   `gorm:"foreignKey:OrderID;references:ID" json:"order"`
	Product   Product `gorm:"foreignKey:ProductID;references:ID" json:"product"`

	// TODO: Add other fields
}

func (OrderDetail) TableName() string {
	return "order_detail_tbl"
}
