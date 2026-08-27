package response

type ErrorResponseFormat struct {
	Status   string   `json:"status"`
	Code     string   `json:"code"`
	Message  string   `json:"message,omitempty"`
	Data     any      `json:"data"`
	Metadata Metadata `json:"metadata"`
}

type ErrorAuthzResponseFormat struct {
	Status   string   `json:"status"`
	Code     string   `json:"code"`
	Message  any      `json:"message,omitempty"`
	Data     any      `json:"data"`
	Metadata Metadata `json:"metadata"`
}

type ErrorResponse struct {
	Code        int    `json:"code,omitempty"`
	Message     string `json:"message,omitempty"`
	Description string `json:"description,omitempty"`

	StackTrace string       `json:"stack_trace,omitempty"`
	FieldError []FieldError `json:"field_error,omitempty"`
	Data       interface{}  `json:"field_error,omitempty"`
}

type FieldError struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Metadata struct {
	ServerTimestamp string `json:"server_timestamp"`
}

type SuccessResponse struct {
	Status   string      `json:"status"`
	Code     string      `json:"code"`
	Message  any         `json:"message,omitempty"`
	Data     interface{} `json:"data"`
	Metadata Metadata    `json:"metadata"`
}

// package response

// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"strings"
// 	"time"

// 	"ecomerce/internal/errors"

// 	validation "github.com/go-playground/validator/v10"
// 	"github.com/joomcode/errorx"
// 	"github.com/spf13/viper"
// )

// func SendSuccessResponse(w http.ResponseWriter, statusCode int, code Code, data any, message any) {
// 	response := SuccessResponse{
// 		Status:  "success",
// 		Code:    string(code),
// 		Message: message,
// 		Data:    data,
// 		Metadata: Metadata{
// 			ServerTimestamp: time.Now().UTC().Format("2006-01-02T15:04:05Z"),
// 		},
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(statusCode)
// 	_ = json.NewEncoder(w).Encode(response)
// }

// func SendErrorResponse(w http.ResponseWriter, err *ErrorResponse) {
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(err.Code)
// 	_ = json.NewEncoder(w).Encode(err)
// }

// func SendErrorResponseFormated(w http.ResponseWriter, err error) {
// 	statusCode := http.StatusInternalServerError
// 	message := "Internal Server Error"
// 	code := "internal_error"

// 	if err != nil {
// 		if validationErrors, ok := err.(validation.ValidationErrors); ok {
// 			statusCode = http.StatusBadRequest
// 			message = "Validation error"
// 			code = "validation_error"

// 			for _, fieldError := range validationErrors {
// 				message = fieldError.Field() + ": " + fieldError.ActualTag()
// 			}
// 		} else {
// 			for _, e := range errors.Error {
// 				if errorx.IsOfType(err, e.Type) {
// 					statusCode = e.StatusCode
// 					message = strings.ReplaceAll(
// 						e.Type.FullName(),
// 						fmt.Sprintf("%s.", e.Type.Namespace().String()),
// 						"",
// 					)

// 					if strings.Contains(err.Error(), "attempt left") ||
// 						strings.Contains(err.Error(), "attempts left") {
// 						message = strings.ReplaceAll(
// 							err.Error(),
// 							"CUSTOMER_NOT_VERIFIED.invalid pin",
// 							"Invalid Pin",
// 						)
// 					}

// 					code = e.Type.Namespace().String()
// 					break
// 				}
// 			}

// 			if message == "Internal Server Error" {
// 				message = err.Error()
// 				code = "unknown_error"
// 			}
// 		}
// 	}

// 	response := ErrorResponseFormat{
// 		Status:  "error",
// 		Code:    code,
// 		Message: message,
// 		Metadata: Metadata{
// 			ServerTimestamp: time.Now().UTC().Format("2006-01-02T15:04:05Z"),
// 		},
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(statusCode)
// 	_ = json.NewEncoder(w).Encode(response)
// }

// func GetErrorFrom(err error) *ErrorResponse {
// 	debugMode := viper.GetBool("debug")

