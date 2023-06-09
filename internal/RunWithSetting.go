package internal

import (
	"fmt"
	"log"
	"os"
	"sync"
)

func RunWithSetting(settings []Setting) {
	//log.Printf("settings: %v", settings)
	//StackTrace()

	for _, setting := range settings {
		// 出力先ディレクトリを作成
		CreateSameDirsOn(setting.OutputDir, setting.InputDir)

		// 入力ディレクトリの中のテキストファイルを検索
		// 拡張子指定があれば、その拡張子のファイルのみ使用する。
		files := GetTargetFiles(setting.InputDir, setting.Extension)
		if len(files) == 0 {
			fmt.Println("入力ディレクトリの中に処理対象となるテキストファイルがありませんでした。\n終了します。")
			os.Exit(0)
		}

		// 一時ファイルを保存するためのディレクトリを作成する。
		os.RemoveAll(setting.Tmp)
		if err := os.Mkdir(setting.Tmp, 0777); err != nil {
			fe, e := "file exists", err.Error()
			if ee := e[len(e)-len(fe):]; ee != fe {
				log.Printf("Error: %v", err)
			}
		}

		var wg sync.WaitGroup
		sem := make(chan struct{}, setting.Concurrency) // concurrency数のバッファ
		// ファイルごとの処理
		for _, f := range files {
			//log.Printf("%v:%v", i, f)

			sem <- struct{}{}
			wg.Add(1)
			go func() {
				defer wg.Done()
				defer func() { <-sem }() // 処理が終わったらチャネルを解放

				Ask(f, setting.InputDir, setting.OutputDir, setting.Tmp, setting.ApiKey, &setting.Prompt)

			}()
		}
		wg.Wait()
	}
}
