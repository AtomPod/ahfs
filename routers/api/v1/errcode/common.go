package errcode

const (
	OK                   ErrorCode = 200000
	InternalServerError  ErrorCode = 500000
	UnauthorizedError    ErrorCode = 400000
	ParameterFormatError ErrorCode = 400001
)
