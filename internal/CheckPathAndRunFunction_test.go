package internal

import (
	"log"
	"path/filepath"
	"testing"
)

func TestCheckPathAndRunFunction(t *testing.T) {

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
				log.Printf("Error: %v", err)
			},
			func() {
				log.Printf("%v: dir : %v", i, p)
			},
			func() {
				log.Printf("%v: file: %v", i, p)
			},
		)
	}

}
