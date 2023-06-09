package internal

import (
	"log"
	"os"

	"github.com/pkg/errors"
)

func CreateDir(path string) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Error: %+v", err)
		}
	}()

	if err := os.Mkdir(path, 0777); err != nil {
		fe, e := "file exists", err.Error()
		if ee := e[len(e)-len(fe):]; ee != fe {
			panic(errors.Errorf("%v", err))
		}
	}

	//panic(errors.Errorf(path))
}
