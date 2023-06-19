package internal

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
			Split:       0,
		},
	}
}
