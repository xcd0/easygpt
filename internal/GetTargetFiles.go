package internal

import (
	"log"
	"path/filepath"

	"github.com/pkg/errors"
)

func GetTargetFiles(dir string, extension string) []string {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Error: %+v", err)
		}
	}()
	// 入力ディレクトリの中のテキストファイルを検索
	fes, err := Dirwalk(dir)
	if err != nil {
		//log.Fatalf("Error: %v", err)
		panic(errors.Errorf("%v", err))
	}

	// 拡張子指定があれば、その拡張子のファイルのみ使用する。
	files := []string{}
	//log.Printf("extension: %v", extension)
	if extension == "" {
		return fes.Files
	} else {
		for _, f := range fes.Files {
			e := filepath.Ext(f)
			if e == extension {
				// ok
				files = append(files, f)
			}
		}
	}
	return files
}
