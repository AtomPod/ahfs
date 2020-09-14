package validator

import (
	"regexp"

	"code.gitea.io/gitea/modules/password"
	"github.com/czhj/ahfs/modules/log"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

var (
	validateFuncMap = map[string]validator.Func{
		"password": Password,
		"username": Username,
		"nickname": Nickname,
	}

	usernameRegexp = regexp.MustCompile("^[a-zA-Z0-9_-]{6,16}$")
	nicknameRegexp = regexp.MustCompile("^[^\\s]{1,16}$")
)

func Password(fl validator.FieldLevel) bool {
	pwd := fl.Field().String()
	return password.IsComplexEnough(pwd)
}

func Username(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	return usernameRegexp.MatchString(username)
}

func Nickname(fl validator.FieldLevel) bool {
	nickname := fl.Field().String()
	return nicknameRegexp.MatchString(nickname)
}

func Register() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		setupValidator(v)
	} else {
		log.Fatal("Cannot convert validator engine to v10")
	}
}

func setupValidator(v *validator.Validate) {
	for name, fn := range validateFuncMap {
		if err := v.RegisterValidation(name, fn); err != nil {
			log.Warn("Failed to register validation", zap.String("name", name), zap.Error(err))
		} else {
			log.Debug("Register validation for gin", zap.String("name", name))
		}
	}
}
