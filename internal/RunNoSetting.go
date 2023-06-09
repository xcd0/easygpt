package internal

import (
	"fmt"
	"path/filepath"
)

func RunNoSetting(settings []SettingForDD) {

	for _, s := range settings {
		outputfilename := AddPostfix(s.Input, s.Postfix)
		output := Question(s.ApiKey, &s.Prompt, s.Input, outputfilename, false)
		outputpath := filepath.Join(filepath.Dir(s.Input), outputfilename)
		OutputTextForCheck(outputpath, output)
	}
}

func AddPostfix(file, postfix string) string {
	ext := filepath.Ext(file)
	name := GetFileNameWithoutExt(file)
	return fmt.Sprintf("%v%v.%v", name, postfix, ext)
}
