package internal

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/alexflint/go-arg"
)

type Args struct {
	InputFiles  []string `arg:"positional"         help:"非フラグ引数はファイルやディレクトリのパスであると見なす。\n                         この方法で指定されたファイル群は、引数の--output-dirの指定を無視して、\n                         指定されたファイルと同じディレクトリに、POSTFIXを付与した名前で出力される。"`
	InputDir    string   `arg:"--input-dir"        help:"入力テキストファイルのあるディレクトリパス。\n                         再帰的にファイルを検索して全て処理する。\n"`
	OutputDir   string   `arg:"--output-dir"       help:"出力テキストファイルのあるディレクトリパス。\n                         input-dirのディレクトリと同じ構造でサブディレクトリを作成する。\n"`
	ApiFile     string   `arg:"--api-file"         help:"テキストファイルによってAPIキーを指定する。\n                         この指定がない時、カレントディレクトリの./apikey.txtがあればそれを使用する。\n"`
	Apikey      string   `arg:"--api-key"          help:"APIキーを直接指定する。\n                         この指定がある時、--api-fileによる指定を無視する。\n"`
	PromptFile  string   `arg:"--prompt-file"      help:"全ての入力テキストファイルの直前に付与したい文字列を書いたテキストファイルパス。\n                         例えば、「以下のテキストを和訳してください。」と書いたテキストファイルを指定して、\n                         入力テキストファイルとして英文テキストファイルを与えれば、翻訳して貰える。\n                         この指定がない時、カレントディレクトリの./prompt.txtがあればそれを使用する。\n"`
	Prompt      string   `arg:"--prompt"           help:"全ての入力テキストファイルの直前に付与したい文字列を指定する。\n                         この指定があるときファイルによる指定を無視する。\n                         プロンプトが指定されていない時、入力テキストファイルの内容がそのまま使用される。\n"`
	Postfix     string   `arg:"--postfix"          help:"出力ファイル名の末尾に付与する文字列。\n                         空の時 _easygpt_output となる。\n"`
	Extension   string   `arg:"--target-extension" help:"入力として使用したいテキストファイルの拡張子。\n                         指定なしの時すべて使用する。\n                         拡張子のドットを含めて.mdのように指定する。\n"`
	Tmpdir      string   `arg:"--tmp-dir"          help:"一時ファイルを保存するディレクトリを指定する。\n                         指定がない時、カレントディレクトリにtmpディレクトリを作成する。\n                         既にあれば、削除して再作成する。\n"`
	Concurrency int      `arg:"--concurrency"      help:"並列処理数を指定する。初期値1。\n                         APIの Rate Limitに引っかからない程度に並列したいところ。しかし、それは入力ファイル次第。\n                         この数値は単純に並行処理のスレッド数だと思ったらよい。\n                         Token/分とRequest/分に配慮する。\n                         小さなファイルは並列数を小さめ、大きなファイルは少し大きく、という感じだと思われる。\n"`
}

func (Args) Description() string {
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

1. ヘルプを見ながら引数を設定する。

### 例

./inputに英文テキストファイルがあるとして、それらをまとめて和訳させる例。

easygpt --input-dir ./input --output-dir ./output --api-key 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 --prompt 以下を和訳してください。

## ヘルプ

`
}

// 戻り値は引数処理結果のsettingと引数にあったファイルのリスト
func Argparse() ([]Setting, []SettingForDD) {
	args := Args{
		InputFiles:  []string{},
		InputDir:    "",
		OutputDir:   "",
		ApiFile:     "",
		Apikey:      "",
		PromptFile:  "",
		Prompt:      "",
		Postfix:     "_easygpt_output",
		Extension:   "",
		Tmpdir:      "",
		Concurrency: 1,
	}

	arg.MustParse(&args)

	// D&Dで実行された場合の対応のために、
	// 引数に--xxxのような形式でない引数があれば、
	// それは個別に実行する。

	if errApikey := GetApikey(&args.Apikey, &args.ApiFile); errApikey != nil {
		log.Fatalf("%v", errApikey) // これは続行不可。
	}
	if errPrompt := GetPrompt(&args.Prompt, &args.PromptFile); errPrompt != nil {
		fmt.Println(errPrompt) // これは続行。
	}
	if GetPostfix(&args.Postfix); len(args.Postfix) != 0 {
		fmt.Printf("出力ファイルのファイル名の末尾に %v を付与します。", args.Postfix)
	}

	//log.Printf("%v", args.InputFiles)
	settingsfordd, files := GenerateSettingForDD(&args.Apikey, &args.Prompt, &args.Postfix, &args.InputFiles)

	if len(files) == 0 {
		// 引数に非フラグ引数でファイルが与えられていなかった。
		if errInputDir := GetInputDir(&args.InputDir); errInputDir != nil {
			// 入力ディレクトリ指定がない
			fmt.Println(errInputDir)
			os.Exit(1) // これは終了。
		}
		if errOutputDir := GetOutputDir(&args.InputDir); errOutputDir != nil {
			// 出力ディレクトリ指定がない
			fmt.Println(errOutputDir)
			os.Exit(1) // これは終了。
		}

		return []Setting{Setting{
			ApiKey:      args.Apikey,
			InputDir:    args.InputDir,
			OutputDir:   args.OutputDir,
			Prompt:      args.Prompt,
			Postfix:     args.Postfix,
			Extension:   args.Extension, // 指定がなければそれはそれでOK
			Tmp:         args.Tmpdir,
			Concurrency: args.Concurrency,
		}}, settingsfordd
	} else {
		return []Setting{}, settingsfordd
	}

}

func GenerateSettingForDD(apikey, prompt, postfix *string, inputArgFiles *[]string) ([]SettingForDD, []string) {
	// D&Dによる実行
	// D&Dの場合、引数を舐めて、ファイルはファイルごとに実行
	// ディレクトリは設定ありとある程度同じように実行

	// ファイル一覧を生成
	var files []string
	for _, p := range *inputArgFiles {
		p, _ = filepath.Abs(p)
		CheckPathAndRunFunction(p,
			func(err error) {
				// 何もしなければcontinue相当の動作になる。
				// 存在しない時無視する。
			},
			func() {
				fes, err := Dirwalk(p)
				if err != nil {
					fmt.Printf("Error: %v", err)
					fmt.Printf("       Skipped.")
				}
				files = append(files, fes.Files...)
			},
			func() {
				// ファイルが存在する
				files = append(files, p)
			},
		)
	}
	var settings []SettingForDD
	for _, f := range files {
		settings = append(settings,
			SettingForDD{
				ApiKey:      *apikey,  // APIキー
				Input:       f,        // 入力ファイル
				Prompt:      *prompt,  // プロンプト
				Postfix:     *postfix, // 出力ファイル名の末尾に付与する文字列。空の時 "_easygpt_output" となる。
				Tmp:         "",       // 一時ファイルを保存するディレクトリ
				Concurrency: 1,        // 最大同時並列実行数
			},
		)
	}

	return settings, files
}
