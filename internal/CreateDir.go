package internal

import (
	"log"
	"os"
	"strings"

	"github.com/pkg/errors"
)

func CreateDir(path string) {
	//log.Printf("debug")
	defer func() {
		//log.Printf("CreateDir : panic")
		if err := recover(); err != nil {
			log.Printf("Error: %+v", err)
		}
	}()

	if err := os.Mkdir(path, 0777); err != nil {
		fe, e := "file exists", err.Error()

		log.Printf("debug: len(e): %v", len(e))
		log.Printf("debug: len(fe): %v", len(fe))
		log.Printf("debug: e-fe: %v", len(e)-len(fe))
		log.Printf("debug: e[e-fe:]: %v", e[len(e)-len(fe):])
		if ee := e[len(e)-len(fe):]; ee != fe {
			// 既に存在していたエラーは無視したい。
			if strings.Contains(err.Error(), "Cannot create a file when that file already exists.") {
				log.Printf("debug")
				// 無視
			} else {
				log.Printf("debug")
				log.Printf("Error: %v", err)
				panic(errors.Errorf("%v", err))
			}
		}
	} else {
		//log.Printf("debug")
	}

	//panic(errors.Errorf(path))
}
