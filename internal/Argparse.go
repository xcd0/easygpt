package internal

import (
	"fmt"
	"log"
	"os"

	"github.com/alexflint/go-arg"
)

type ArgsAll struct {
	ArgsSetting
	ArgsDD
	ArgsCommandLine

	InputText string `arg:"-i,--input-text" help:"質問したい文字列を直接指定。\n                         この引数がある場合、APIキー指定、OpenAIのURL指定、AIのモデル指定、AIのTemperature、\n                         これら以外の引数と設定を無視します。\n"`
	Readme    bool   `arg:"-r,--readme" help:"詳しい説明文を出力する。長いため\"./easygpt --readme | less\"などで見るのがおすすめ。\n"`
}

type ArgsSetting struct {
	CreateSetting bool   `arg:"-c,--create-setting"   help:"設定ファイルの雛形を生成する。\n"`
	Setting       string `arg:"-s,--setting"          help:"設定ファイルを指定する。\n                         指定がなければカレントディレクトリか、ホームディレクトリか、実行ファイルのあるディレクトリを探す。\n                         探すファイル名は'easygpt.hjson'、'.easygpt.hjson'、'.easygpt'の3つ。\n"`
}

type ArgsDD struct {
	InputFiles []string `arg:"positional"         help:"非フラグ引数はファイルやディレクトリのパスであると見なす。\n                         この方法で指定されたファイル群は、引数の--output-dirの指定を無視して、\n                         指定されたファイルと同じディレクトリに、POSTFIXを付与した名前で出力される。"`
}

type ArgsCommandLine struct {
	InputDir    string  `arg:"--input-dir"        help:"入力テキストファイルのあるディレクトリパス。\n                         再帰的にファイルを検索して全て処理する。\n"`
	OutputDir   string  `arg:"--output-dir"       help:"出力テキストファイルのあるディレクトリパス。\n                         input-dirのディレクトリと同じ構造でサブディレクトリを作成する。\n"`
	ApiFile     string  `arg:"--apikey-file"      help:"テキストファイルによってAPIキーを指定する。\n                         この指定がない時、カレントディレクトリの./apikey.txtがあればそれを使用する。\n"`
	Apikey      string  `arg:"-a,--apikey"        help:"APIキーを直接指定する。\n                         この指定がある時、--api-fileによる指定を無視する。\n"`
	PromptFile  string  `arg:"--prompt-file"      help:"全ての入力テキストファイルの直前に付与したい文字列を書いたテキストファイルパス。\n                         例えば、「以下のテキストを和訳してください。」と書いたテキストファイルを指定して、\n                         入力テキストファイルとして英文テキストファイルを与えれば、翻訳して貰える。\n                         この指定がない時、カレントディレクトリの./prompt.txtがあればそれを使用する。\n"`
	Prompt      string  `arg:"-p,--prompt"        help:"全ての入力テキストファイルの直前に付与したい文字列を指定する。\n                         この指定があるときファイルによる指定を無視する。\n                         プロンプトが指定されていない時、入力テキストファイルの内容がそのまま使用される。\n"`
	Postfix     string  `arg:"--postfix"          help:"出力ファイル名の末尾に付与する文字列。\n                         空の時 _easygpt_output となる。\n"`
	Extension   string  `arg:"--extension"        help:"入力として使用したいテキストファイルの拡張子。\n                         指定なしの時すべて使用する。\n                         拡張子のドットを含めて.mdのように指定する。\n"`
	Tmpdir      string  `arg:"--tmp-dir"          help:"一時ファイルを保存するディレクトリを指定する。\n                         指定がない時、カレントディレクトリにtmpディレクトリを作成する。\n                         既にあれば、削除して再作成する。\n"`
	Concurrency int     `arg:"--concurrency"      help:"並列処理数を指定する。初期値1。\n                         APIの Rate Limitに引っかからない程度に並列したいところ。しかし、それは入力ファイル次第。\n                         この数値は単純に並行処理のスレッド数だと思ったらよい。\n                         Token/分とRequest/分に配慮する。\n                         小さなファイルは並列数を小さめ、大きなファイルは少し大きく、という感じだと思われる。\n"`
	Temperature float64 `args:"-t,--temperature"  help:"これは、AIに与える変数で、返答のランダム性を制御するパラメータである。\n                         値が小さいほどよくある解答など決まりきった解答を返し、\n                         値が大きいほど奇抜な返答が返ってきやすくなる。\n                         0から2の範囲の値を設定する。詳細は\n                         https://platform.openai.com/docs/api-reference/completions/create#completions/create-temperature\n                         を参照。\n"`
}

