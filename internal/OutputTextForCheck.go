package internal

import (
	"log"
	"os"
)

func OutputTextForCheck(file path, str string) {
	f, err := os.Create(file)
	if err != nil {
		log.Printf("Warning: 確認用の一時ファイルが作成できませんでした。動作に支障はないので続行します。\n%v", err)
		//panic(err)
	} else {
		if _, err := f.Write([]byte(str)); err != nil {
			log.Printf("Warning: 確認用の一時ファイルに書き込めませんでした。動作に支障はないので続行します。\n%v", err)
		}
	}
}
