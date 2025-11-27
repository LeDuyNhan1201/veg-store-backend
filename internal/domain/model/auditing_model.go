package model

import (
	"time"

	"github.com/google/uuid"
)

// AllowedId Define a constraint for allowed ID types
type AllowedId interface {
	string | UUID | int
}

// UUID Define custom UUID type for GORM
type UUID string

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

type AuditingModel struct {
	CreatedAt time.Time `gorm:"autoCreateTime;not null"`
	CreatedBy string    `gorm:"type:varchar(100);not null"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt,omitempty"`
	UpdatedBy *string   `gorm:"type:varchar(100)" json:"UpdatedBy,omitempty"`
	IsDeleted bool      `gorm:"default:false;not null;index"`
	Version   int64     `gorm:"not null"`
}

func (m *AuditingModel) Created() *AuditingModel {
	m.IsDeleted = false
	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now()
	}
	if m.CreatedBy == "" {
		m.CreatedBy = "system"
	}
	m.Version = 1

	return m
}

func (m *AuditingModel) Updated() *AuditingModel {
	m.UpdatedAt = time.Now()
	m.Version++
	return m
}

func (m *AuditingModel) Deleted() *AuditingModel {
	m.UpdatedAt = time.Now()
	m.IsDeleted = true
	m.Version++
	return m
}