// 	for _, e := range errors.Error {
// 		if errorx.IsOfType(err, e.Type) {
// 			er := errorx.Cast(err)

// 			res := ErrorResponse{
// 				Code:       e.StatusCode,
// 				Message:    er.Message(),
// 				FieldError: ErrorFields(er.Cause()),
// 			}

// 			if debugMode {
// 				res.Description = fmt.Sprintf("Error: %v", er)
// 				res.StackTrace = fmt.Sprintf("%+v", errorx.EnsureStackTrace(err))
// 			}

// 			return &res
// 		}
// 	}

// 	return &ErrorResponse{
// 		Code:    http.StatusInternalServerError,
// 		Message: "Unknown server error",
// 	}
// }

// func ErrorFields(err error) []FieldError {
// 	var errs []FieldError

// 	if data, ok := err.(validation.ValidationErrors); ok {
// 		for i, v := range data {
// 			errs = append(errs, FieldError{
// 				Name:        fmt.Sprintf("%d", i+1),
// 				Description: v.Error(),
// 			})
// 		}

// 		return errs
// 	}

// 	return nil
// }

// func SendAuthzResponseErr(w http.ResponseWriter, err error, message any) {
// 	statusCode := http.StatusInternalServerError
// 	code := "internal_error"

// 	if err != nil {
// 		for _, e := range errors.Error {
// 			if errorx.IsOfType(err, e.Type) {
// 				statusCode = e.StatusCode
// 				code = e.Type.Namespace().String()
// 				break
// 			}
// 		}

// 		if message == "Internal Server Error" {
// 			message = err.Error()
// 			code = "unknown_error"
// 		}
// 	}

// 	response := ErrorAuthzResponseFormat{
// 		Status:  "error",
// 		Code:    code,
// 		Message: message,
// 		Metadata: Metadata{
// 			ServerTimestamp: time.Now().UTC().Format("2006-01-02T15:04:05Z"),
// 		},
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(statusCode)
// 	_ = json.NewEncoder(w).Encode(response)
// }

// With Chi, you use the standard `net/http` package instead of `gin.Context`. The biggest change is replacing:

// * `ctx.JSON(...)` → `json.NewEncoder(w).Encode(...)`
// * `ctx.AbortWithStatusJSON(...)` → write the status code and encode the JSON response.

// Here's the Chi version:

// ```go
// package response

// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"strings"
// 	"time"

// 	"ecomerce/internal/errors"

// 	validation "github.com/go-playground/validator/v10"
// 	"github.com/joomcode/errorx"
// 	"github.com/spf13/viper"
// )

// func SendSuccessResponse(w http.ResponseWriter, statusCode int, code Code, data any, message any) {
// 	response := SuccessResponse{
// 		Status:  "success",
// 		Code:    string(code),
// 		Message: message,
// 		Data:    data,
// 		Metadata: Metadata{
// 			ServerTimestamp: time.Now().UTC().Format("2006-01-02T15:04:05Z"),
// 		},
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(statusCode)
// 	_ = json.NewEncoder(w).Encode(response)
// }

// func SendErrorResponse(w http.ResponseWriter, err *ErrorResponse) {
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(err.Code)
// 	_ = json.NewEncoder(w).Encode(err)
// }

// func SendErrorResponseFormated(w http.ResponseWriter, err error) {
// 	statusCode := http.StatusInternalServerError
// 	message := "Internal Server Error"
// 	code := "internal_error"

// 	if err != nil {
// 		if validationErrors, ok := err.(validation.ValidationErrors); ok {
// 			statusCode = http.StatusBadRequest
// 			message = "Validation error"
// 			code = "validation_error"

// 			for _, fieldError := range validationErrors {
// 				message = fieldError.Field() + ": " + fieldError.ActualTag()
// 			}
// 		} else {
// 			for _, e := range errors.Error {
// 				if errorx.IsOfType(err, e.Type) {
// 					statusCode = e.StatusCode
// 					message = strings.ReplaceAll(
// 						e.Type.FullName(),
// 						fmt.Sprintf("%s.", e.Type.Namespace().String()),
// 						"",
// 					)

