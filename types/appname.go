package types

import (
	"os"
	"path"
)

func ExecutableName() string {
	if name, err := os.Executable(); err == nil {
		return path.Base(name)
	} else {
		panic(err)
	}
}
