package chat

import (
	"context"
	"fmt"
	"github.com/fish/ai-tools/service/openai"

	goopenai "github.com/sashabaranov/go-openai"

	"github.com/fish/ai-tools/service"
)

const HistoryMsgMaxCount = 20

type Chat struct {
	Scene      int64             `json:"scene"`
	Content    string            `json:"-"`
	HistoryMsg []ChatHistoryMsg `json:"history_msg"`
}

type ChatHistoryMsg struct {
	Role     string `json:"role"`     // 仅支持 user 和 assistant
	KeyWords string `json:"keywords"` // 聊天信息抽取关键词填入
}

// NewChat new chat
func NewChat(scene int64, content string)  *Chat  {
	if content == "" {
		panic(fmt.Sprintf("NewChat keywords param is invalid: %s", content))
	}
	return &Chat{
		Content: content,
		Scene:   scene,
	}
}

// AppendHistory Append history msg
func (c *Chat) AppendHistoryMsg(role, content string) {
	if role != goopenai.ChatMessageRoleUser && role != goopenai.ChatMessageRoleAssistant {
		panic(fmt.Sprintf("AppendHistoryMsg role param is invalid: %s", role))
	}
	if content == "" {
		panic(fmt.Sprintf("AppendHistoryMsg keywords param is invalid: %s", content))
	}
	// FIXME: 理论上聊天记录顺序应该按照一问一答的方式一一记录，但是考虑到在回复前可能多次提问，此处先不做聊天顺序校验
	c.HistoryMsg = append(c.HistoryMsg, ChatHistoryMsg{
		Role: role,
		KeyWords: content, // TODO: 提取关键字存入以节省 token，提高效率
	})
}

// Reply
func (c *Chat) Reply() (string, error) {
	ctx := context.Background()
	messages := make([]goopenai.ChatCompletionMessage, 0, len(c.HistoryMsg)+2)
	// 添加系统设置
	messages = append(messages, goopenai.ChatCompletionMessage{
		Role:    goopenai.ChatMessageRoleSystem,
		Content: openai.PromptFreeChat,
	})
	// 添加历史消息
	for _, msg := range c.HistoryMsg {
		messages = append(messages, goopenai.ChatCompletionMessage{
			Role:    msg.Role,
			Content: msg.KeyWords,
		})
	}
	// 添加此次询问消息
	messages = append(messages, goopenai.ChatCompletionMessage{
		Role:    goopenai.ChatMessageRoleUser,
		Content: c.Content,
	})
	req := goopenai.ChatCompletionRequest{
		Model:     goopenai.GPT3Dot5Turbo0301,
		MaxTokens: 500,
		Messages: messages,
	}
	resp, err := service.OpenAI.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", err
	}
	if len(resp.Choices) == 0 || resp.Choices[0].Message.Content == "" {
		return "", openai.ErrOpenAIResp
	}
	reply := resp.Choices[0].Message.Content
	// c.Reply = strings.Replace(resp.Choices[0].Message.Content, "\n", "<br>", -1)
	// Add history
	c.AppendHistoryMsg(goopenai.ChatMessageRoleUser, c.Content)
	c.AppendHistoryMsg(goopenai.ChatMessageRoleAssistant, reply)
	return reply, nil
}
