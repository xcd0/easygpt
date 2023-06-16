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

func RunWithSetting(setting *Setting) {

	if errInputDir := CheckInputDir(&setting.InputDir); errInputDir != nil {
		fmt.Printf("%v\n\n", errInputDir)
		//parser.WriteUsage(os.Stdout)
		ShowUsage()
		os.Exit(1)
	}
	if errOutputDir := CheckOutputDir(&setting.OutputDir); errOutputDir != nil {
		fmt.Printf("%v\n\n", errOutputDir)
		//parser.WriteUsage(os.Stdout)
		ShowUsage()
		os.Exit(1)
	}

	// 出力先ディレクトリを作成
	CreateSameDirsOn(setting.InputDir, setting.OutputDir)

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

			Ask(&f, setting)

		}()
	}
	wg.Wait()
}

func Ask(input_file *string, setting *Setting) {
	// このファイルのファイルパスから固有の文字列を生成し、一時ディレクトリの名前とする。
	id := strings.ReplaceAll(filepath.ToSlash(strings.Replace(*input_file, setting.InputDir, "", 1)), "/", "-")[1:]
	// 一時ディレクトリにこのファイル固有のディレクトリを作成する
	id = filepath.Join(setting.Tmp, id)
	if len(id) == 0 {
		// なぜか文字がすべて消えた。ファイル名が\とかだったらありえる。
		// ランダム文字列生成してディレクトリ名にする
		uuidObj, _ := uuid.NewUUID()
		id = filepath.Join(setting.Tmp, uuidObj.String())
	}
	CreateDir(id)
	p := filepath.Join(id, "original_path.txt")
	OutputTextForCheck(&p, input_file)

	// 出力
	outputpath := strings.Replace(*input_file, setting.InputDir, setting.OutputDir, 1)
	inputText := GetTextNoError(input_file)
	outputText := QuestionByText(inputText, setting, true)
	OutputTextForCheck(&outputpath, outputText)
}
