package internal

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/google/uuid"
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
		if setting.Tmp == "" {
			setting.Tmp = "tmp"
		}
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

				Ask(&f, &setting.InputDir, &setting.OutputDir, &setting.Tmp, &setting.OpenaiURL, &setting.AiModel, &setting.Apikey, &setting.Prompt, setting.Temperature)

			}()
		}
		wg.Wait()
	}
}

func Ask(input_file, input_dir, output_dir, tmp_dir, openaiURL, aiModel, apikey, prompt *string, temperature float64) {
	// このファイルのファイルパスから固有の文字列を生成し、一時ディレクトリの名前とする。
	id := strings.ReplaceAll(filepath.ToSlash(strings.Replace(*input_file, *input_dir, "", 1)), "/", "-")[1:]
	// 一時ディレクトリにこのファイル固有のディレクトリを作成する
	id = filepath.Join(*tmp_dir, id)
	if len(id) == 0 {
		// なぜか文字がすべて消えた。ファイル名が\とかだったらありえる。
		// ランダム文字列生成してディレクトリ名にする
		uuidObj, _ := uuid.NewUUID()
		id = filepath.Join(*tmp_dir, uuidObj.String())
	}
	CreateDir(id)
	p := filepath.Join(id, "original_path.txt")
	OutputTextForCheck(&p, input_file)

	// 出力
	outputpath := strings.Replace(*input_file, *input_dir, *output_dir, 1)

	output := Question(openaiURL, aiModel, apikey, prompt, input_file, &id, temperature, true)

	OutputTextForCheck(&outputpath, output)
}
