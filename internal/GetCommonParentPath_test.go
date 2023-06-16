package internal

import (
	"log"
	"testing"
)

func TestGetCommonParentPath(t *testing.T) {
	log.SetFlags(log.Ltime | log.Lshortfile) // ログの出力書式を設定する

	files := []string{
		"/hoge/a/foo",
		"/hoge/a/b/bar",
		"/hoge/a/c/hoge",
		"/hoge/a/b/c/piyo",
	}
	p := GetCommonParentPath(files)
	log.Printf("%v", p)
}
