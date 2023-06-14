package internal

import (
	"log"
	"strings"

	"github.com/hjson/hjson-go"
	"github.com/pkg/errors"
)

func CreateSettingHjsonTemplate(paths []string) {
	setting := Setting{
		Apikey:      "ここに発行したAPIキーを記述する",
		InputDir:    "./入力ファイルがあるディレクトリのパス",
		OutputDir:   "./出力先ディレクトリのパス",
		Prompt:      "",
		Extension:   "",
		Postfix:     "",
		Concurrency: 1,
		Tmp:         "",
		OpenaiURL:   "https://api.openai.com/v1/chat/completions",
		AiModel:     "gpt-3.5-turbo-16k",
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

	for _, path := range paths {
		OutputTextForCheck(path, h)
	}
}
