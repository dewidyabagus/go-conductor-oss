package main

import (
	"time"

	"github.com/dewidyabagus/go-payout-workflow/sources/constant"
	"github.com/dewidyabagus/go-payout-workflow/sources/pkg/workflow"
)

func workerDefinitions(handler *handler) []workflow.WorkerDefinition {
	return []workflow.WorkerDefinition{
		{TaskName: constant.TaskHarsyaPaymentGateway, Handler: handler.DummyCompletedTask, BatchSize: 1, PollInterval: time.Second},
		{TaskName: constant.TaskPpobPrepaidPayment, Handler: handler.DummyCompletedTask, BatchSize: 1, PollInterval: time.Second},
		{TaskName: constant.TaskCreateTransaction, Handler: handler.CreateTransactionHandler, BatchSize: 1, PollInterval: time.Second},
		{TaskName: constant.TaskCreateInventory, Handler: handler.CreateInventoryHandler, BatchSize: 1, PollInterval: time.Second},
		{TaskName: constant.TaskCreateLedger, Handler: handler.CreateLedgerHandler, BatchSize: 1, PollInterval: time.Second},
		{TaskName: constant.TaskSuccessNotify, Handler: handler.SuccessNotificationHandler, BatchSize: 1, PollInterval: time.Second},
		{TaskName: constant.TaskDeleteTransaction, Handler: handler.DummyCompletedTask, BatchSize: 1, PollInterval: time.Second},
		{TaskName: constant.TaskDeleteInventory, Handler: handler.DummyCompletedTask, BatchSize: 1, PollInterval: time.Second},
		{TaskName: constant.TaskFailureNotify, Handler: handler.DummyCompletedTask, BatchSize: 1, PollInterval: time.Second},
		{TaskName: constant.TaskCreateCallbackLog, Handler: handler.DummyCompletedTask, BatchSize: 24, PollInterval: time.Second},
		{TaskName: constant.TaskSendMerchantCallback, Handler: handler.DummyCompletedTask, BatchSize: 24, PollInterval: time.Second},
		{TaskName: constant.TaskEmailAlert, Handler: handler.DummyCompletedTask, BatchSize: 1, PollInterval: time.Second},
	}
}
