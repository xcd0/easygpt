package internal

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/alexflint/go-arg"
)

const (
	Exit_success                   = iota
	Error_incomplete_configuration // 必須パラメータが引数や設定ファイル、環境変数で設定されなかった。
)

var ShowUsage func()

// 戻り値は引数処理結果のsettingと引数にあったファイルのリスト
func Argparse(readme *string) (*Setting, []SettingForDD, []string, *ArgsAll) {
	var argsAll ArgsAll = ArgsAll{}
	parser := arg.MustParse(&argsAll)
	//log.Printf("argsAll:%v", argsAll)
	// 関数として定義しないのは、文字列をキャプチャして、引数なしで呼べるようにするため。
	ShowUsage = func() { fmt.Printf("%v", *CreateUsageText(parser)) } // parser.WriteUsage(buf)を呼んで文字列を編集している。

	// 即終了系の処理 OpenAIに投げないで終了する処理 設定ファイル生成とか
	// なんかいい名前があったらリネームしたい
	RunImmidiateTerminate(&argsAll, readme)

	// 設定ファイル探索
	setting, err := GetSetting(&argsAll, parser) // 設定ファイルがなかったときはここでは無視する。nilが入る。
	if err != nil {
		fmt.Printf("%v", err) // 必須パラメータが引数や設定ファイル、環境変数で設定されなかった。
		log.Printf("setting : %v", setting)
		ShowUsage() // parser.WriteUsage(os.Stdout) をいい感じに見た目を弄った関数を呼び出す
		os.Exit(1)
	}

	// 引数に直接(非フラグ引数として)与えられた入力ファイル毎に設定ファイルを生成する。
	settingsfordd, files := GetSettingsForDD(&argsAll, setting)
	if files == nil {
		return setting, nil, nil, &argsAll //引数ファイルが与えられなかった。
	}

	return setting, settingsfordd, files, &argsAll
}

// 即終了系の処理 OpenAIに投げないで終了する処理
func RunImmidiateTerminate(argsAll *ArgsAll, readme *string) {

	// 引数--readmeがあれば、詳細な説明文を出力して終了する。
	if argsAll.Readme {
		fmt.Println(*readme)
		os.Exit(0)
	}

	// 引数で、設定ファイルの雛形を生成する指示があれば生成して終了。
	if argsAll.CreateSetting {
		CreateSettingHjsonTemplate("./easygpt.hjson")
		fmt.Printf("設定ファイルの雛形を以下に出力しました。\n./easygpt.hjson\n出力した設定ファイルの雛形は編集が必要です。\nAPIキーやプロンプトを記載して使用してください。\n")
		os.Exit(0)
	}
}

func CreateUsageText(parser *arg.Parser) *string {
	outputs := make([]string, 0, 100)
	buf := new(bytes.Buffer)
	//parser.WriteUsage(os.Stdout) // キャプチャしているつもり。
	parser.WriteUsage(buf)
	strs := strings.Split(buf.String(), " ")
	for i, s := range strs {
		if i == 0 {
			outputs = append(outputs, s, "\n")
		} else if i == 1 {
			outputs = append(outputs, fmt.Sprintf(`  $ %v [options...] [files...]

  * 引数は全て順不同。混在可。
  * 引数による設定ではなく、設定ファイルによる指定を推奨する。
  * 設定ファイルは、'--create-setting'によって雛形を生成可能。
  * APIキーは環境変数'OPENAI_API_KEY'に設定する事を推奨。設定ファイルに指定してもよい。
  * 詳しくは '--help' や '--readme' を参照。

OPTIONS:
`, s))
		} else {
			if s[0] == '[' {
				if s[1] == '-' {
					outputs = append(outputs, fmt.Sprintf("    %v ", s))
				} else {
					// これ以降は直書き...orz
					outputs = append(outputs, fmt.Sprintf(`    [--help]
FILES:
    [任意個の入力ファイル、入力ファイルディレクトリ...]
`))
					break
				}
				if s[len(s)-1] == ']' {
					outputs = append(outputs, "\n")
				}
			} else if s[len(s)-1] == ']' {
				outputs = append(outputs, s, "\n")
			} else {
				outputs = append(outputs, " ", s)
			}
		}
	}
	outputs = append(outputs, "\n")
	return StringJoin(outputs)
}
