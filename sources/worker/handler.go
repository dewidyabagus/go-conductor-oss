package main

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	conductor "github.com/conductor-sdk/conductor-go/sdk/model"
	"github.com/dewidyabagus/go-payout-workflow/sources/model"
	"github.com/google/uuid"
)

type handler struct {
	service *service
}

func (handler) DummyCompletedTask(_ context.Context, task *conductor.Task) (*conductor.TaskResult, error) {
	taskResult := conductor.NewTaskResultFromTask(task)
	taskResult.Status = conductor.CompletedTask

	return taskResult, nil
}

func (handler) DummyFailedTask(_ context.Context, task *conductor.Task) (*conductor.TaskResult, error) {
	return conductor.NewTaskResultFromTaskWithError(task, errors.New("test error")), nil
}

func (handler) CreateTransactionHandler(_ context.Context, t *conductor.Task) (*conductor.TaskResult, error) {
	result := conductor.NewTaskResultFromTask(t)

	raw, err := json.Marshal(t.InputData)
	if err != nil {
		result.ReasonForIncompletion = err.Error()
		result.Status = conductor.FailedWithTerminalErrorTask

		return result, err
	}

	request := model.TransactionRequest{}
	if err = json.Unmarshal(raw, &request); err != nil {
		result.ReasonForIncompletion = err.Error()
		result.Status = conductor.FailedWithTerminalErrorTask

		return result, err
	}
	result.Status = conductor.CompletedTask

	if request.OrderId == "000002" {
		result.Status = conductor.FailedWithTerminalErrorTask
		result.ReasonForIncompletion = "Failed when create transaction"
		return result, err
	}

	result.OutputData = map[string]interface{}{
		"id":   uuid.NewString(),
		"date": time.Now().Format(time.DateTime),
	}
	return result, nil
}

func (handler) CreateInventoryHandler(_ context.Context, t *conductor.Task) (*conductor.TaskResult, error) {
	result := conductor.NewTaskResultFromTask(t)

	raw, err := json.Marshal(t.InputData)
	if err != nil {
		result.ReasonForIncompletion = err.Error()
		result.Status = conductor.FailedWithTerminalErrorTask

		return result, err
	}

	request := model.InventoryRequest{}
	if err = json.Unmarshal(raw, &request); err != nil {
		result.ReasonForIncompletion = err.Error()
		result.Status = conductor.FailedWithTerminalErrorTask

		return result, err
	}

	if request.OrderId == "000003" {
		result.Status = conductor.FailedWithTerminalErrorTask
		result.ReasonForIncompletion = "Failed when create inventory"
		return result, err
	}
	result.Status = conductor.CompletedTask

	result.OutputData = map[string]interface{}{
		"items": request.Items,
	}
	return result, nil
}

func (handler) CreateLedgerHandler(_ context.Context, t *conductor.Task) (*conductor.TaskResult, error) {
	result := conductor.NewTaskResultFromTask(t)

	raw, err := json.Marshal(t.InputData)
	if err != nil {
		result.ReasonForIncompletion = err.Error()
		result.Status = conductor.FailedWithTerminalErrorTask

		return result, err
	}

	request := model.LedgerRequest{}
	if err = json.Unmarshal(raw, &request); err != nil {
		result.ReasonForIncompletion = err.Error()
		result.Status = conductor.FailedWithTerminalErrorTask

		return result, err
	}

	if request.OrderId == "000004" {
		result.Status = conductor.FailedWithTerminalErrorTask
		result.ReasonForIncompletion = "Failed when create ledger"
		return result, err
	}
	result.Status = conductor.CompletedTask

	result.OutputData = map[string]interface{}{
		"journalNo": uuid.NewString(),
	}
	return result, nil
}

func (handler) SuccessNotificationHandler(_ context.Context, t *conductor.Task) (*conductor.TaskResult, error) {
	result := conductor.NewTaskResultFromTask(t)
	result.Status = conductor.CompletedTask

	return result, nil
}
