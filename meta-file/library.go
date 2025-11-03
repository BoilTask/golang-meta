package metafile

import (
	"path/filepath"
	"strings"
)

func GetFileNameWithoutExt(path string) string {
	filename := filepath.Base(path)
	return strings.TrimSuffix(filename, filepath.Ext(filename))
}
