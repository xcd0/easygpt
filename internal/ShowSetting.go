package internal

import (
	"fmt"
	"strings"
)

func ShowSetting(s *Setting) {
	p := strings.ReplaceAll(s.SettingCommon.Prompt, "\n", "\\n")
	if len(p) > 30 {
		p = p[:30]
	}
	fmt.Printf("------------------------------------------------\n")
	fmt.Printf("InputDir    : '%v'\n", s.InputDir)
	fmt.Printf("OutputDir   : '%v'\n", s.OutputDir)
	fmt.Printf("Apikey      : '%v'\n", s.SettingCommon.Apikey)
	fmt.Printf("Prompt      : '%v'\n", p)
	fmt.Printf("Extension   : '%v'\n", s.SettingCommon.Extension)
	fmt.Printf("Postfix     : '%v'\n", s.SettingCommon.Postfix)
	fmt.Printf("Concurrency : '%v'\n", s.SettingCommon.Concurrency)
	fmt.Printf("Tmp         : '%v'\n", s.SettingCommon.Tmp)
	fmt.Printf("Temperature : '%v'\n", s.SettingCommon.Temperature)
	fmt.Printf("OpenaiURL   : '%v'\n", s.SettingCommon.OpenaiURL)
	fmt.Printf("AiModel     : '%v'\n", s.SettingCommon.AiModel)
	fmt.Printf("Split       : '%v'\n", s.SettingCommon.Split)
	fmt.Printf("------------------------------------------------\n")
}
