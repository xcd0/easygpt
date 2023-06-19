package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"

	"github.com/pkg/errors"
)

type OpenaiRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature"`
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

// Questionからしか呼ばないつもり。
func GetOpenAIResponse(messages *[]Message, openaiURL, aiModel, apiKey, tmpdir *string, temperature float64, tmpflag bool) *OpenaiResponse {
	var req *http.Request = CreateHttpRequest(messages, openaiURL, aiModel, apiKey, tmpdir, temperature, tmpflag)
	if req == nil {
		return nil
	}

	//log.Printf("\nreq:%v\n\n*tmpdir:%v\n\ntmpflag:%v\n\n", req, *tmpdir, tmpflag)

	var body []byte = GetResponseBody(req, tmpdir, tmpflag)

	return func(body []byte, messages *[]Message) *OpenaiResponse {
		var response OpenaiResponse
		if err := json.Unmarshal(body, &response); err != nil {
			log.Printf("Error: %v", err.Error())
			return nil
		}
		if len(response.Choices) == 0 {
			//log.Printf("Error: レスポンスがありませんでした。")
			//log.Printf("       %v", response)
			return nil
		}
		*messages = append(*messages, Message{
			Role:    "assistant",
			Content: response.Choices[0].Messages.Content,
		})
		return &response
	}(body, messages)
}

func CreateHttpRequest(messages *[]Message, openaiURL, aiModel, apiKey, tmpdir *string, temperature float64, tmpflag bool) *http.Request {
	requestBody := OpenaiRequest{
		Model:       *aiModel,
		Messages:    *messages,
		Temperature: temperature,
	}

	requestJSON, _ := json.Marshal(requestBody)

	if tmpflag {
		p := filepath.Join(*tmpdir, "request.json")
		j := JsonFormat(requestJSON)
		OutputTextForCheck(&p, &j)
	}

	if len(*openaiURL) == 0 {
		fmt.Println("URL指定がありません。")
		return nil
	}
	if len(*aiModel) == 0 {
		fmt.Println("AIのモデル指定がありません。")
		return nil
	}
	if len(*apiKey) == 0 {
		fmt.Println("APIキー指定がありません。")
		return nil
	}

	req, err := http.NewRequest("POST", *openaiURL, bytes.NewBuffer(requestJSON))
	if err != nil {
		log.Printf("Error: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", *apiKey))
	return req
}

func GetResponseBody(req *http.Request, tmpdir *string, tmpflag bool) []byte {
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error: %v", err)
	}
	if resp == nil {
		log.Printf("レスポンスがnilです。")
		log.Printf("Request: %v", *req)
	}

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

	if tmpflag {
		p := filepath.Join(*tmpdir, "response.json")
		j := JsonFormat(body)
		OutputTextForCheck(&p, &j)
	}

	/*
		defer func(Body io.ReadCloser) {
			if err := Body.Close(); err != nil {
				log.Printf("Error: %v", err)
			}
		}(resp.Body)
	*/

	if err := resp.Body.Close(); err != nil {
		log.Printf("Error: %v", err)
	}

	return body
}
