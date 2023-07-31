package winfsinjector

import (
	"os"
)

func SetReadFile(f func(string) ([]byte, error)) {
	readFile = f
}

func ResetReadFile() {
	readFile = os.ReadFile
}

func SetRemoveAll(f func(string) error) {
	removeAll = f
}

func ResetRemoveAll() {
	removeAll = os.RemoveAll
}
