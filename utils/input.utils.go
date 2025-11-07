package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/etsrohan/Rohan-Srivastava_Golang-Backend-Practical-Task/config"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func StrictBindJSON[T any](ctx *gin.Context, obj *T) *[]config.APIError {
	var validationErrs []config.APIError

	// Read raw request body bytes
	bodyBytes, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		validationErrs = append(validationErrs, config.APIError{Code: http.StatusBadRequest, Message: "Failed to read request body"})
		return &validationErrs
	}

	// Reset the body for further use downstream
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	// Check for unknown fields via findUnknownFields
	unknownFields, err := findUnknownFields[T](bodyBytes)
	if err != nil {
		validationErrs = append(validationErrs, config.APIError{Code: http.StatusBadRequest, Message: "Malformed JSON"})
		return &validationErrs
	}
	if len(unknownFields) > 0 {
		for _, f := range unknownFields {
			validationErrs = append(validationErrs, config.APIError{
				Code:    http.StatusBadRequest,
				Message: f,
			})
		}
		return &validationErrs
	}

	// Decode JSON into the struct normally
	decoder := json.NewDecoder(bytes.NewBuffer(bodyBytes))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(obj); err != nil {
		validationErrs = append(validationErrs, config.APIError{Code: http.StatusBadRequest, Message: err.Error()})
		return &validationErrs
	}

	// Validate struct fields via validator
	validate := validator.New()
	if err := validate.Struct(obj); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			validationErrs = append(validationErrs, config.APIError{Code: http.StatusBadRequest, Message: err.Error()})
		}
	}

	if len(validationErrs) > 0 {
		return &validationErrs
	}
	return nil
}

// Workaround for catching multiple unknown fields
func findUnknownFields[T any](jsonData []byte) ([]string, error) {
	var dataMap map[string]interface{}
	if err := json.Unmarshal(jsonData, &dataMap); err != nil {
		return nil, err
	}

	var knownFields []string
	var t T
	rt := reflect.TypeOf(t)
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}

	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		tag := field.Tag.Get("json")
		if tag == "-" || tag == "" {
			continue
		}
		tagName := strings.Split(tag, ",")[0]
		knownFields = append(knownFields, tagName)
	}

	var unknownFields []string
	for key := range dataMap {
		found := false
		for _, known := range knownFields {
			if key == known {
				found = true
				break
			}
		}
		if !found {
			unknownFields = append(unknownFields, key)
		}
	}
	return unknownFields, nil
}
