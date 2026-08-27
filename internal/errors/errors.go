package errors

import (
	"net/http"

	"github.com/joomcode/errorx"
)

type ErrorType struct {
	StatusCode int
	Type       *errorx.Type
}

var Error = []ErrorType{

	{
		StatusCode: http.StatusBadRequest,
		Type:       ErrInvalidUserInput,
	},

	{
		StatusCode: http.StatusBadRequest,
		Type:       ErrTextractFailed,
	},

	{
		StatusCode: http.StatusBadRequest,
		Type:       ErrInvalidDocument,
	},

	{
		StatusCode: http.StatusBadRequest,
		Type:       ErrUnableToGet,
	},
}

var (
	invalidInput    = errorx.NewNamespace("validation error").ApplyModifiers(errorx.TypeModifierOmitStackTrace)
	textractfailed  = errorx.NewNamespace("Textract Failed.")
	invaliddocument = errorx.NewNamespace("invalid document")
	unabletoget     = errorx.NewNamespace("Error_UNABLE_TO_GET")
)

var (
	ErrInvalidUserInput = errorx.NewType(invalidInput, "invalid user input")
	ErrTextractFailed   = errorx.NewType(textractfailed, "Error Textract failed.")
	ErrInvalidDocument  = errorx.NewType(invaliddocument, "Error invalid document.")
	ErrUnableToGet      = errorx.NewType(unabletoget, "Error unable to get.")
)
