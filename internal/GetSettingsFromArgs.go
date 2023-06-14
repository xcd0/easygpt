package internal

import (
	"fmt"

	"github.com/pkg/errors"
)

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

	if _, errApikey := GetApikey(&args.Apikey, &args.ApiFile); errApikey != nil {
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
		return nil, errInputDir // これは続行不可。
	}
	if errOutputDir := GetOutputDir(&args.InputDir); errOutputDir != nil {
		// 出力ディレクトリ指定がない
		//fmt.Println(errOutputDir)
		//os.Exit(1) // これは終了。
		return nil, errOutputDir
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
