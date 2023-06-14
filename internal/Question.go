package internal

import (
	"fmt"
	"os"
	"path/filepath"
)

func Question(openaiURL, aiModel, apikey, prompt, input, tmpdir *string, temperature float64, tmpflag bool) *string {
	inputText := GetTextNoError(input)
	return QuestionByText(openaiURL, aiModel, apikey, prompt, inputText, tmpdir, temperature, tmpflag)
}

func QuestionByText(openaiURL, aiModel, apikey, prompt, inputText, tmpdir *string, temperature float64, tmpflag bool) *string {
	//log.Printf("temperature:%v", temperature)
	messages := []Message{}
	if len(*prompt) > 0 {
		messages = append(messages,
			Message{
				Role:    "system",
				Content: *prompt,
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
		p := filepath.Join(*tmpdir, "input.txt")
		q := fmt.Sprintf("%v", messages)
		OutputTextForCheck(&p, &q)
	}
	//log.Printf("%v", messages)

	response := GetOpenAIResponse(&messages, openaiURL, aiModel, apikey, tmpdir, temperature, tmpflag)
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
		p := filepath.Join(*tmpdir, "output.txt")
		OutputTextForCheck(&p, &output)
	}
	return &output
}
