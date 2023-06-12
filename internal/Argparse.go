package internal

import (
	"log"

	"github.com/alexflint/go-arg"
)

type ArgsAll struct {
	ArgsSetting //
	/* type ArgsSetting struct {
		CreateSetting string `arg:"--create-setting"   help:"設定ファイルの雛形を生成する。\n"`
		Setting       string `arg:"--setting"          help:"設定ファイルを指定する。\n                         指定がなければカレントディレクトリか、ホームディレクトリか、実行ファイルのあるディレクトリを探す。\n                         探すファイル名は'easygpt.hjson'、'.easygpt.hjson'、'.easygpt'の3つ。\n"`
	}*/
	ArgsDD

	/*type ArgsDD struct {
		InputFiles []string `arg:"positional"         help:"非フラグ引数はファイルやディレクトリのパスであると見なす。\n                         この方法で指定されたファイル群は、引数の--output-dirの指定を無視して、\n                         指定されたファイルと同じディレクトリに、POSTFIXを付与した名前で出力される。"`
	}*/
	ArgsCommandLine
	/* type ArgsCommandLine struct {
		InputDir    string `arg:"--input-dir"        help:"入力テキストファイルのあるディレクトリパス。\n                         再帰的にファイルを検索して全て処理する。\n"`
		OutputDir   string `arg:"--output-dir"       help:"出力テキストファイルのあるディレクトリパス。\n                         input-dirのディレクトリと同じ構造でサブディレクトリを作成する。\n"`
		ApiFile     string `arg:"--apikey-file"      help:"テキストファイルによってAPIキーを指定する。\n                         この指定がない時、カレントディレクトリの./apikey.txtがあればそれを使用する。\n"`
		Apikey      string `arg:"--apikey"           help:"APIキーを直接指定する。\n                         この指定がある時、--api-fileによる指定を無視する。\n"`
		PromptFile  string `arg:"--prompt-file"      help:"全ての入力テキストファイルの直前に付与したい文字列を書いたテキストファイルパス。\n                         例えば、「以下のテキストを和訳してください。」と書いたテキストファイルを指定して、\n                         入力テキストファイルとして英文テキストファイルを与えれば、翻訳して貰える。\n                         この指定がない時、カレントディレクトリの./prompt.txtがあればそれを使用する。\n"`
		Prompt      string `arg:"--prompt"           help:"全ての入力テキストファイルの直前に付与したい文字列を指定する。\n                         この指定があるときファイルによる指定を無視する。\n                         プロンプトが指定されていない時、入力テキストファイルの内容がそのまま使用される。\n"`
		Postfix     string `arg:"--postfix"          help:"出力ファイル名の末尾に付与する文字列。\n                         空の時 _easygpt_output となる。\n"`
		Extension   string `arg:"--extension" help:"入力として使用したいテキストファイルの拡張子。\n                         指定なしの時すべて使用する。\n                         拡張子のドットを含めて.mdのように指定する。\n"`
		Tmpdir      string `arg:"--tmp-dir"          help:"一時ファイルを保存するディレクトリを指定する。\n                         指定がない時、カレントディレクトリにtmpディレクトリを作成する。\n                         既にあれば、削除して再作成する。\n"`
		Concurrency int    `arg:"--concurrency"      help:"並列処理数を指定する。初期値1。\n                         APIの Rate Limitに引っかからない程度に並列したいところ。しかし、それは入力ファイル次第。\n                         この数値は単純に並行処理のスレッド数だと思ったらよい。\n                         Token/分とRequest/分に配慮する。\n                         小さなファイルは並列数を小さめ、大きなファイルは少し大きく、という感じだと思われる。\n"`
	} */
}

// 戻り値は引数処理結果のsettingと引数にあったファイルのリスト
func Argparse() ([]Setting, []SettingForDD) {
	var settings []Setting
	var settingsfordd []SettingForDD
	var argsAll ArgsAll

	//./easygpt --create-setting _createsetting --setting _setting positional1 positional2 --input-dir _inputdir --output-dir _outputdir --apikey-file _apikeyfile --apikey _apikey --prompt-file _promptfile --prompt _prompt --postfix _postfix --extension  _extension --tmp-dir _tmpdir --concurrency 1111
	arg.MustParse(&argsAll)
	//log.Printf("argsAll:%v", argsAll)
	//log.Printf("--------------------------------------")

	// 引数解析は段階を踏む
	// ①設定ファイルからの読み込み
	// ②引数からの設定の読み込み
	// ③引数からのファイルの読み込み

	// まず設定ファイルの指定があるかどうか

	if settingFilePath := GetSettingFilePathFromArgs(&argsAll); len(settingFilePath) > 0 {
		// 引数で設定ファイルが指定された
		settings = ReadSettingHjson(settingFilePath)
		//log.Printf("Debug: 設定ファイルが見つかりました。 : %v", settingFilePath)
		//log.Printf("--------------------------------------")
	} else if settingFilePath, err := GetSettingFilePath(); len(settingFilePath) > 0 {
		if err != nil {
			log.Printf("%v", err)
		} else {
			// 既定の設定ファイルがあった。
			settings = ReadSettingHjson(settingFilePath)
			//log.Printf("Debug: 設定ファイルが見つかりました。 : %v", settingFilePath)
		}
		//log.Printf("--------------------------------------")
	} else {
		//log.Printf("debug")
		// 設定ファイルがなかった。
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
			log.Printf("%+v", err) // 引数からの設定に失敗。
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
