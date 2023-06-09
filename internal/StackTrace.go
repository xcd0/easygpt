package internal

import (
	"log"

	"github.com/pkg/errors"
)

func StackTrace() {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("%+v", err)
		}
	}()
	panic(errors.Errorf(""))
}
