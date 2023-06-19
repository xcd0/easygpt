package internal

import (
	"log"
	"os"
)

func OutputTextForCheck(file, str *string) {
	f, err := os.Create(*file)
	if err != nil {
		log.Printf("Warning: %v", err)
		//panic(err)
	} else {
		if _, err := f.Write([]byte(*str)); err != nil {
			log.Printf("Warning: %v", err)
		}
	}
	//log.Printf("Debug: output : %v", file)
}
