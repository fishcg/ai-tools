package openai

import (
	"context"
	"fmt"

	goopenai "github.com/sashabaranov/go-openai"
)

var promptCodeVariable2Name = "你是一个代码生成器，将输入的文本总结命名为使用驼峰格式的程序变量，并将它输出，你需要遵循以下规则：" +
	"1.仅输入变量本身;" +
	"2. \"剧集\"翻译为 drama, \"音单\" 翻译为 album, \"音频\" 翻译为 sound, \"关注\" 翻译为 follow, \"动态\" 翻译为 feed, \"消息\" 翻译为 msg, \"提醒\" 翻译为 notice"

// CodeVariable2Name .
func (c *Client) CodeVariable2Name(content string) (string, error) {
	ctx := context.Background()
	content = formatContent(content, SceneVariable2Name)
	req := goopenai.ChatCompletionRequest{
		Model:     goopenai.GPT3Dot5Turbo0301,
		MaxTokens: 1000,
		Messages: []goopenai.ChatCompletionMessage{
			{
				Role:    goopenai.ChatMessageRoleSystem,
				Content: promptCodeVariable2Name,
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
		return "", errOpenAIResp
	}
	return resp.Choices[0].Message.Content, nil
}

var promptCodeCreateGoUnitTestFunc = "你是一个 go 语言代码生成器，对输入的方法或函数进行分析，并输出该方法或函数的单元测试方法，你需要遵循以下规则：" +
	"1.仅输入单元测试方法或函数本身; 2. 但测试函数入参为 t，类型为 *testing.T;" +
	"2. 单元测试开始时需要先定义 assert := assert.New(t) 或者 require := require.New(t) 并合理使用它们;" +
	"3. 变量使用驼峰命名" +
	"4. 输入的函数或方法可能使用了 gin、gorm、redis/v7 包" +
	"5. 不能使用 sqlmock 包"

// CodeCreateGoUnitTestFunc .
func (c *Client) CodeCreateGoUnitTestFunc(content string) (string, error) {
	return c.codeCreateUnitTestFunc(content, SceneGoUnitTestFunc)
}

var promptCodeCreatePHPUnitTestFunc = "你是一个 PHP 代码生成器，对输入的方法或函数进行分析，并输出该方法或函数的 phpinit 单元测试方法，你需要遵循以下规则：" +
	"1.仅输入单元测试方法或函数本身;" +
	"2. 变量不使用驼峰命名，使用下划线分割命名" +
	"3. 输入的函数或方法可能使用了 Yii1.1 或 Yii2 框架"

// CodeCreatePHPUnitTestFunc .
func (c *Client) CodeCreatePHPUnitTestFunc(content string) (string, error) {
	return c.codeCreateUnitTestFunc(content, ScenePHPUnitTestFunc)
}

// CodeCreateUnitTestFunc .
func (c *Client) codeCreateUnitTestFunc(content string, scene int64) (string, error) {
	var prompt string
	switch scene {
	case SceneGoUnitTestFunc:
		prompt = promptCodeCreateGoUnitTestFunc
	case ScenePHPUnitTestFunc:
		prompt = promptCodeCreatePHPUnitTestFunc
	default:
		panic(fmt.Sprintf("Scene param err: %d", scene))
	}
	ctx := context.Background()
	req := goopenai.ChatCompletionRequest{
		Model:     goopenai.GPT3Dot5Turbo0301,
		MaxTokens: 1000,
		Messages: []goopenai.ChatCompletionMessage{
			{
				Role:    goopenai.ChatMessageRoleSystem,
				Content: prompt,
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
		return "", errOpenAIResp
	}
	return fmt.Sprintf("%s%s%s", "<pre><code class=\"launght\">", resp.Choices[0].Message.Content, "</code></pre>"), nil
}
