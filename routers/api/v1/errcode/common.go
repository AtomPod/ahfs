package errcode

const (
	OK                   ErrorCode = 200000
	InternalServerError  ErrorCode = 500000
	UnauthorizedError    ErrorCode = 400000
	ParameterFormatError ErrorCode = 400001
	PermissionDenied     ErrorCode = 400002
	VisitTooFrequently   ErrorCode = 400003
)
