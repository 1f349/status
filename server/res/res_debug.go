//go:build debug

package res

import (
	"os"
	"path"
)

func Open(p string) (*os.File, error) {
	return os.Open(path.Join("res/public", p))
}
