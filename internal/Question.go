package internal

import (
	"fmt"
	"os"
	"path/filepath"
)

func QuestionByText(inputText *string, setting *Setting, tmpflag bool) *string {
	//log.Printf("temperature:%v", temperature)
	messages := []Message{}
	if len(setting.Prompt) > 0 {
		messages = append(messages,
			Message{
				Role:    "system",
				Content: setting.Prompt,
			})
	}
	if len(*inputText) > 0 {
		messages = append(messages,
			Message{
				Role:    "user",
				Content: *inputText,
			})
	}

	if len(messages) == 0 {
		fmt.Println("入力テキスト空でした。")
		os.Exit(0)
	}

	if tmpflag {
		p := filepath.Join(setting.Tmp, "input.txt")
		q := fmt.Sprintf("%v", messages)
		OutputTextForCheck(&p, &q)
	}
	//log.Printf("%v", messages)

	response := GetOpenAIResponse(&messages, &setting.OpenaiURL, &setting.AiModel, &setting.Apikey, &setting.Tmp, setting.Temperature, tmpflag)
	//log.Printf("response: %v", response)
	//log.Printf(": %v", response)
	if response == nil {
		return nil
	}

	if len((*response).Choices) == 0 {
		//log.Printf("Error: no response.")
		//log.Printf("       %v", response)
		return nil
	}

	output := (*response).Choices[0].Messages.Content
	if tmpflag {
		p := filepath.Join(setting.Tmp, "output.txt")
		OutputTextForCheck(&p, &output)
	}
	return &output
}
