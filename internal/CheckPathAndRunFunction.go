package internal

import "os"

// 引数で与えられたpathを調べ、
// 存在しなければfuncNotExistを実行し、
// ディレクトリであればfuncIsDirを実行し、
// ファイルであればfuncIsFileを実行する関数。
// いい名前があればリネームしたい。
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
