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

	"github.com/pkg/errors"
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

const (
	openaiURL = "https://api.openai.com/v1/chat/completions"
	aiModel   = "gpt-3.5-turbo"
)

func GetOpenAIResponse(messages *[]Message, apiKey, tmpdir string, tmpflag bool) OpenaiResponse {
	var req *http.Request = CreateHttpRequest(messages, apiKey, tmpdir, tmpflag)

	var body []byte = GetResponseBody(req, tmpdir, tmpflag)

	return func(body []byte, messages *[]Message) OpenaiResponse {
		var response OpenaiResponse
		if err := json.Unmarshal(body, &response); err != nil {
			log.Printf("Error: %v", err.Error())
			return OpenaiResponse{}
		}
		if len(response.Choices) == 0 {
			//log.Printf("Error: レスポンスがありませんでした。")
			//log.Printf("       %v", response)
			return OpenaiResponse{}
		}
		*messages = append(*messages, Message{
			Role:    "assistant",
			Content: response.Choices[0].Messages.Content,
		})
		return response
	}(body, messages)
}

func CreateHttpRequest(messages *[]Message, apiKey, tmpdir string, tmpflag bool) *http.Request {
	requestBody := OpenaiRequest{
		Model:    aiModel,
		Messages: *messages,
	}

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
	return req
}

func GetResponseBody(req *http.Request, tmpdir string, tmpflag bool) []byte {
	client := &http.Client{}
	//log.Printf("Request: %v", *req)
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

	func() {
		var responseError struct {
			Error struct {
				Message string `json:"message"`
				Type    string `json:"type"`
				Param   string `json:"param"`
				Code    string `json:"code"`
			} `json:"error"`
		}
		err := json.Unmarshal(body, &responseError)
		if err != nil {
			log.Printf("Error: %v", err)
		}
		if responseError.Error.Code == "invalid_api_key" {
			log.Printf("%v", string(body))
			//log.Printf("%v", responseError.Error.Code)
			err := errors.Errorf("APIキーが不正です。公式サイトからAPIキーを再発行して、設定してください。\nhttps://platform.openai.com/account/api-keys\n\n")
			//fmt.Printf("%+v", err)
			fmt.Printf("%v", err)
		}
	}()

	j := JsonFormat(body)
	//log.Printf("response body:\n%v", j)
	if tmpflag {
		OutputTextForCheck(filepath.Join(tmpdir, "response.json"), j)
	}
	return body
}
