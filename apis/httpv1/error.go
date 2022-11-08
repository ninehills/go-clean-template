package httpv1

type ErrorResponse struct {
	Message   string `json:"message" example:"message"`
	Code      string `json:"code" example:"Conflict"`
	RequestID string `json:"requestId" example:"b5953bf0-9f15-4c42-afb4-1c125b40d7ce"`
}
