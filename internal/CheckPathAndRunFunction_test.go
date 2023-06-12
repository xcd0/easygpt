package internal

import (
	"log"
	"path/filepath"
	"testing"
)

func TestCheckPathAndRunFunction(t *testing.T) {
	log.SetFlags(log.Ltime | log.Lshortfile) // ログの出力書式を設定する

	paths := []string{
		"CheckPathAndRunFunction.go",
		"../internal",
		"../internal/hoge",
	}
	for i, p := range paths {
		p, _ := filepath.Abs(p)
		CheckPathAndRunFunction(
			p,
			func(err error) {
				t.Logf("Error: %v", err)
			},
			func() {
				t.Logf("%v: dir : %v", i, p)
			},
			func() {
				t.Logf("%v: file: %v", i, p)
			},
		)
	}

}
