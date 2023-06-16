package internal

import (
	"fmt"
	"log"
	"os"
)

func RunInputTextOnArgs(argsAll *ArgsAll, setting *Setting) {
	if len(argsAll.Apikey) == 0 { // 引数でapiキーが指定されておらず
		if len(setting.Apikey) != 0 { // 設定ファイルでapiキーが指定されていたら、
			argsAll.Apikey = setting.Apikey // 既に環境変数で設定されている可能性があるが、設定ファイルの値で上書きする。
		}
	}
	//log.Printf("&argsAll.Apikey:%v, &argsAll.ApiFile:%v", argsAll.Apikey, argsAll.ApiFile)
	if _, errApikey := GetApikey(&argsAll.Apikey, &argsAll.ApiFile); errApikey != nil {
		log.Printf("%v", errApikey)
		os.Exit(1)
	}
	//log.Printf("&argsAll.Apikey:%v, &argsAll.ApiFile:%v", argsAll.Apikey, argsAll.ApiFile)

	s := GetTemplateSetting()
	s.Apikey = argsAll.Apikey
	s.AiModel = setting.AiModel
	s.Temperature = setting.Temperature
	s.OpenaiURL = setting.OpenaiURL

	if argsAll.Temperature != 0 { // 初期値でない
		s.Temperature = argsAll.Temperature
	}
	s.Temperature = Clamp(s.Temperature, 0, 2)

	output := QuestionByText(&argsAll.InputText, s, false)
	if output == nil {
		fmt.Fprintln(os.Stderr, "エラー終了")
		os.Exit(1)
	}

	fmt.Println(*output)
}
