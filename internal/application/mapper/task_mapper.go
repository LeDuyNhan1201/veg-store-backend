package mapper

import (
	"veg-store-backend/internal/application/dto"
	"veg-store-backend/internal/domain/model"
	"veg-store-backend/util"
)

func ToTaskItem(task *model.Task) dto.TaskItem {
	var status dto.PreviewStatus
	if task.Status.ID == "" {
		status = dto.PreviewStatus{
			ID: task.StatusID.String(),
		}
	} else {
		status = ToPreviewTaskStatus(&task.Status)
	}

	return dto.TaskItem{
		ID:        task.ID.String(),
		Title:     task.Title,
		StartDay:  task.StartDay,
		EndDay:    task.EndDay,
		Status:    status,
		TargetDay: task.TargetDay,
		Priority:  task.Priority,
	}
}

func ToPreviewTaskStatus(status *model.TaskStatus) dto.PreviewStatus {
	return dto.PreviewStatus{
		ID:    status.ID.String(),
		Title: status.Title,
	}
}

func ToTask(request dto.CreateTaskRequest) model.Task {
	startDay := util.ParseDay(request.StartDay)
	endDay := util.ParseDay(request.EndDay)
	targetDay := util.ParseDay(request.TargetDay)

	return model.Task{
		Title:     request.Title,
		StatusID:  model.ToUUID(request.StatusID),
		StartDay:  &startDay,
		EndDay:    &endDay,
		TargetDay: &targetDay,
	}
}
