package main

import (
	"log"

	"github.com/xcd0/easygpt/internal"
)

func main() {
	log.SetFlags(log.Ltime | log.Lshortfile)             // ログの出力書式を設定する
	setting, settingsfordd, files := internal.Argparse() // 引数処理
	if len(settingsfordd) != 0 {
		internal.RunByDD(settingsfordd, files, setting) // D&Dで実行されたとき
	} else if setting != nil {
		internal.RunWithSetting(setting) // コマンドラインから実行されたとき
	}

}
