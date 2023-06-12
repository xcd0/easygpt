package internal

import (
	"bytes"
	"encoding/json"
	"log"
)

func JsonFormat(body []byte) string {
	var buf bytes.Buffer
	err := json.Indent(&buf, body, "", "\t")
	if err != nil {
		log.Printf("Error: %v", err)
	}
	return buf.String()
}
