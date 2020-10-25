package storage

import (
	"fmt"
)

func RegisterStorageGenerator(typ Type, generator Generator) {
	if _, ok := storageMap[typ]; ok {
		panic(fmt.Sprintf("storage: cannot register generator twice for %s", typ))
	}

	storageMap[typ] = generator
}
