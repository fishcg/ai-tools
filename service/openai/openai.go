package openai

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	goopenai "github.com/sashabaranov/go-openai"
)

// Doc: https://platform.openai.com/docs/guides/chat

// Scene
const (
	SceneFreeChat = iota + 1
	SceneTextLint
	SceneVariable2Name
	SceneTranslate
	SceneGoUnitTestFunc
	ScenePHPUnitTestFunc
)

var (
	ErrOpenAIResp = errors.New("OpenAI API resp error")
)

// Config contains configuration for Client
type Config struct {
	// TODO: yaml 解析字段暂未改动
	Token string `yaml:"token"`
	Proxy string `yaml:"proxy"`
}

// Client wraps the CommonAPI client
type Client struct {
	*goopenai.Client
}

// NewClient returns a new Client
func NewClient(conf *Config) (c *Client) {
	config := goopenai.DefaultConfig(conf.Token)

	if conf.Proxy != "" {
		proxyUrl, err := url.Parse(conf.Proxy)
		if err != nil {
			panic(err)
		}
		transport := &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		}
		config.HTTPClient = &http.Client{
			Transport: transport,
		}
	}

	return &Client{
		goopenai.NewClientWithConfig(config),
	}
}

func formatContent(content string, scene int64) string {
	// 添加全局 prompt
	switch scene {
	case SceneTextLint:
		return fmt.Sprintf("对以下文本进行格式检查，输出正确文本并指出错误的地方：\n\n%s", content)
	default:
		return content
	}
}

// TODO: 从配置或数据库中加载
// 文本 格式检查 prompt
var promptTextLint = `You are a text format checker，please output chinese. When conducting a text format check, The principles should be followed, the principle is to check that the format adheres to the following specifications:
1. There should be a space between half-width symbols and the content.
2. Correct capitalization should be used. For example, "UP 主" must be capitalized and cannot be written as "Up 主" or "up 主".
4. There should be no space between English units and Arabic numerals, such as "45s", "1.2w", and "150min".
5. Chinese characters for years, months, and days are not considered units, such as "2002年".
6. There should be a space between half-width symbols and content, such as "Views: 233". If half-width parentheses are used, a space should be added between the outside and the content (no space is needed on the inside). When connected to other symbols, spaces should not be used.
7. There should be no space between the minutes and seconds in the time format "12:56", but there should be spaces on both sides.
8. There should be no space between the date separators in formats such as "2019/0/11" or "2019-10-11".
9. A space should be added after half-width commas and periods. For example: "Download Elasticsearch, Logstash. Elasticsearch is developed in Java. 1. The entries must be original works."
10. When listing points or conditions, it is recommended to use periods at the end of each sentence (or semicolons at the end of all but the last point).
11.Avoid using full-width tilde "" in copywriting and use half-width tilde "" instead.
12. Avoid using periods as the end of prompts. Incorrect example: "Please enter your password." Correct example: "Please enter your password".
13. When expressing quantities of items, avoid using "x" or "*", and use the full-width "×" instead.
14. The second principle is to follow the following specifications when outputting text:`

// TextLint .
func (c *Client) TextLint(content string) (string, error) {
	ctx := context.Background()
	content = formatContent(content, SceneTextLint)
	req := goopenai.ChatCompletionRequest{
		Model:     goopenai.GPT3Dot5Turbo0301,
		MaxTokens: 1000,
		Messages: []goopenai.ChatCompletionMessage{
			{
				Role:    goopenai.ChatMessageRoleSystem,
				Content: promptTextLint,
			},
			{
				Role:    goopenai.ChatMessageRoleUser,
				Content: content,
			},
		},
	}
	resp, err := c.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", err
	}
	if len(resp.Choices) == 0 || resp.Choices[0].Message.Content == "" {
		return "", ErrOpenAIResp
	}
	return strings.Replace(resp.Choices[0].Message.Content, "\n", "<br>", -1), nil
}

var PromptFreeChat = "You should say Chinese"

// FreeChat .
func (c *Client) FreeChat(content string) (string, error) {
	ctx := context.Background()
	content = formatContent(content, SceneFreeChat)
	req := goopenai.ChatCompletionRequest{
		Model:     goopenai.GPT3Dot5Turbo0301,
		MaxTokens: 1000,
		Messages: []goopenai.ChatCompletionMessage{
			{
				Role:    goopenai.ChatMessageRoleSystem,
				Content: PromptFreeChat,
			},
			{
				Role:    goopenai.ChatMessageRoleUser,
				Content: content,
			},
		},
	}
	resp, err := c.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", err
	}
	if len(resp.Choices) == 0 || resp.Choices[0].Message.Content == "" {
		return "", ErrOpenAIResp
	}
	return strings.Replace(resp.Choices[0].Message.Content, "\n", "<br>", -1), nil
}

var promptTranslate = "你是一个翻译器，将输入的文本翻译为英文，遵循以下规则：" +
	"1.仅输入翻译后的语句本身;" +
	"2. 句子首字母不需要大写" +
	"2. \"剧集\"翻译为 drama, \"音单\" 翻译为 album, \"音频\" 翻译为 sound, \"关注\" 翻译为 follow, \"动态\" 翻译为 feed, \"消息\" 翻译为 msg, \"提醒\" 翻译为 notice"

// Translate .
func (c *Client) Translate(content string) (string, error) {
	ctx := context.Background()
	content = formatContent(content, SceneFreeChat)
	req := goopenai.ChatCompletionRequest{
		Model:     goopenai.GPT3Dot5Turbo0301,
		MaxTokens: 1000,
		Messages: []goopenai.ChatCompletionMessage{
			{
				Role:    goopenai.ChatMessageRoleSystem,
				Content: promptTranslate,
			},
			{
				Role:    goopenai.ChatMessageRoleUser,
				Content: content,
			},
		},
	}
	resp, err := c.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", err
	}
	if len(resp.Choices) == 0 || resp.Choices[0].Message.Content == "" {
		return "", ErrOpenAIResp
	}
	return strings.Replace(resp.Choices[0].Message.Content, "\n", "<br>", -1), nil
}
