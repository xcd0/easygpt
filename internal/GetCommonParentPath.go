package internal

import (
	"path/filepath"
	"strings"
)

// https://go.dev/play/p/0atgRGr8nhh

func GetCommonParentPath(paths []string) string {
	if len(paths) == 0 {
		return ""
	}
	commonParts := strings.Split(filepath.ToSlash(paths[0]), "/")
	for _, path := range paths[1:] {
		pathParts := strings.Split(filepath.ToSlash(path), "/")
		for i := len(commonParts); i > 0; i-- {
			if strings.Join(commonParts[:i], "/") == strings.Join(pathParts[:i], "/") {
				commonParts = commonParts[:i]
				break
			}
			if i == 1 {
				return "/"
			}
		}
	}
	return strings.Join(commonParts, "/")
}
