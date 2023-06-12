package internal

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/hjson/hjson-go/v4"
	"github.com/pkg/errors"
)

type ArgsSetting struct {
	CreateSetting string `arg:"--create-setting"   help:"設定ファイルの雛形を生成する。\n"`
	Setting       string `arg:"--setting"          help:"設定ファイルを指定する。\n                         指定がなければカレントディレクトリか、ホームディレクトリか、実行ファイルのあるディレクトリを探す。\n                         探すファイル名は'easygpt.hjson'、'.easygpt.hjson'、'.easygpt'の3つ。\n"`
}

func GetSettingFilePath() (string, error) {
	dirs := []string{GetCurrentDir(), GetHomeDir(), GetBinDir()}
	settingFileNames := []string{
		"easygpt.hjson",
		".easygpt.hjson",
		".easygpt",
	}
	// 設定ファイルは、カレントディレクトリ -> ホームディレクトリ -> 実行ファイルがあるディレクトリ
	// の順で調べる。
	for _, d := range dirs {
		for _, n := range settingFileNames {
			p := filepath.Join(d, n)
			//log.Printf("Debug: check path : %v", p)
			if _, err := os.Stat(p); err == nil {
				// 設定ファイルがあった。
				return p, nil
			}
		}
	}
	err := errors.Errorf("設定ファイルが見つかりませんでした。")
	//log.Printf("Debug: %v", err)
	return "", err
}

func GetSettingFilePathFromArgs(argsAll *ArgsAll) string {

	//log.Printf("--------------------------------------")
	argsSetting := ArgsSetting{
		CreateSetting: argsAll.CreateSetting,
		Setting:       argsAll.Setting,
	}

	//log.Printf("argsSetting:%v", argsSetting)
	//arg.MustParse(&argsSetting)
	//log.Printf("--------------------------------------")
	if len(argsSetting.CreateSetting) > 0 {
		// 雛形生成
		CreateSettingHjsonTemplate(argsSetting.CreateSetting)
		return ""
	} else if len(argsSetting.Setting) > 0 {
		return argsSetting.Setting
	}
	return ""
}

func CreateSettingHjsonTemplate(path string) {
	setting := Setting{
		Apikey:      "ここに発行したAPIキーを記述する",
		InputDir:    "./入力ファイルがあるディレクトリのパス",
		OutputDir:   "./出力先ディレクトリのパス",
		Prompt:      "",
		Postfix:     "",
		Extension:   "",
		Tmp:         "",
		Concurrency: 1,
	}
	b, err := hjson.Marshal(setting)
	if err != nil {
		log.Printf("%+v", errors.Errorf("%v", err))
	}
	//log.Printf("b:\n%v", string(b))
	h := string(b)
	//log.Printf("h:\n%v", h)
	h = h[1:]                           // hjsonの{}を削除
	h = h[:len(h)-2] + "\n\n"           // hjsonの{}を削除
	h = strings.ReplaceAll(h, "  ", "") // hjsonのインデントを削除
	//log.Printf("h:\n%v", h)

	OutputTextForCheck(path, h)
}
