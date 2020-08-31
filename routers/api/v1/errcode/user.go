package errcode

const (
	EmailAlreadyExists     ErrorCode = 400101
	EmailNotFound          ErrorCode = 400102
	UsernameAlreadyExists  ErrorCode = 400103
	UsernameNotFound       ErrorCode = 400104
	EmailActiveCodeError   ErrorCode = 400105
	ResetPwdCodeError      ErrorCode = 400106
	IncorrectUserNameOrPwd ErrorCode = 400107
	EmailFormatError       ErrorCode = 400108
	EmailResetPwdCodeError ErrorCode = 400109
)
