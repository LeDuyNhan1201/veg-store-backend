package dto

import (
	"time"
)

type TaskItem struct {
	ID        string        `json:"id"`
	Title     string        `json:"title"`
	Status    PreviewStatus `json:"status"`
	StartDay  *time.Time    `json:"startDay,omitempty"`
	TargetDay *time.Time    `json:"targetDay,omitempty"`
	EndDay    *time.Time    `json:"endDay,omitempty"`
	Priority  int8          `json:"priority"`
}

type PreviewStatus struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type CreateTaskRequest struct {
	Title     string `json:"title" binding:"required,range=3-200" example:"Create Task"`
	StatusID  string `json:"statusId" binding:"required,uuid" example:"5259ac80-1823-44d1-a701-0ed1e36fb38c"`
	StartDay  string `json:"startDay" binding:"required" example:"2001-12-31"`
	TargetDay string `json:"targetDay" binding:"required" example:"2001-12-31"`
	EndDay    string `json:"endDay" binding:"required" example:"2001-12-31"`
}

type UpdateTaskStatusRequest struct {
	ID       string `json:"id" binding:"required,uuid" example:"5259ac80-1823-44d1-a701-0ed1e36fb38c"`
	StatusID string `json:"statusId" binding:"required,uuid" example:"5259ac80-1823-44d1-a701-0ed1e36fb38c"`
}

type SortTasksOption struct {
	Field     string    `json:"field" binding:"taskFields"`
	Direction Direction `json:"direction"`
}

type AdvancedFilterTaskRequest struct {
	Keyword  string            `json:"keyword" example:"Task"`
	FromDate string            `json:"fromDate" example:"2001-12-31"`
	ToDate   string            `json:"toDate" example:"2001-12-31"`
	Sorts    []SortTasksOption `json:"sorts" binding:"dive"` // dive to validate each element of slice
}
