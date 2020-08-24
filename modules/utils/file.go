package utils

import (
	"path/filepath"
	"strconv"

	"github.com/rs/xid"
)

func GetFilename(filename string) string {
	return filepath.Base(filename)
}

func GenerateFileID(uid uint) string {
	return strconv.FormatUint(uint64(uid), 16) + xid.New().String()
}

func GenerateLockID(fid uint) string {
	return strconv.FormatUint(uint64(fid), 16) + xid.New().String()
}
