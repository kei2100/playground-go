package response

import (
	"encoding/json"
	"log"
	"net/http"
)

// SendJSON sends an application/json message
func SendJSON(w http.ResponseWriter, status int, body interface{}) {
	if body == nil {
		body = struct{}{}
	}
	raw, err := json.Marshal(body)
	if err != nil {
		log.Printf("app: failed to marshal json %v : %v", body, err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json charset=utf-8")
	w.WriteHeader(status)
	w.Write(raw)
}

// errorResponse is a type of the error response
type errorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// apply error options
func (r *errorResponse) apply(opts []ErrorOpts) {
	for _, o := range opts {
		o(r)
	}
}

// ErrorOpts is a optional function for the errorResponse
type ErrorOpts func(response *errorResponse)

// WithMessage option
func WithMessage(message string) ErrorOpts {
	return func(res *errorResponse) {
		res.Message = message
	}
}

// WithErrorCode option
func WithErrorCode(code string) ErrorOpts {
	return func(res *errorResponse) {
		res.Code = code
	}
}

// kind of error code
const (
	CodeBadRequest          = "40000"
	CodeNotFound            = "40400"
	CodePayloadTooLarge     = "41300"
	CodeConflict            = "40900"
	CodeInternalServerError = "50000"
)

// SendBadRequest sends 400 Bad request
func SendBadRequest(w http.ResponseWriter, opts ...ErrorOpts) {
	res := errorResponse{
		Code:    CodeBadRequest,
		Message: "Bad request",
	}
	res.apply(opts)
	SendJSON(w, 400, &res)
}

// SendNotFound sends 404 Not found
func SendNotFound(w http.ResponseWriter, opts ...ErrorOpts) {
	res := errorResponse{
		Code:    CodeNotFound,
		Message: "Not found",
	}
	res.apply(opts)
	SendJSON(w, 404, &res)
}

// SendPayloadTooLarge sends 413 Payload too large
func SendPayloadTooLarge(w http.ResponseWriter, opts ...ErrorOpts) {
	res := errorResponse{
		Code:    CodePayloadTooLarge,
		Message: "Payload too large",
	}
	res.apply(opts)
	SendJSON(w, 413, &res)
}

// SendConflict sends 409 conflict
func SendConflict(w http.ResponseWriter, opts ...ErrorOpts) {
	res := errorResponse{
		Code:    CodeConflict,
		Message: "Conflict",
	}
	res.apply(opts)
	SendJSON(w, 409, &res)
}

// SendInternalServerError sends 500 Internal server error
func SendInternalServerError(w http.ResponseWriter, opts ...ErrorOpts) {
	res := errorResponse{
		Code:    CodeInternalServerError,
		Message: "Internal server error",
	}
	res.apply(opts)
	SendJSON(w, 500, &res)
}
