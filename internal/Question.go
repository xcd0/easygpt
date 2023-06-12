package internal

import (
	"fmt"
	"path/filepath"
)

type path = string

func Question(apikey string, prompt *string, input path, tmpdir path, tmpflag bool) string {
	in := GetTextNoError(input)

	question := fmt.Sprintf("%v\n%v", *prompt, in)

	if tmpflag {
		OutputTextForCheck(filepath.Join(tmpdir, "input.txt"), question)
	}

	messages := []Message{
		Message{
			Role:    "user",
			Content: question,
		},
	}

	response := GetOpenAIResponse(&messages, apikey, tmpdir, tmpflag)
	//log.Printf("response: %v", response)
	//log.Printf(": %v", response)

	if len(response.Choices) == 0 {
		//log.Printf("Error: no response.")
		//log.Printf("       %v", response)
		return ""
	}

	output := response.Choices[0].Messages.Content
	if tmpflag {
		OutputTextForCheck(filepath.Join(tmpdir, "output.txt"), output)
	}
	return output
}
