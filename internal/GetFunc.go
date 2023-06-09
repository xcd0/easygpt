package internal

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

func GetApikey(apikey, apiFile *string) error { // APIキー
	if len(*apikey) != 0 {
		// ok
	} else if ret, err := GetText(*apiFile); err == nil {
		// ok
		*apikey = ret
	} else if ret, err = GetText("./apikey.txt"); err == nil {
		// ok
		*apikey = ret
	} else {
		// APIキーが指定されておらず、APIキーを記載したテキストファイルのパスも与えられておらず、カレントディレクトリにもAPIキーを書いたテキストファイルがない
		return errors.Errorf("Error: APIキーを指定してください。")
	}
	*apikey = strings.ReplaceAll(*apikey, "\n", "") // 改行コードを削除
	*apikey = strings.ReplaceAll(*apikey, "\r", "") // 改行コードを削除
	return nil
}

func GetPrompt(prompt, promptFile *string) error { // prompt
	if len(*prompt) != 0 { // ok
	} else if ret, err := GetText(*promptFile); err == nil {
		*prompt = ret // ok
	} else if ret, err = GetText("./prompt.txt"); err == nil {
		*prompt = ret // ok
	} else { // プロンプトが与えられていない。
		return errors.Errorf("Warnig: プロンプトが与えられていません。\n        プロンプトがない場合は、入力テキストファイルがそのまま使用されます。")
	}
	return nil
}

func GetPostfix(postfix *string) error {
	if len(*postfix) != 0 {
		// ok
	} else if ret, err := GetText("./postfix.txt"); err == nil {
		// ok
		*postfix = ret
	} else {
		// postfixがない。別に問題ない。
	}
	return nil
}

func GetInputDir(inputDir *string) error { // 入出力テキストファイルディレクトリパス
	var err error
	if len(*inputDir) == 0 {
		return errors.Errorf("Error: 入出力ディレクトリを指定してください。")
	}
	if *inputDir, err = filepath.Abs(*inputDir); err != nil {
		return errors.Errorf("Error: 入力ディレクトリの指定が不正です。\n       指定されたディレクトリパス: %v", *inputDir)
	}
	return nil
}

func GetOutputDir(outputDir *string) error { // 入出力テキストファイルディレクトリパス
	var err error
	if len(*outputDir) == 0 {
		return errors.Errorf("Error: 入出力ディレクトリを指定してください。")
	}
	if *outputDir, err = filepath.Abs(*outputDir); err != nil {
		return errors.Errorf("Error: 出力ディレクトリの指定が不正です。\n       指定されたディレクトリパス: %v", *outputDir)
	}
	return nil
}

func GetText(filepath string) (string, error) {
	b, err := os.ReadFile(filepath) // https://pkg.go.dev/os@go1.20.5#ReadFile
	if err != nil {
		//log.Print("Error: %v, file: %v", err, filepath)
		return "", err
	}
	return string(b), err
}

func GetTextNoError(filepath string) string {
	defer func() {
		if err := recover(); err != nil {
			log.Fatalf("%+v", err)
		}
	}()

	b, err := os.ReadFile(filepath) // https://pkg.go.dev/os@go1.20.5#ReadFile
	if err != nil {
		panic(errors.Errorf("Error: %v, file: %v", err, filepath))
	}
	return string(b)
}

func GetFileNameWithoutExt(path string) string {
	return filepath.Base(path[:len(path)-len(filepath.Ext(path))])
}
