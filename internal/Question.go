package internal

import (
	"fmt"
	"path/filepath"

	"github.com/pkg/errors"
)

func QuestionByText(inputText *string, setting *Setting, tmpflag bool, id *string) (*string, error) {
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
		return nil, errors.Errorf("入力テキスト空でした。")
	}

	if tmpflag {
		p := filepath.Join(*id, "input.txt")
		q := fmt.Sprintf("%v", messages)
		OutputTextForCheck(&p, &q)
	}
	//log.Printf("%v", messages)

	// idは一時ディレクトリのパス
	response := GetOpenAIResponse(&messages, &setting.OpenaiURL, &setting.AiModel, &setting.Apikey, id, setting.Temperature, tmpflag)
	//log.Printf("response: %v", response)
	//log.Printf(": %v", response)
	if response == nil {
		return nil, errors.Errorf("no response.")
	}

	if len((*response).Choices) == 0 {
		//log.Printf("Error: no response.")
		return nil, errors.Errorf("response:%v", response)
	}

	output := (*response).Choices[0].Messages.Content
	if tmpflag {
		p := filepath.Join(*id, "output.txt")
		OutputTextForCheck(&p, &output)
	}
	return &output, nil
}
