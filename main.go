package main

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/xcd0/easygpt/internal"
)

//go:embed readme.md
var readme string

func main() {

	// Ctrl-Cを捕捉する。
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	done := make(chan interface{})
	go run(done)

	select {
	case <-done:
		//fmt.Println("すべての処理が完了")
		os.Exit(0)
	case s := <-c:
		// シグナル受信時
		fmt.Printf("割り込み終了: %v\n", s)
		os.Exit(1)
	}
}

func run(done chan interface{}) {
	log.SetFlags(log.Ltime | log.Lshortfile) // ログの出力書式を設定する

	setting, settingsfordd, files, argsAll := internal.Argparse(&readme) // 引数処理

	//log.Println("setting:\n%v", setting)
	//log.Println("settingdd:\n%v", settingsfordd)
	//log.Println("files:\n%v", files)
	//log.Println("args:\n%v", argsAll)

	if len(argsAll.InputText) != 0 {
		// 引数で質問文が与えられている場合
		// 設定ファイルのAPIキー、URL、temperature、aimodel以外の設定を無視して実行する。
		// つまり、引数で与えられた文字列だけを使用してAIに処理されて返答を得て終了する。
		// ほかの引数で入力ファイルなどが与えられていても無視する。
		internal.RunInputTextOnArgs(argsAll, setting)
	} else if len(settingsfordd) != 0 {
		// D&Dで実行されたとき
		internal.RunByDD(settingsfordd, files, setting)
	} else if setting != nil {
		// コマンドラインから実行されたとき
		internal.RunWithSetting(setting)
	}

	done <- struct{}{}
}
