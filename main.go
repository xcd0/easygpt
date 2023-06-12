package main

import (
	"fmt"
	"log"

	"github.com/xcd0/easygpt/internal"
)

func main() {
	log.SetFlags(log.Ltime | log.Lshortfile) // ログの出力書式を設定する

	// 引数処理
	settings, settingsfordd := internal.Argparse()

	//log.Println(settings)
	//log.Println(settingsfordd)

	// D&Dで実行されたとき用一括実行
	if len(settingsfordd) != 0 {
		fmt.Println("オプションなし引数モードで実行します。\n")
		internal.RunNoSetting(settingsfordd)
	} else if len(settings) != 0 {
		fmt.Println("オプションあり引数モードで実行します。\n")
		// 引数を設定して実行されたとき用一括実行
		internal.RunWithSetting(settings)
	}

}
