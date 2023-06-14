package internal

import (
	"fmt"
	"path/filepath"
)

func RunNoSetting(settings []SettingForDD) {

	for _, s := range settings {
		outputfilename := AddPostfix(&s.Input, &s.Postfix)
		output := Question(&s.OpenaiURL, &s.AiModel, &s.ApiKey, &s.Prompt, &s.Input, outputfilename, s.Temperature, false)
		outputpath := filepath.Join(filepath.Dir(s.Input), *outputfilename)
		OutputTextForCheck(&outputpath, output)
		//log.Printf("============================")
	}
}

func AddPostfix(file, postfix *string) *string {

	// これはすでに最低限必要な設定が得られているときにのみ生成する。最低限の設定が得られていないときnilになる。
	// postfixが空の時、入力ファイルに上書きしてしまうので勝手に文字列を付与する。
	if postfix == nil {
		return file
	}
	if len(*postfix) == 0 {
		*postfix = "_easygpt_output"
	}
	// どうしても上書きしたいときは `--postfix *` のように設定する。
	if *postfix == "*" {
		*postfix = ""
	}
	ext := filepath.Ext(*file)
	name := GetFileNameWithoutExt(file)
	n := fmt.Sprintf("%v%v%v", name, *postfix, ext)
	return &n
}
