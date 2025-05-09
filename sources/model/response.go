package model

type PrepaidPaymentResponse struct {
	ReferenceId string `json:"referenceId"`
	PrepaidPaymentRequest
	Amount          float64 `json:"amount"`
	BankReferenceNo string  `json:"bankReferenceNo,omitempty"`
	Status          string  `json:"status"`
}

type HarsyaPaymentNotificationOutput struct {
	ReferenceNo     string `json:"referenceNo"`
	BankReferenceNo string `json:"bankReferenceNo"`
}
