package errcode

import (
	"github.com/czhj/ahfs/modules/log"
	"go.uber.org/zap"
)

// 错误码格式 x yyy zz 开头第一位为错误级别
// 2 成功
// 4 应用错误/参数
// 5 系统错误
// 后三位为服务类型，服务类型全000为特殊通用，不允许使用，最后两位为错误码序号
type ErrorCode int

func (e ErrorCode) Message() string {
	msg, ok := errorMessage[e]
	if !ok {
		msg = "Unknow error"
	}
	return msg
}

var (
	errorMessage = map[ErrorCode]string{}
)

func addErrorMessage(e ErrorCode, m string) {
	if _, ok := errorMessage[e]; ok {
		log.Fatal("Error code register message twice, conflicts may occur.", zap.Int("error_code", int(e)), zap.String("message", m))
	}
	errorMessage[e] = m
}
