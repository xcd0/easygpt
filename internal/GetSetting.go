package internal

import (
	"fmt"
	"os"

	"github.com/alexflint/go-arg"
)

// もし、初期値が設定されている値が、設定ファイルにおいて指定されていない場合、初期値を指定しなおす。
func ResetValueIfValueIsEmpty(setting *Setting) {
	s := GetTemplateSetting()
	if setting.Temperature == 0 {
		setting.Temperature = s.Temperature
	}
	if len(setting.OpenaiURL) == 0 {
		setting.OpenaiURL = s.OpenaiURL
	}
	if len(setting.AiModel) == 0 {
		setting.AiModel = s.AiModel
	}
}

func HasSetting() bool {
	s, err := GetSettingFilePath()
	return err == nil && len(*s) > 0
}

// 引数と設定ファイルから設定を取得する。
func GetSetting(argsAll *ArgsAll, parser *arg.Parser) (*Setting, error) {
	setting := GetTemplateSetting() // 最低限初期値を設定しておく

	// 設定ファイルによる指定。
	if s, err := GetSettingFilePath(); err == nil && len(*s) > 0 { // 既定の設定ファイルがあった。
		setting = ReadSettingHjson(s)
	}
	if len(argsAll.Setting) > 0 { // 引数で設定ファイルが指定された
		setting = ReadSettingHjson(&argsAll.Setting)
	}
	// もし、初期値が設定されている値が、設定ファイルにおいて指定されていない場合、初期値を指定しなおす。
	ResetValueIfValueIsEmpty(setting)

	//ShowSetting(setting)
	//os.Exit(0)

	// 引数による指定。設定ファイルによる設定を上書きする。
	args := &argsAll.ArgsCommandLine
	if len(args.Prompt) != 0 {
		setting.Prompt = args.Prompt
	} else if errPrompt := GetPrompt(&args.Prompt, &args.PromptFile); errPrompt != nil {
		fmt.Println(errPrompt) // これは続行。
	}
	if len(args.Postfix) != 0 {
		setting.Postfix = args.Postfix
	}
	if len(args.Extension) != 0 {
		setting.Extension = args.Extension
	}
	if len(args.Tmp) != 0 {
		setting.Tmp = args.Tmp
	}
	if args.Concurrency != 0 {
		setting.Concurrency = args.Concurrency
	}
	if len(args.AiModel) != 0 {
		setting.AiModel = args.AiModel
	}
	if args.Temperature != 0 { // args.Temperatureの初期値は0のはず setting.Temperatureの初期値は0.7
		setting.Temperature = args.Temperature
	}
	if len(args.OpenaiURL) != 0 {
		setting.OpenaiURL = args.OpenaiURL
	}
	if args.Split != 0 {
		setting.Split = args.Split
	}
	if len(args.Move) != 0 {
		setting.Move = args.Move
	}

	if GetPostfix(&args.Postfix); len(args.Postfix) != 0 {
		fmt.Printf("出力ファイルのファイル名の末尾に %v を付与します。", args.Postfix)
	}

	if len(args.Apikey) != 0 {
		setting.Apikey = args.Apikey
	} else {
		// 環境変数OPENAI_API_KEYが設定されていたら読み込む。
		if apikey := os.Getenv("OPENAI_API_KEY"); len(apikey) != 0 {
			setting.Apikey = apikey
		} else
		// ファイルからの取得
		if ret, errApikey := GetApikey(&args.Apikey, &args.ApiFile); errApikey != nil {
			if len(setting.Apikey) == 0 { // 設定ファイルで既に指定されていなければ
				// 引数からも、環境変数からも、設定ファイルからも、APIキーが取得できなかった。
				// APIキーがないと何もできないので終了する。
				return setting, errApikey // 入力ディレクトリ指定がないこれは続行不可。
			}
		} else {
			setting.InputDir = *ret // 指定されていたので設定ファイルの設定を上書きする。
		}
	}
	if len(args.InputDir) != 0 {
		setting.InputDir = args.InputDir // 引数で指定されていたので設定ファイルの設定を上書きする。
	}
	if len(args.OutputDir) != 0 {
		setting.OutputDir = args.OutputDir // 引数で指定されていたので設定ファイルの設定を上書きする。
	}
	return setting, nil
}
