package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"aidoc/internal/errors"

	validation "github.com/go-playground/validator/v10"
	"github.com/joomcode/errorx"
	"github.com/spf13/viper"
)

func SendSuccessResponse(w http.ResponseWriter, statusCode int, code Code, data any, message any) {

	response := SuccessResponse{
		Status:  "success",
		Code:    string(code),
		Message: message,
		Data:    data,
		Metadata: Metadata{
			ServerTimestamp: time.Now().UTC().Format("2006-01-02T15:04:05Z"),
		},
	}
	w.Header().Set("Content-Type", "application/hson")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(response)

}

func SendErrorResponse(w http.ResponseWriter, err *ErrorResponse) {
	w.Header().Set("Contet-Type", "application/json")
	w.WriteHeader(err.Code)
	_ = json.NewEncoder(w).Encode(err)
}

func SendErrorResponseFormated(w http.ResponseWriter, err error) {
	statusCode := http.StatusInternalServerError
	message := "Internal Server Error"
	code := "internal_error"
	if err != nil {
		if validationErrors, ok := err.(validation.ValidationErrors); ok {
			statusCode = http.StatusBadRequest
			message = "Validation error"
			code = "validation_error"
			for _, fieldError := range validationErrors {
				message = fieldError.Field() + ": " + fieldError.ActualTag()
			}
		} else {
			for _, e := range errors.Error {
				if errorx.IsOfType(err, e.Type) {
					statusCode = e.StatusCode
					message = strings.ReplaceAll(e.Type.FullName(), fmt.Sprintf("%s.", e.Type.Namespace().String()), "")

					if strings.Contains(err.Error(), "attempt left") || strings.Contains(err.Error(), "attempts left") {
						message = strings.ReplaceAll(err.Error(), "CUSTOMER_NOT_VERIFIED.invalid pin", "Invalid Pin")
					}

					code = e.Type.Namespace().String()
					break
				}
			}

			if message == "Internal Server Error" {
				message = err.Error()
				code = "unknown_error"
			}
		}
	}

	response := ErrorResponseFormat{
		Status:  "error",
		Code:    code,
		Message: message,
		Metadata: Metadata{
			ServerTimestamp: time.Now().UTC().Format("2006-01-02T15:04:05Z"),
		},
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(response)

}

func GetErrorFrom(err error) *ErrorResponse {
	debugMode := viper.GetBool("debug")

	for _, e := range errors.Error {
		if errorx.IsOfType(err, e.Type) {
			er := errorx.Cast(err)
			res := ErrorResponse{
				Code:       e.StatusCode,
				Message:    er.Message(),
				FieldError: ErrorFields(er.Cause()),
			}

			if debugMode {
				res.Description = fmt.Sprintf("Error: %v", er)
				res.StackTrace = fmt.Sprintf("%+v", errorx.EnsureStackTrace(err))
			}

			return &res
		}
	}

	return &ErrorResponse{
		Code:    http.StatusInternalServerError,
		Message: "Unknown server error",
	}
}

func ErrorFields(err error) []FieldError {
	var errs []FieldError

	if data, ok := err.(validation.ValidationErrors); ok {
		for i, v := range data {
			errs = append(errs, FieldError{
				Name:        fmt.Sprintf("%d", i+1),
				Description: v.Error(),
			},
			)
		}

		return errs
	}

	return nil
}

func SendAuthzResponseErr(w http.ResponseWriter, err error, message any) {

	statusCode := http.StatusInternalServerError
	code := "internal_error"

	if err != nil {
		for _, e := range errors.Error {
			if errorx.IsOfType(err, e.Type) {
				statusCode = e.StatusCode
				code = e.Type.Namespace().String()
				break
			}
		}

		if message == "Internal Server Error" {
			message = err.Error()
			code = "unknown_error"
		}

	}

	response := ErrorAuthzResponseFormat{
		Status:  "error",
		Code:    code,
		Message: message,
		Metadata: Metadata{
			ServerTimestamp: time.Now().UTC().Format("2006-01-02T15:04:05Z"),
		},
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(response)

}
