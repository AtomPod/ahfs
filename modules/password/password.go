package password

import (
	"strings"
	"sync"

	"github.com/czhj/ahfs/modules/setting"
)

type complexity struct {
	ValidChars string
}

var (
	initComplexityOnce sync.Once
	requiredList       []complexity
	charComplexities   = map[string]complexity{
		"lower": {
			"abcdefghijklmnopqrstuvwxyz",
		},
		"upper": {
			"ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		},
		"digit": {
			"0123456789",
		},
		"spec": {
			`!"#$%&'()*+,-./:;<=>?@[\]^_{|}~` + "`",
		},
	}
)

func NewComplexity() {
	initComplexityOnce.Do(func() {
		setupComplexity(setting.PasswordComplexity)
	})
}

func setupComplexity(vals []string) {
	for _, val := range vals {
		requiredList = append(requiredList, charComplexities[val])
	}

	if len(requiredList) == 0 {
		for _, c := range charComplexities {
			requiredList = append(requiredList, c)
		}
	}
}

func IsEnoughComplexity(pwd string) bool {
	NewComplexity()

	for _, required := range requiredList {
		if !strings.ContainsAny(required.ValidChars, pwd) {
			return false
		}
	}

	return true
}
