package api

import (
	"github.com/go-openapi/strfmt"
	"github.com/nvloff/f3_payments_service/gen/models"
)

const DomainErrorCode strfmt.UUID = "08e64539-3b7c-45fb-9470-68ec736ef9de"

// Convert go err to API response Error
func NewApiError(err error) *models.APIError {
	return &models.APIError{
		ErrorCode:    DomainErrorCode,
		ErrorMessage: err.Error(),
	}
}
