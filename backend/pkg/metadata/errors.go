package metadata

import "github.com/go-kratos/kratos/v2/errors"

const (
	reason string = "UNAUTHORIZED"
)

var (
	ErrNoMetadata          = errors.Unauthorized(reason, "metadata: missing server metadata")
	ErrCodecNotInitialized = errors.Unauthorized(reason, "metadata: codec not initialized")

	ErrNoOperatorHeader = errors.Unauthorized(reason, "metadata: missing operator header")
	ErrInvalidOperator  = errors.Unauthorized(reason, "metadata: invalid operator header")

	ErrNilOperatorMetadata         = errors.Unauthorized(reason, "metadata: nil operator metadata")
	ErrEmptyOperatorMetadataString = errors.Unauthorized(reason, "metadata: empty operator metadata string")
	ErrNilRequestHeader            = errors.Unauthorized(reason, "metadata: nil request headers map")

	ErrSignatureMismatch = errors.Unauthorized(reason, "metadata: signature mismatch")
)
