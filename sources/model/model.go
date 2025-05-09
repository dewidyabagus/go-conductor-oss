package model

type Response struct {
	Status string         `json:"status"`
	Error  *ErrorResponse `json:"error,omitempty"`
	Data   any            `json:"data"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	Details any    `json:"details"`
}
