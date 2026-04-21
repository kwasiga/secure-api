// Package validator provides request decoding and struct validation helpers.
package validator

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

type validationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Validate decodes the JSON request body into dst and runs struct validation.
// On failure it writes the appropriate error response and returns a non-nil error,
// allowing handlers to return early with a single if-check.
func Validate(w http.ResponseWriter, r *http.Request, dst any) error {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}
	if err := validate.Struct(dst); err != nil {
		if validateErrs, ok := err.(validator.ValidationErrors); ok {
			var errs []validationError
			for _, ve := range validateErrs {
				errs = append(errs, validationError{
					Field:   ve.Field(),
					Message: ve.Error(),
				})
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(errs)
		}
		return err
	}
	return nil
}
