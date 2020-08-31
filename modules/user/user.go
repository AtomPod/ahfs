package user

import (
	"os"
	"os/user"
	"runtime"
	"strings"
)

func CurrentUser() string {
	u, err := user.Current()
	if err != nil {
		return fallbackCurrentUser()
	}

	username := u.Username
	if runtime.GOOS == "windows" {
		parts := strings.Split(username, "\\")
		username = parts[len(parts)-1]
	}
	return username
}

func fallbackCurrentUser() string {
	username := os.Getenv("USER")
	if len(username) > 0 {
		return username
	}
	return os.Getenv("USERNAME")
}
