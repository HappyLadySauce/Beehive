package core

import (
	"net/http"
	"strings"

	"github.com/HappyLadySauce/errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"k8s.io/klog/v2"

	"github.com/HappyLadySauce/Beehive/internal/pkg/code"
)

// ErrResponse defines the return messages when an error occurred.
// Reference will be omitted if it does not exist.
// swagger:model
type ErrResponse struct {
	// Code defines the business error code.
	Code int `json:"code"`

	// Message contains the detail of this message.
	// This message is suitable to be exposed to external
	Message string `json:"message"`

	// Reference returns the reference document which maybe useful to solve this error.
	Reference string `json:"reference,omitempty"`

	// Details contains detailed validation error information.
	// This field is only present when validation errors occur.
	Details map[string]string `json:"details,omitempty"`
}

// SuccessResponse defines the return messages when an operation succeeds.
// swagger:model
type SuccessResponse struct {
	// Code defines the business success code.
	Code int `json:"code"`

	// Message contains the success message.
	Message string `json:"message"`
}

// FormatValidationError formats validator.ValidationErrors into a user-friendly map.
// Returns a map where keys are field names (JSON field names) and values are error messages.
func FormatValidationError(err error) map[string]string {
	details := make(map[string]string)

	// Token format: "validation.<tag>|k=v|k2=v2"
	// Frontend should parse and translate with i18n.
	toToken := func(tag string, kv map[string]string) string {
		key := "validation." + tag
		if len(kv) == 0 {
			return key
		}
		parts := []string{key}
		for k, v := range kv {
			if k == "" || v == "" {
				continue
			}
			parts = append(parts, k+"="+v)
		}
		return strings.Join(parts, "|")
	}

	// Check if the error is a validator.ValidationErrors
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrors {
			// Get JSON field name from struct field tag
			jsonFieldName := getJSONFieldName(fieldError)

			// Generate user-friendly error message based on the validation tag
			var message string
			switch fieldError.Tag() {
			case "required":
				message = toToken("required", nil)
			case "min":
				message = toToken("min", map[string]string{"min": fieldError.Param()})
			case "max":
				message = toToken("max", map[string]string{"max": fieldError.Param()})
			case "email":
				message = toToken("email", nil)
			case "len":
				message = toToken("len", map[string]string{"len": fieldError.Param()})
			case "oneof":
				// validator oneof param is space separated: "a b c"
				values := strings.ReplaceAll(fieldError.Param(), " ", ",")
				message = toToken("oneof", map[string]string{"values": values})
			default:
				message = toToken("invalid", map[string]string{"tag": fieldError.Tag()})
			}

			details[jsonFieldName] = message
		}
	} else {
		// If it's not a ValidationErrors, just use the error message
		details["error"] = err.Error()
	}

	return details
}

// getJSONFieldName extracts the JSON field name from a validation error.
// It uses the namespace to get the field name and converts it to camelCase JSON field name.
func getJSONFieldName(fieldError validator.FieldError) string {
	// Use Namespace() to get the full path (e.g., "CreateWGPeerRequest.AllowedIPs")
	namespace := fieldError.Namespace()
	if namespace != "" {
		// Parse namespace to get the field name (last part after the dot)
		parts := strings.Split(namespace, ".")
		if len(parts) > 0 {
			fieldName := parts[len(parts)-1]
			// Convert struct field name to JSON field name (camelCase)
			// e.g., "AllowedIPs" -> "allowedIPs"
			if len(fieldName) > 0 {
				jsonFieldName := strings.ToLower(fieldName[:1]) + fieldName[1:]
				return jsonFieldName
			}
		}
	}

	// Fallback: use Field() method and convert to camelCase
	fieldName := fieldError.Field()
	if idx := strings.Index(fieldName, "."); idx != -1 {
		fieldName = fieldName[idx+1:]
	}
	if len(fieldName) > 0 {
		jsonFieldName := strings.ToLower(fieldName[:1]) + fieldName[1:]
		return jsonFieldName
	}

	return fieldName
}

// WriteResponseWithDetails write an error or the response data into http response body with details.
// The details parameter can be used to provide additional error information, such as validation errors.
func WriteResponseWithDetails(c *gin.Context, err error, data interface{}, details map[string]string) {
	if err != nil {
		klog.Errorf("%#+v", err)
		coder := errors.ParseCoder(err)

		response := ErrResponse{
			Code:      coder.Code(),
			Message:   coder.String(),
			Reference: coder.Reference(),
		}

		// Only include details if there are validation errors
		if len(details) > 0 {
			response.Details = details
		}

		c.JSON(coder.HTTPStatus(), response)
		return
	}

	// If data is nil, return a default success response
	if data == nil {
		response := SuccessResponse{
			Code:    code.ErrSuccess,
			Message: code.Message(code.ErrSuccess),
		}
		c.JSON(http.StatusOK, response)
		return
	}

	c.JSON(http.StatusOK, data)
}

// WriteResponse write an error or the response data into http response body.
// It use errors.ParseCoder to parse any error into errors.Coder
// errors.Coder contains error code, user-safe error message and http status code.
func WriteResponse(c *gin.Context, err error, data interface{}) {
	WriteResponseWithDetails(c, err, data, nil)
}

// WriteResponseBindErr writes a validation error response with details.
// It use FormatValidationError to format the validation errors into a map.
// Then it use WriteResponseWithDetails to write the response with the details.
func WriteResponseBindErr(c *gin.Context, err error, data interface{}) {
	details := FormatValidationError(err)
	WriteResponseWithDetails(c, errors.WithCode(code.ErrBind, "Validation failed"), data, details)
}
