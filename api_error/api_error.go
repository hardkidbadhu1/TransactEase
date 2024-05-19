package api_error

type ErrorCode string

type ErrorResponse struct {
	HttpStatusCode int       `json:"-"`
	ErrorCode      ErrorCode `json:"error_code"`
	ErrorMessage   string    `json:"error_message"`
}

const (
	InvalidRequest      ErrorCode = "INVALID_REQUEST"
	InternalServerError ErrorCode = "INTERNAL_SERVER_ERROR"
)

func NewErrorResponse(httpStatusCode int, errorCode ErrorCode, errorMessage string) ErrorResponse {
	return ErrorResponse{
		HttpStatusCode: httpStatusCode,
		ErrorCode:      errorCode,
		ErrorMessage:   errorMessage,
	}
}

func NewInternalServerError(errorMessage string) ErrorResponse {
	return NewErrorResponse(500, InternalServerError, errorMessage)
}

var (
	InvalidDocumentNumber = NewErrorResponse(400, InvalidRequest, "document number is required")
	InvalidParams         = NewErrorResponse(400, InvalidRequest, "invalid params")
)
