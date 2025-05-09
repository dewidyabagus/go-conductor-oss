package constant

const (
	// Workflow
	WorkflowPrepaidPaymentHttpTrigger = "prepaid_payment_http_trigger"

	// Task
	TaskHarsyaPaymentGateway      = "harsya_payment_gateway_task"
	TaskHarsyaPaymentNotification = "harsya_payment_notification"
	TaskPpobPrepaidPayment        = "ppob_prepaid_payment_task"
	TaskCreateTransaction         = "create_transaction"
	TaskCreateInventory           = "create_inventory"
	TaskCreateLedger              = "create_ledger"
	TaskSuccessNotify             = "success_notification"
	TaskFailureNotify             = "failure_notification_task"
	TaskDeleteInventory           = "delete_inventory_task"
	TaskDeleteTransaction         = "delete_transaction_task"

	// Task Domain
	TaskDomainOpenAPI = "OpenApi"
)

const (
	HeaderXMerchantId = "X-Merchant-Id"
)

const (
	MerchantInsufficientBalance = "INSUFFICIENT_BALANCE"
	MerchantSufficientBalance   = "SUFFICIENT_BALANCE"
)

const (
	StatusPending = "PENDING"
	StatusSuccess = "SUCCESS"
	StatusFailed  = "FAILED"
)
