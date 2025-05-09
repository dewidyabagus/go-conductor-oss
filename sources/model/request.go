package model

type PrepaidPaymentRequest struct {
	CustomerId           string `json:"customerId"`
	ProductId            string `json:"productId"`
	PaymentMethod        string `json:"paymentMethod"`
	PaymentMethodOptions any    `json:"paymentMethodOptions"`
}

type HarsyaPaymentNotificationRequest struct {
	ReferenceNo     string `json:"referenceNo" validate:"required"`
	BankReferenceNo string `json:"bankReferenceNo" validate:"required"`
	Status          string `json:"status" validate:"required"`
}
