package internal

import (
	"fmt"
	"log"
	"os/exec"
)

func RunOnBash(format string, a ...any) string {
	cmd := fmt.Sprintf(format, a...)
	ret, _ := exec.Command("bash", "-c", cmd).Output()
	str := string(ret)
	log.Printf("\n$ %v\n%v", cmd, str)
	return str
}
