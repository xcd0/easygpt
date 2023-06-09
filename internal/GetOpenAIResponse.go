package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
)

type OpenaiRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type OpenaiResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int      `json:"created"`
	Choices []Choice `json:"choices"`
	Usages  Usage    `json:"usage"`
}

type Choice struct {
	Index        int     `json:"index"`
	Messages     Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

func GetOpenAIResponse(messages *[]Message, apiKey, tmpdir string, tmpflag bool) OpenaiResponse {
	requestBody := OpenaiRequest{
		Model:    "gpt-3.5-turbo",
		Messages: *messages,
	}

	const openaiURL = "https://api.openai.com/v1/chat/completions"

	requestJSON, _ := json.Marshal(requestBody)

	if tmpflag {
		OutputTextForCheck(filepath.Join(tmpdir, "request.json"), JsonFormat(requestJSON))
	}

	req, err := http.NewRequest("POST", openaiURL, bytes.NewBuffer(requestJSON))
	if err != nil {
		log.Printf("Error: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", apiKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Error: %v", err)
		}
	}(resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error: %v", err)
	}

	if tmpflag {
		OutputTextForCheck(filepath.Join(tmpdir, "response.json"), JsonFormat(requestJSON))
	}

	var response OpenaiResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("Error: %v", err.Error())
		return OpenaiResponse{}
	}

	if len(response.Choices) == 0 {
		log.Printf("Error: レスポンスがありませんでした。")
		return OpenaiResponse{}
	}
	*messages = append(*messages, Message{
		Role:    "assistant",
		Content: response.Choices[0].Messages.Content,
	})

	return response
}

func JsonFormat(body []byte) string {
	var buf bytes.Buffer
	json.Indent(&buf, body, "", "\t")
	return buf.String()
}
