package utils

import (
	conductor "github.com/conductor-sdk/conductor-go/sdk/model"
)

func TaskCompletedWithDataResponse(task *conductor.Task, data any) (*conductor.TaskResult, error) {
	taskResult := conductor.NewTaskResultFromTask(task)
	if err := CopyValueWithJSONTags(&taskResult.OutputData, data); err != nil {
		return TaskNonRetryableErrorResponse(task, err)
	}
	taskResult.Status = conductor.CompletedTask

	return taskResult, nil
}

func TaskCompletedResponse(task *conductor.Task) (*conductor.TaskResult, error) {
	taskResult := conductor.NewTaskResultFromTask(task)
	taskResult.Status = conductor.CompletedTask

	return taskResult, nil
}
