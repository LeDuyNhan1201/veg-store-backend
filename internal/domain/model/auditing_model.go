package model

import "time"

type IAuditingModel[TId comparable] interface {
	GetId() TId
	Created() IAuditingModel[TId]
	Updated() IAuditingModel[TId]
	Deleted() IAuditingModel[TId]
}

type AuditingModel[TId comparable] struct {
	Id        TId       `gorm:"primaryKey"`       // Primary key
	CreatedAt time.Time `gorm:"auto_create_time"` // Auto-set on insert
	UpdatedAt time.Time `gorm:"auto_update_time"` // Auto-set on update
	IsDeleted bool      `gorm:"default:false"`    // Soft delete flag
	Version   int64     `gorm:"version"`          // Optimistic lock
}

func (model *AuditingModel[TId]) GetId() TId {
	return model.Id
}

func (model *AuditingModel[TId]) Created() IAuditingModel[TId] {
	model.IsDeleted = false
	model.CreatedAt = time.Now()
	model.UpdatedAt = time.Now()
	model.Version = 1
	return model
}

func (model *AuditingModel[TId]) Updated() IAuditingModel[TId] {
	model.UpdatedAt = time.Now()
	model.Version = 2
	return model
}

func (model *AuditingModel[TId]) Deleted() IAuditingModel[TId] {
	model.UpdatedAt = time.Now()
	model.IsDeleted = true
	model.Version = 3
	return model
}
