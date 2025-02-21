package openaigen

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/metatube-community/metatube-sdk-go/translate"
	"io"
	"net/http"
)

var _ translate.Translator = (*OpenAIGen)(nil)

type OpenAIGen struct {
	Url   string `json:"openaigen-basic-url"`
	Auth  string `json:"openaigen-auth-bearer"`
	Model string `json:"openaigen-model"`
}
type Request struct {
	Prompt   string   `json:"prompt"`   // 用户输入的对话内容
	Model    string   `json:"model"`    // 使用的模型
	Messages Messages `json:"messages"` // 最大生成token数
}

type Resp struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
type Messages []Message

func (oa *OpenAIGen) Translate(q, source, target string) (result string, err error) {
	prompt := "将以下文本翻译成中文，并直接输出：" + q
	requestBody := Request{
		Model: oa.Model,
		Messages: Messages{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %v", err)
	}
	req, err := http.NewRequest(http.MethodPost, oa.Url+"/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+oa.Auth)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}
	var response Resp
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %v", err)
	}
	if len(response.Choices) > 0 {
		return response.Choices[0].Message.Content, nil
	}
	return "", errors.New("no response from Tencent AI")
}

func init() {
	translate.Register(&OpenAIGen{})
}
