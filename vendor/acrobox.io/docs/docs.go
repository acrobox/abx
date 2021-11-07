package docs // import "acrobox.io/docs"

import (
	"embed"
	"io/fs"
)

//go:embed docs
var docs embed.FS

// FS is a io/fs.FS that contains the Acrobox documentation.
var FS fs.FS

func init() {
	fs, err := fs.Sub(docs, "docs")
	if err != nil {
		panic(err)
	}
	FS = fs
}
