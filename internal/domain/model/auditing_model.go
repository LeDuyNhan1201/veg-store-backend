package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// CompositeKey Define a composite key type
type CompositeKey struct {
	Key1 string `gorm:"type:uuid;primaryKey"`
	Key2 string `gorm:"type:uuid;primaryKey"`
}

// UUID Define custom UUID type for GORM
type UUID string

func (UUID) GormDataType() string                               { return "uuid" }
func (UUID) GormDBDataType(db *gorm.DB, _ *schema.Field) string { return "uuid" }
func (u UUID) String() string {
	return string(u)
}

func NewUUID() UUID {
	return UUID(uuid.New().String())
}

func ToUUID(str string) UUID {
	_, err := uuid.Parse(str)
	if err != nil {
		panic("Invalid UUID format: " + str)
	}
	return UUID(str)
}

// AllowedId Define a constraint for allowed ID types
type AllowedId interface {
	~string | CompositeKey
}

type AuditingModel[TId AllowedId] struct {
	Id        TId       `gorm:"primaryKey"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	CreatedBy string    `gorm:"type:varchar(100);not null"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	UpdatedBy string    `gorm:"type:varchar(100);not null"`
	IsDeleted bool      `gorm:"default:false"`
	Version   int64     `gorm:"version"`
}

func (m *AuditingModel[TId]) GetId() TId {
	return m.Id
}

func (m *AuditingModel[TId]) Created() *AuditingModel[TId] {
	m.IsDeleted = false
	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now()
	}
	if m.CreatedBy == "" {
		m.CreatedBy = "system"
	}
	m.Version = 1

	// Auto generate ID only for UUID
	switch id := any(m.Id).(type) {

	case UUID:
		// Convert ID to underlying string before comparison
		if string(id) == "" {
			newId := UUID(uuid.New().String())
			m.Id = any(newId).(TId)
		}

	case CompositeKey:
		// Composite key → do nothing (must be set externally)

	default:
		panic("Unsupported TId type for AuditingModel")
	}

	return m
}

func (m *AuditingModel[TId]) Updated() *AuditingModel[TId] {
	m.UpdatedAt = time.Now()
	m.Version++
	return m
}

func (m *AuditingModel[TId]) Deleted() *AuditingModel[TId] {
	m.UpdatedAt = time.Now()
	m.IsDeleted = true
	m.Version++
	return m
}
