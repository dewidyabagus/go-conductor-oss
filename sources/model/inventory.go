package model

type TransactionRequest struct {
	OrderId string                   `json:"orderId"`
	Items   []map[string]interface{} `json:"items"`
}

type TransactionResponse struct {
	Id   string `json:"id"`
	Date string `json:"date"`
}

type InventoryRequest struct {
	TransactionRequest
	TransactionId   string `json:"transactionId"`
	TransactionDate string `json:"transactionDate"`
}

type LedgerRequest InventoryRequest
