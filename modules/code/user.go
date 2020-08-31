package code

import (
	"fmt"
	"time"

	"github.com/czhj/ahfs/modules/cache"
	"github.com/czhj/ahfs/modules/log"
	"github.com/czhj/ahfs/modules/setting"
	"github.com/czhj/ahfs/modules/utils"
	"go.uber.org/zap"
)

type emailCodeType string

var (
	emailActiveCode   emailCodeType = "email:%s:active:%s"
	emailResetPwdCode emailCodeType = "email:%s:reset_pwd:%s"
)

func CreateEmailActiveCode(email string) (string, error) {
	return createEmailCode(email, emailActiveCode, setting.Service.ActiveCodeLive)
}

func VerifyEmailActiveCode(email, code string) (bool, error) {
	return verifyEmailCode(email, code, emailActiveCode)
}

func RemoveEmailActiveCode(email, code string) {
	removeEmailCode(email, code, emailActiveCode)
}

func CreateEmailResetPwdCode(email string) (string, error) {
	return createEmailCode(email, emailResetPwdCode, setting.Service.ResetPasswordCodeLive)
}

func VerifyEmailResetPwdCode(email, code string) (bool, error) {
	return verifyEmailCode(email, code, emailResetPwdCode)
}

func RemoveEmailResetPwdCode(email, code string) {
	removeEmailCode(email, code, emailResetPwdCode)
}

func createEmailCode(email string, typ emailCodeType, d time.Duration) (string, error) {
	code := utils.RandomCode()
	key := fmt.Sprintf(string(typ), email, code)
	if err := cache.Set(key, true, d); err != nil {
		return "", err
	}
	return code, nil
}

func verifyEmailCode(email string, code string, typ emailCodeType) (bool, error) {
	key := fmt.Sprintf(string(typ), email, code)

	var exists bool
	if err := cache.Get(key, &exists); err != nil || !exists {
		if exists || err == cache.ErrNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func removeEmailCode(email, code string, typ emailCodeType) {
	key := fmt.Sprintf(string(typ), email, code)

	if err := cache.Delete(key); err != nil {
		if err != cache.ErrNotFound {
			log.Error("Failed to delete cache key", zap.String("key", key), zap.Error(err))
		}
	}
}
