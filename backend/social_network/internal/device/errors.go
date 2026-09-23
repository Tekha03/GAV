package device

import apperrors "shared/app_errors"

var ErrInvalidToken = apperrors.New(apperrors.Validation, "invalid device token")
