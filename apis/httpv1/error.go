package httpv1

type ErrorResponse struct {
	Message   string `example:"message"                              json:"message"`
	Code      string `example:"Conflict"                             json:"code"`
	RequestID string `example:"b5953bf0-9f15-4c42-afb4-1c125b40d7ce" json:"requestId"`
}
