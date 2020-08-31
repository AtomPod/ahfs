package utils

import (
	"math/rand"
	"strconv"
)

func RandomCode() string {
	return strconv.FormatUint(uint64(rand.Uint32()%900000+100000), 10)
}
