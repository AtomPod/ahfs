package sso

import (
	"reflect"

	"github.com/czhj/ahfs/modules/log"
	"go.uber.org/zap"
)

var (
	ssoMethods = []SingleSignOn{
		&Session{},
		&Auth{},
	}
)

func Methods() []SingleSignOn {
	return ssoMethods
}

func Init() {
	for _, method := range Methods() {
		if err := method.Init(); err != nil {
			log.Error("Could not initialize SSO method", zap.String("name", reflect.TypeOf(ssoMethods).String()), zap.Error(err))
		}
	}
}

func Close() {
	for _, method := range Methods() {
		if err := method.Close(); err != nil {
			log.Error("Could not close SSO method", zap.String("name", reflect.TypeOf(ssoMethods).String()), zap.Error(err))
		}
	}
}

func Register(m SingleSignOn) {
	ssoMethods = append(ssoMethods, m)
}
