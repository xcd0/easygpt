package internal

import (
	"fmt"

	"github.com/pkg/errors"
)

type ArgsCommandLine struct {
	InputDir    string `arg:"--input-dir"        help:"入力テキストファイルのあるディレクトリパス。\n                         再帰的にファイルを検索して全て処理する。\n"`
	OutputDir   string `arg:"--output-dir"       help:"出力テキストファイルのあるディレクトリパス。\n                         input-dirのディレクトリと同じ構造でサブディレクトリを作成する。\n"`
	ApiFile     string `arg:"--apikey-file"      help:"テキストファイルによってAPIキーを指定する。\n                         この指定がない時、カレントディレクトリの./apikey.txtがあればそれを使用する。\n"`
	Apikey      string `arg:"--apikey"           help:"APIキーを直接指定する。\n                         この指定がある時、--api-fileによる指定を無視する。\n"`
	PromptFile  string `arg:"--prompt-file"      help:"全ての入力テキストファイルの直前に付与したい文字列を書いたテキストファイルパス。\n                         例えば、「以下のテキストを和訳してください。」と書いたテキストファイルを指定して、\n                         入力テキストファイルとして英文テキストファイルを与えれば、翻訳して貰える。\n                         この指定がない時、カレントディレクトリの./prompt.txtがあればそれを使用する。\n"`
	Prompt      string `arg:"--prompt"           help:"全ての入力テキストファイルの直前に付与したい文字列を指定する。\n                         この指定があるときファイルによる指定を無視する。\n                         プロンプトが指定されていない時、入力テキストファイルの内容がそのまま使用される。\n"`
	Postfix     string `arg:"--postfix"          help:"出力ファイル名の末尾に付与する文字列。\n                         空の時 _easygpt_output となる。\n"`
	Extension   string `arg:"--extension"        help:"入力として使用したいテキストファイルの拡張子。\n                         指定なしの時すべて使用する。\n                         拡張子のドットを含めて.mdのように指定する。\n"`
	Tmpdir      string `arg:"--tmp-dir"          help:"一時ファイルを保存するディレクトリを指定する。\n                         指定がない時、カレントディレクトリにtmpディレクトリを作成する。\n                         既にあれば、削除して再作成する。\n"`
	Concurrency int    `arg:"--concurrency"      help:"並列処理数を指定する。初期値1。\n                         APIの Rate Limitに引っかからない程度に並列したいところ。しかし、それは入力ファイル次第。\n                         この数値は単純に並行処理のスレッド数だと思ったらよい。\n                         Token/分とRequest/分に配慮する。\n                         小さなファイルは並列数を小さめ、大きなファイルは少し大きく、という感じだと思われる。\n"`
}

func (ArgsCommandLine) Description() string { // {{{
	return `# easygpt

chatgptのapiを使ってテキストファイルをまとめて一括で処理させるアプリ。
翻訳や、ソースコードにコメントを付けさせたりと、使い方次第で色々できる。

## 使い方1 D&D
1. 実行ファイルと同じディレクトリに2つテキストファイルを作成する。
	* "./apikey.txt"
		* APIキーを書き込む。  
		"echo "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" > ./apikey.txt"
		* APIキーは https://platform.openai.com/account/api-keys から発行できる。
	* "./prompt.txt"
		* これはなくてもよい。
		* 入力テキストファイルの前に与えたい文字列を記載する。
			* 例) 英文テキストファイルを翻訳してほしい場合、"./prompt.txt"に"以下を和訳してください。"と書く。

1. gptに投げたいテキストファイル、またはそれが含まれるディレクトリを"easygpt.exe"の実行ファイルにドラッグアンドドロップする。
1. 投げたファイルと同じディレクトリに、入力ファイルに"_easygpt_output"を付与した名前で処理結果を出力する。

## 使い方2 コマンドラインから実行

詳しくは "easygpt -h" を参照

### 例

./inputに英文テキストファイルがあるとして、それらをまとめて和訳させる例。

easygpt --input-dir ./input --output-dir ./output --api-key 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 --prompt 以下を和訳してください。

## ヘルプ
`
} //}}}

func GetSettingsFromArgs(argsAll *ArgsAll, settings []Setting) ([]Setting, error) {
	var args ArgsCommandLine
	if len(settings) == 0 {
		args = ArgsCommandLine{
			InputDir:    argsAll.InputDir,
			OutputDir:   argsAll.OutputDir,
			ApiFile:     argsAll.ApiFile,
			Apikey:      argsAll.Apikey,
			PromptFile:  argsAll.PromptFile,
			Prompt:      argsAll.Prompt,
			Postfix:     argsAll.Postfix,
			Extension:   argsAll.Extension,
			Tmpdir:      argsAll.Tmpdir,
			Concurrency: argsAll.Concurrency,
		}
	} else {
		args = ArgsCommandLine{
			InputDir:    settings[0].InputDir,
			OutputDir:   settings[0].OutputDir,
			ApiFile:     argsAll.ApiFile,
			Apikey:      settings[0].Apikey,
			PromptFile:  argsAll.PromptFile,
			Prompt:      settings[0].Prompt,
			Postfix:     settings[0].Postfix,
			Extension:   settings[0].Extension,
			Tmpdir:      settings[0].Tmp,
			Concurrency: settings[0].Concurrency,
		}
	}

	//arg.MustParse(&args)

	if errApikey := GetApikey(&args.Apikey, &args.ApiFile); errApikey != nil {
		//log.Fatalf("%v", errApikey) // これは続行不可。
		return nil, errors.Errorf("%v", errApikey) // これは続行不可。
	}
	if errPrompt := GetPrompt(&args.Prompt, &args.PromptFile); errPrompt != nil {
		fmt.Println(errPrompt) // これは続行。
	}
	if GetPostfix(&args.Postfix); len(args.Postfix) != 0 {
		fmt.Printf("出力ファイルのファイル名の末尾に %v を付与します。", args.Postfix)
	}

	// 引数に非フラグ引数でファイルが与えられていなかった。
	if errInputDir := GetInputDir(&args.InputDir); errInputDir != nil {
		// 入力ディレクトリ指定がない
		//fmt.Println(errInputDir)
		//os.Exit(1) // これは終了。
		return nil, errors.Errorf("%v", errInputDir) // これは続行不可。
	}
	if errOutputDir := GetOutputDir(&args.InputDir); errOutputDir != nil {
		// 出力ディレクトリ指定がない
		//fmt.Println(errOutputDir)
		//os.Exit(1) // これは終了。
		return nil, errors.Errorf("%v", errOutputDir)
	}

	return []Setting{Setting{
		Apikey:      args.Apikey,
		InputDir:    args.InputDir,
		OutputDir:   args.OutputDir,
		Prompt:      args.Prompt,
		Postfix:     args.Postfix,
		Extension:   args.Extension, // 指定がなければそれはそれでOK
		Tmp:         args.Tmpdir,
		Concurrency: args.Concurrency,
	}}, nil
}
