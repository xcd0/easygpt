package internal

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/alexflint/go-arg"
)

var ShowUsage func()

// 戻り値は引数処理結果のsettingと引数にあったファイルのリスト
func Argparse() (*Setting, []SettingForDD, []string) {
	var argsAll ArgsAll = ArgsAll{}
	parser := arg.MustParse(&argsAll)
	//log.Printf("argsAll:%v", argsAll)
	ShowUsage = func() {
		buf := new(bytes.Buffer)
		//parser.WriteUsage(os.Stdout) // キャプチャしているつもり。
		parser.WriteUsage(buf)
		strs := strings.Split(buf.String(), " ")
		for i, s := range strs {
			if i == 0 {
				fmt.Printf("%v\n", s)
			} else if i == 1 {
				fmt.Printf(`  $ %v [options...] [files...]

  * 引数は全て順不同。混在可。
  * 引数による設定ではなく、設定ファイルによる指定を推奨する。
  * 設定ファイルは、'--create-setting'によって雛形を生成可能。
  * APIキーは環境変数'OPENAI_API_KEY'に設定する事を推奨。設定ファイルに指定してもよい。
  * 詳しくは '--help' を参照。

OPTIONS:
`, s)
			} else {
				if s[0] == '[' {
					if s[1] == '-' {
						fmt.Printf("    %v ", s)
					} else {
						// これ以降は直書き...orz
						fmt.Printf(`    [--help]
FILES:
    [任意個の入力ファイル、入力ファイルディレクトリ...]
`)
						break
					}
					if s[len(s)-1] == ']' {
						fmt.Printf("\n")
					}
				} else if s[len(s)-1] == ']' {
					fmt.Printf("%v\n", s)
				} else {
					fmt.Printf(" %v", s)
				}
			}
		}
		fmt.Printf("\n")
	}

	// 即終了系の処理 OpenAIに投げないで終了する処理 設定ファイル生成とか
	RunImmidiateTerminate(&argsAll)

	// 設定ファイル探索
	setting, err := GetSetting(&argsAll, parser) // 設定ファイルがなかったときはここでは無視する。nilが入る。
	if err != nil {
		fmt.Printf("%v", err) // 必須パラメータが引数や設定ファイル、環境変数で設定されなかった。
		//parser.WriteUsage(os.Stdout)
		ShowUsage()
		os.Exit(1)
	}

	if len(argsAll.InputText) != 0 { // 引数で質問文が与えられている場合
		// 設定ファイルのAPIキー、URL、temperature、aimodel以外の設定を無視して実行する。
		// つまり、引数で与えられた文字列だけを使用してAIに処理されて返答を得て終了する。
		// ほかの引数で入力ファイルなどが与えられていても無視する。
		RunInputTextOnArgs(&argsAll, setting)
		os.Exit(0)
	}

	// 引数に直接(非フラグ引数として)与えられた入力ファイル毎に設定ファイルを生成する。
	settingsfordd, files := GetSettingsForDD(&argsAll, setting)
	if files == nil {
		return setting, nil, nil //引数ファイルが与えられなかった。
	}

	return setting, settingsfordd, files
}

// 即終了系の処理 OpenAIに投げないで終了する処理
func RunImmidiateTerminate(argsAll *ArgsAll) {
	/*
		// 引数--readmeがあれば、詳細な説明文を出力して終了する。
		if argsAll.Readme {
			ShowDescription()
			os.Exit(0)
		}
	*/

	// 引数で、設定ファイルの雛形を生成する指示があれば生成して終了。
	if argsAll.CreateSetting {
		CreateSettingHjsonTemplate("./easygpt.hjson")
		fmt.Printf("設定ファイルの雛形を以下に出力しました。\n./easygpt.hjson\n出力した設定ファイルの雛形は編集が必要です。\nAPIキーやプロンプトを記載して使用してください。\n")
		os.Exit(0)
	}
}
