package main

import (
	"easygpt/internal"
	"log"
)

func main() {
	log.SetFlags(log.Ltime | log.Lshortfile) // ログの出力書式を設定する

	// 引数処理
	settings, settingsfordd := internal.Argparse()

	// D&Dで実行されたとき用一括実行
	if len(settingsfordd) != 0 {
		internal.RunNoSetting(settingsfordd)
	}
	if len(settings) != 0 {
		// 引数を設定して実行されたとき用一括実行
		internal.RunWithSetting(settings)
	}

}
