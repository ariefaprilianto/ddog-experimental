package err

import (
	"encoding/json"
	"net/http"
)

//
// An ErrorMessage represents an error message format and list of error.
//
type ErrorMessage struct {
	ErrorList []*ErrorFormat `json:"error_list"`
	Code      int            `json:"code"`
}

//
// An ErrorFormat represents an error message format and code that we used.
//
type ErrorFormat struct {
	ErrorName        string `json:"error_name"`
	ErrorDescription string `json:"error_description"`
}

//
// Common internal server error message
//
const (
	InternalServerName        = "internal_server_error"
	InternalServerDescription = "The server is unable to complete your request"
)

//
// Common error unautorized
//
const (
	UnauthorizedErrorName        = "access_denied"
	UnauthorizedErrorDescription = "Authorization failed by filter."
)

//
// Common bad request error message
//
var DefaultBadRequest = ErrorFormat{
	ErrorName:        "bad_request",
	ErrorDescription: "Your request resulted in error",
}

//
// Create new error message
//
func NewErrorMessage() *ErrorMessage {
	return &ErrorMessage{}
}

//
// Set bad request
//
func (em *ErrorMessage) SetBadRequest() *ErrorMessage {
	em.Code = http.StatusBadRequest
	return em
}

//
// SetNewError is function return new error message.
// It support to set code, error name, and error description
//
func SetNewError(code int, errorName, errDesc string) *ErrorMessage {
	return &ErrorMessage{
		Code: code,
		ErrorList: []*ErrorFormat{
			&ErrorFormat{
				ErrorName:        errorName,
				ErrorDescription: errDesc,
			},
		},
	}
}

//
// SetNewBadRequest is function return new error message with bad request standard code(400).
// It support to set error name and error description
//
func SetNewBadRequest(errorName, errDesc string) *ErrorMessage {
	return SetNewError(http.StatusBadRequest, errorName, errDesc)
}

//
// SetNewBadRequest is function return new error message with bad request standard code(400).
// It support to set error name and error description using error format
//
func SetNewBadRequestByFormat(ef *ErrorFormat) *ErrorMessage {
	return &ErrorMessage{
		Code: http.StatusBadRequest,
		ErrorList: []*ErrorFormat{
			ef,
		},
	}
}

func SetDefaultNewBadRequest() *ErrorMessage {
	return SetNewBadRequestByFormat(&DefaultBadRequest)
}

//
// SetNewInternalError is function return new error message with internal server error standard code(500).
//
func SetNewInternalError() *ErrorMessage {
	return SetNewError(http.StatusInternalServerError, InternalServerName, InternalServerDescription)
}

//
// SetNewUnauthorizedError is function return new error message with unauthorized error code(400).
// It support to set error name and error description
//
func SetNewUnauthorizedError(errorName, errDesc string) *ErrorMessage {
	return SetNewError(http.StatusUnauthorized, errorName, errDesc)
}

//
// SetNewUnauthorizedError is function return new error message with unauthorized error code(400).
// It support to set error name and error description
//
func SetDefaultUnauthorized() *ErrorMessage {
	return SetNewUnauthorizedError(UnauthorizedErrorName, UnauthorizedErrorDescription)
}

//
// Append is function add error to existing error message.
// It support to set error name and error description.
//
func (errorMessage *ErrorMessage) Append(errorName, errDesc string) *ErrorMessage {
	errorMessage.ErrorList = append(errorMessage.ErrorList, &ErrorFormat{
		ErrorName:        errorName,
		ErrorDescription: errDesc,
	})
	return errorMessage
}

//
// AppendFormat is function add error to existing error message.
// It support to set error name and error description using error format
//
func (errorMessage *ErrorMessage) AppendFormat(ef *ErrorFormat) *ErrorMessage {
	errorMessage.ErrorList = append(errorMessage.ErrorList, ef)
	return errorMessage
}

//
// GetListError is function to get list error message.
//
func (errorMessage *ErrorMessage) GetListError() []*ErrorFormat {
	return errorMessage.ErrorList
}

//
// GetCode is function to get code.
//
func (errorMessage *ErrorMessage) GetCode() int {
	return errorMessage.Code
}

//
// Get error byte
//
func (errorMessage *ErrorMessage) Marshal() []byte {
	b, _ := json.Marshal(errorMessage)
	return b
}

//
// Get string
//
func (errorMessage *ErrorMessage) ToString() string {
	return string(errorMessage.Marshal())
}

// Error to implement error interface
func (errorMessage *ErrorMessage) Error() string {
	return errorMessage.ToString()
}
