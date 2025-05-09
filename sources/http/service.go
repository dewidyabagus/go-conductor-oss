package main

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/dewidyabagus/go-payout-workflow/sources/constant"
	"github.com/dewidyabagus/go-payout-workflow/sources/data"
	"github.com/dewidyabagus/go-payout-workflow/sources/model"
	"github.com/dewidyabagus/go-payout-workflow/sources/pkg/utils"
	"github.com/dewidyabagus/go-payout-workflow/sources/pkg/workflow"

	workflowModel "github.com/conductor-sdk/conductor-go/sdk/model"
	"github.com/google/uuid"
)

type service struct {
	workflow *workflow.WorkflowExecutor
}

func (s *service) HelloWorld(ctx context.Context, message string) (result any, err error) {
	workflowId, err := s.workflow.StartWorkflowWithContext(ctx, &workflowModel.StartWorkflowRequest{
		Name:    "hello_world_workflow",
		Version: 1,
		Input: map[string]string{
			"message": message,
		},
		Priority:      1,
		CorrelationId: uuid.NewString(),
	})
	if err != nil {
		fmt.Println("Error:", err.Error())
		return nil, utils.WrapError(utils.ErrInternal, err)
	}

	channel, err := s.workflow.MonitorExecution(workflowId)
	if err != nil {
		return nil, utils.WrapError(utils.ErrInternal, err)
	}
	run := <-channel

	return run.Output, nil
}

func (s *service) PrepaidPayment(ctx context.Context, request model.PrepaidPaymentRequest) (result model.PrepaidPaymentResponse, err error) {
	product, ok := data.GetPrepaidProductById(request.ProductId)
	if !ok {
		return result, utils.WrapError(utils.ErrDataNotFound, errors.New("product not foundß"))
	}

	id, _ := uuid.NewV7()
	referenceId := utils.GenerateULID()

	data.Transactions = append(data.Transactions, data.Transaction{
		Id:          id.String(),
		ReferenceId: referenceId,
		Status:      "PENDING",
	})
	idx := len(data.Transactions) - 1

	result = model.PrepaidPaymentResponse{
		ReferenceId:           referenceId,
		PrepaidPaymentRequest: request,
		Amount:                product.Amount + product.Fee,
	}
	defer func() {
		result.Status = data.Transactions[idx].Status
	}()

	workflowRequest := &workflowModel.StartWorkflowRequest{
		Name:     constant.WorkflowPrepaidPaymentHttpTrigger,
		Version:  1,
		Priority: 1,
		Input: map[string]any{
			"referenceId":          referenceId,
			"customerId":           request.CustomerId,
			"productId":            request.ProductId,
			"amount":               product.Amount + product.Fee,
			"paymentMethod":        request.PaymentMethod,
			"paymentMethodOptions": request.PaymentMethodOptions,
		},
	}
	workflowId, err := s.workflow.StartWorkflowWithContext(ctx, workflowRequest)
	if err != nil {
		fmt.Println("Start Workflow Error:", err.Error())

		data.Transactions[idx].Status = "FAILED" // Ketika dalam session transaction db bisa di rollback
		return result, utils.WrapError(utils.ErrInternal, err)
	}
	data.Transactions[idx].WorkflowId = workflowId

	workflowChannel, err := s.workflow.MonitorExecution(workflowId)
	if err != nil {
		fmt.Println("Monitor Execution Error:", err.Error())
		return result, nil
	}

	timeout := time.After(60 * time.Second)

	select {
	case <-timeout:
		// No action, transaction status is PENDING

	case workflowResult := <-workflowChannel:
		result.BankReferenceNo, _ = workflowResult.Output["bankReferenceNo"].(string)
		switch workflowResult.Status {
		case workflowModel.CompletedWorkflow:
			data.Transactions[idx].Status = constant.StatusSuccess

		case workflowModel.FailedWorkflow, workflowModel.TerminatedWorkflow, workflowModel.TimedOutWorkflow:
			data.Transactions[idx].Status = constant.StatusFailed

		case workflowModel.PausedWorkflow, workflowModel.RunningWorkflow:
			data.Transactions[idx].Status = constant.StatusPending
		}
	}
	return
}

func (s *service) HarsyaPaymentNotification(ctx context.Context, request model.HarsyaPaymentNotificationRequest) error {
	idx := slices.IndexFunc(data.Transactions, func(trx data.Transaction) bool {
		return trx.ReferenceId == request.ReferenceNo
	})
	if idx < 0 {
		return utils.WrapError(utils.ErrDataNotFound, errors.New("transaction not found"))
	}
	status := workflowModel.CompletedTask
	switch request.Status {
	case constant.StatusPending:
		status = workflowModel.InProgressTask

	case constant.StatusFailed:
		status = workflowModel.FailedTask
	}
	return s.workflow.UpdateTaskByRefNameWithContext(
		ctx, constant.TaskHarsyaPaymentNotification, data.Transactions[idx].WorkflowId, status, model.HarsyaPaymentNotificationOutput{
			ReferenceNo:     request.ReferenceNo,
			BankReferenceNo: request.BankReferenceNo,
		},
	)
}