// 					if strings.Contains(err.Error(), "attempt left") ||
// 						strings.Contains(err.Error(), "attempts left") {
// 						message = strings.ReplaceAll(
// 							err.Error(),
// 							"CUSTOMER_NOT_VERIFIED.invalid pin",
// 							"Invalid Pin",
// 						)
// 					}

// 					code = e.Type.Namespace().String()
// 					break
// 				}
// 			}

// 			if message == "Internal Server Error" {
// 				message = err.Error()
// 				code = "unknown_error"
// 			}
// 		}
// 	}

// 	response := ErrorResponseFormat{
// 		Status:  "error",
// 		Code:    code,
// 		Message: message,
// 		Metadata: Metadata{
// 			ServerTimestamp: time.Now().UTC().Format("2006-01-02T15:04:05Z"),
// 		},
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(statusCode)
// 	_ = json.NewEncoder(w).Encode(response)
// }

// func GetErrorFrom(err error) *ErrorResponse {
// 	debugMode := viper.GetBool("debug")

// 	for _, e := range errors.Error {
// 		if errorx.IsOfType(err, e.Type) {
// 			er := errorx.Cast(err)

// 			res := ErrorResponse{
// 				Code:       e.StatusCode,
// 				Message:    er.Message(),
// 				FieldError: ErrorFields(er.Cause()),
// 			}

// 			if debugMode {
// 				res.Description = fmt.Sprintf("Error: %v", er)
// 				res.StackTrace = fmt.Sprintf("%+v", errorx.EnsureStackTrace(err))
// 			}

// 			return &res
// 		}
// 	}

// 	return &ErrorResponse{
// 		Code:    http.StatusInternalServerError,
// 		Message: "Unknown server error",
// 	}
// }

// func ErrorFields(err error) []FieldError {
// 	var errs []FieldError

// 	if data, ok := err.(validation.ValidationErrors); ok {
// 		for i, v := range data {
// 			errs = append(errs, FieldError{
// 				Name:        fmt.Sprintf("%d", i+1),
// 				Description: v.Error(),
// 			})
// 		}

// 		return errs
// 	}

// 	return nil
// }

// func SendAuthzResponseErr(w http.ResponseWriter, err error, message any) {
// 	statusCode := http.StatusInternalServerError
// 	code := "internal_error"

// 	if err != nil {
// 		for _, e := range errors.Error {
// 			if errorx.IsOfType(err, e.Type) {
// 				statusCode = e.StatusCode
// 				code = e.Type.Namespace().String()
// 				break
// 			}
// 		}

// 		if message == "Internal Server Error" {
// 			message = err.Error()
// 			code = "unknown_error"
// 		}
// 	}

// 	response := ErrorAuthzResponseFormat{
// 		Status:  "error",
// 		Code:    code,
// 		Message: message,
// 		Metadata: Metadata{
// 			ServerTimestamp: time.Now().UTC().Format("2006-01-02T15:04:05Z"),
// 		},
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(statusCode)
// 	_ = json.NewEncoder(w).Encode(response)
// }
// ```

// ### What's changed

// * Removed the Gin import:

//   ```go
//   import "github.com/gin-gonic/gin"
//   ```

// * Added:

//   ```go
//   import "encoding/json"
//   import "net/http"
//   ```

// * Every function signature changed from:

//   ```go
//   func SendSuccessResponse(ctx *gin.Context, ...)
//   ```

//   to

//   ```go
//   func SendSuccessResponse(w http.ResponseWriter, ...)
//   ```

// * Replaced:

//   ```go
//   ctx.JSON(...)
//   ctx.AbortWithStatusJSON(...)
//   ```

//   with:

//   ```go
//   w.Header().Set("Content-Type", "application/json")
//   w.WriteHeader(statusCode)
//   json.NewEncoder(w).Encode(response)
//   ```

// This is the idiomatic way to send JSON responses when using Chi, since Chi builds directly on Go's standard `net/http` interfaces.
