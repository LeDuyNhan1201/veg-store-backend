package model

import (
	"time"
)

type Task struct {
	AuditingModel
	ID        UUID       `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Title     string     `gorm:"type:varchar(255);not null;uniqueIndex" json:"title"`
	StatusID  UUID       `gorm:"type:uuid;not null" json:"statusId"`                        // FK column
	Status    TaskStatus `gorm:"foreignKey:StatusID;references:ID" json:"status,omitempty"` // relation
	StartDay  *time.Time `gorm:"type:date" json:"startDay,omitempty"`
	TargetDay *time.Time `gorm:"type:date" json:"targetDay,omitempty"`
	EndDay    *time.Time `gorm:"type:date" json:"endDay,omitempty"`
	Priority  int8       `gorm:"type:smallint;not null;default:0" json:"priority"`
}

func (Task) TableName() string {
	return "task_tbl"
}

type TaskField string

const (
	FieldTitle     TaskField = "title"
	FieldStatusID  TaskField = "status_id"
	FieldStartDay  TaskField = "start_day"
	FieldTargetDay TaskField = "target_day"
	FieldEndDay    TaskField = "end_day"
	FieldPriority  TaskField = "priority"
)
