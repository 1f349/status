//go:build !debug

package res

import (
	"embed"
	"io/fs"
	"path"
)

var (
	//go:embed public
	publicFiles embed.FS
)

func Open(p string) (fs.File, error) {
	return publicFiles.Open(path.Join("public", p))
}