// 戻り値は引数処理結果のsettingと引数にあったファイルのリスト
func Argparse() ([]Setting, []SettingForDD) {
	var settings []Setting
	var settingsfordd []SettingForDD
	var argsAll ArgsAll = ArgsAll{}
	var err error

	parser := arg.MustParse(&argsAll)

	//log.Printf("argsAll:%v", argsAll)
	//log.Printf("--------------------------------------")

	if argsAll.Readme {
		ShowDescription()
		os.Exit(0)
	}

	// 引数解析は段階を踏む
	// ①設定ファイルからの読み込み
	// ②引数からの設定の読み込み
	// ③引数からのファイルの読み込み

	// まず設定ファイルの指定があるかどうか

	// 引数で、設定ファイルの雛形を生成する指示があるか。あれば生成して終了。
	if argsAll.CreateSetting {
		// 雛形生成
		CreateSettingHjsonTemplate("./easygpt.hjson")
		fmt.Println("設定ファイルの雛形を以下に出力しました。")
		fmt.Println("./easygpt.hjson")
		fmt.Println("出力した設定ファイルの雛形は編集が必要です。")
		fmt.Println("APIキーやプロンプトを記載して使用してください。")
		os.Exit(0)
	}

	// 引数で設定ファイル指定があるか。
	if len(argsAll.Setting) > 0 {
		// 引数で設定ファイルが指定された
		settings = ReadSettingHjson(&argsAll.Setting)
		//log.Printf("Debug: 設定ファイルが見つかりました。 : %v", argsAll.Setting)
		//log.Printf("--------------------------------------")
	} else if argsAll.Setting, err = GetSettingFilePath(); err != nil {
		log.Printf("%v", err)
	} else if len(argsAll.Setting) > 0 {
		// 既定の設定ファイルがあった。
		settings = ReadSettingHjson(&argsAll.Setting)
		//log.Printf("Debug: 設定ファイルが見つかりました。 : %v", argsAll.Setting)
	} else {
		//log.Printf("debug")
		// 設定ファイルがなかった。
	}

	// 引数で質問文が与えられている場合、APIキー以外の設定を無視する。
	if len(argsAll.InputText) != 0 {
		//log.Printf("settings:%v", settings)
		if len(argsAll.Apikey) == 0 {
			argsAll.Apikey = settings[0].Apikey
		}
		//log.Printf("&argsAll.Apikey:%v, &argsAll.ApiFile:%v", argsAll.Apikey, argsAll.ApiFile)
		if _, errApikey := GetApikey(&argsAll.Apikey, &argsAll.ApiFile); errApikey != nil {
			log.Printf("%v", errApikey)
			os.Exit(1)
		}
		//log.Printf("&argsAll.Apikey:%v, &argsAll.ApiFile:%v", argsAll.Apikey, argsAll.ApiFile)
		s := Setting{
			Apikey:      argsAll.Apikey,
			InputDir:    "",
			OutputDir:   "",
			Prompt:      "",
			Postfix:     "",
			Extension:   "",
			Tmp:         "",
			Concurrency: 1,
			AiModel:     settings[0].AiModel,
			Temperature: settings[0].Temperature,
			OpenaiURL:   settings[0].OpenaiURL,
		}
		if argsAll.Temperature != 0 { // 初期値でない
			s.Temperature = argsAll.Temperature
		}
		s.Temperature = Clamp(s.Temperature, 0, 2)

		output := QuestionByText(&s.OpenaiURL, &s.AiModel, &s.Apikey, &s.Prompt, &argsAll.InputText, &s.Tmp, s.Temperature, false)
		if output == nil {
			fmt.Fprintln(os.Stderr, "エラー終了")
			os.Exit(1)
		}

		fmt.Println(*output)
		os.Exit(0)
	}

	//log.Printf("--------------------------------------")
	if true { // 設定ファイルがあっても、引数で指定された場合読み込んだ設定を上書きする。
		if ret, err := GetSettingsFromArgs(&argsAll, settings); err == nil { // 引数から設定を読み込む。
			//log.Printf("debug")
			if len(settings) == 0 {
				settings = ret
			} else {
				if len(ret[0].Apikey) != 0 {
					settings[0].Apikey = ret[0].Apikey // 上書き
				}
				if len(ret[0].InputDir) != 0 {
					settings[0].InputDir = ret[0].InputDir // 上書き
				}
				if len(ret[0].OutputDir) != 0 {
					settings[0].OutputDir = ret[0].OutputDir // 上書き
				}
				if len(ret[0].Prompt) != 0 {
					settings[0].Prompt = ret[0].Prompt // 上書き
				}
				if len(ret[0].Postfix) != 0 {
					settings[0].Postfix = ret[0].Postfix // 上書き
				}
				if len(ret[0].Extension) != 0 {
					settings[0].Extension = ret[0].Extension // 上書き
				}
				if len(ret[0].Tmp) != 0 {
					settings[0].Tmp = ret[0].Tmp // 上書き
				}
				if ret[0].Concurrency != 1 {
					settings[0].Concurrency = ret[0].Concurrency // 上書き
				}
			}
		} else {
			//log.Printf("%+v", err) // 引数からの設定に失敗。
			if len(os.Args) == 1 {
				// 引数がない
				parser.WriteUsage(os.Stdout)
				os.Exit(1)

				os.Exit(1)
			}
		}
	}

	// 引数に与えられた入力ファイル毎に設定ファイルを生成する。
	settingsfordd, files := GetSettingsForDD(&argsAll, settings)
	if files == nil {
		//引数ファイルが与えられなかった。
		return settings, nil
	} else {
		return settings, settingsfordd
	}

}
