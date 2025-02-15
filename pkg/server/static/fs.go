package static

import (
	"embed"

	"github.com/benbjohnson/hashfs"
)

//go:embed style.css
var fs embed.FS
var HashFS *hashfs.FS = hashfs.NewFS(fs)

func GetStaticPath(resource string) string {
	return "/static/" + HashFS.HashName(resource)
}
