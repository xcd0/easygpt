package internal

import "os"

func CheckPathAndRunFunction(
	path string,
	funcNotExist func(err error),
	funcIsDir func(),
	funcIsFile func(),
) {
	if fi, err := os.Stat(path); err != nil {
		funcNotExist(err)
	} else if fi.IsDir() {
		funcIsDir()
	} else {
		funcIsFile()
	}
}
