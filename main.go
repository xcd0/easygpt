package main

import (
	_ "embed"
	"log"

	"github.com/xcd0/easygpt/internal"
)

//go:embed readme.md
var readme string

func main() {
	log.SetFlags(log.Ltime | log.Lshortfile)             // ログの出力書式を設定する
	setting, settingsfordd, files := internal.Argparse( &readme) // 引数処理
	if len(settingsfordd) != 0 {
		internal.RunByDD(settingsfordd, files, setting) // D&Dで実行されたとき
	} else if setting != nil {
		internal.RunWithSetting(setting) // コマンドラインから実行されたとき
	}

}
