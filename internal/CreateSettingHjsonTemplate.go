package internal

import (
	"fmt"
	"log"
	"strings"

	"github.com/hjson/hjson-go/v4"
	"github.com/pkg/errors"
)

func GetTemplateSetting() *Setting {
	return &Setting{
		InputDir:  "",
		OutputDir: "",
		SettingCommon: SettingCommon{
			Apikey:      "",
			Prompt:      "",
			Extension:   "",
			Postfix:     "",
			Concurrency: 1,
			Tmp:         "",
			Temperature: 0.7,
			OpenaiURL:   "https://api.openai.com/v1/chat/completions",
			AiModel:     "gpt-3.5-turbo-16k",
		},
	}
}

func CreateSettingHjsonTemplate(path string) {
	setting := GetTemplateSetting()
	b, err := hjson.Marshal(*setting)
	if err != nil {
		log.Printf("%+v", errors.Errorf("%v", err))
	}
	//log.Printf("b:\n%v", string(b))
	hj := string(b)
	//log.Printf("h:\n%v", h)
	hj = hj[1:]                           // hjsonの{}を削除
	hj = hj[:len(hj)-2] + "\n\n"          // hjsonの{}を削除
	hj = strings.ReplaceAll(hj, "  ", "") // hjsonのインデントを削除
	discription :=
		`# easygptの設定ファイル
#
# これはhjson形式のテキストファイルです。hjsonについては https://hjson.github.io/ を参照してください。。
# hjsonはjsonやyamlやtomlよりも人が編集しやすい形式であるので採用しています。hjsonの読み方は "aitch-jason" だそうです。
`
	hj = fmt.Sprintf("%v%v", discription, hj)
	//log.Printf("h:\n%v", h)

	OutputTextForCheck(&path, &hj)
}
