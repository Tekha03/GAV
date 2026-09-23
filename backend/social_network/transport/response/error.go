package response

import apperrors "shared/app_errors"

type ErrorResponse struct {
	Error ErrorBody `json:"error"`
}

type ErrorBody struct {
	Code     apperrors.Code     `json:"code"`
	Category apperrors.Category `json:"category"`
	Message  string             `json:"message"`
	Details  map[string]any     `json:"details,omitempty"`
}
