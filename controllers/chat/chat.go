package chat

import (
	"encoding/json"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"github.com/fish/ai-tools/controllers/helper"
	"github.com/fish/ai-tools/controllers/helper/action"
	"github.com/fish/ai-tools/logger"
	"github.com/fish/ai-tools/service"
	"github.com/fish/ai-tools/service/openai"
	openaichat "github.com/fish/ai-tools/service/openai/chat"
)

// SessionChatInfoKey 存储聊天记录的 session key
const SessionChatInfoKey = "chat_info"

func NewController() *helper.Controller {
	return &helper.Controller{
		"chat",
		gin.HandlersChain{},
		action.Actions{
			"/index": helper.NewAction(action.GET, actionIndex),
			"/get":   helper.NewAction(action.GET, actionGet),
			"/error": helper.NewAction(action.GET, actionError),
		},
	}
}

func actionIndex(c *helper.Context) (*helper.HTMLResp, helper.RespInfo, error) {
	return &helper.HTMLResp{
		Path: "chat/index.tmpl",
		Obj: gin.H{
			"title": "Posts",
		},
	}, nil, nil
}

func actionGet(c *helper.Context) (*helper.HTMLResp, helper.RespInfo, error) {
	content := c.GetQuery("content")
	scene, err := c.GetQueryInt64("scene")
	if content == "" || err != nil {
		return nil, nil, helper.ErrInvalidParam
	}
	var res string
	switch scene {
	case openai.SceneFreeChat:
		session := sessions.Default(c.C)
		chatInfo := session.Get(SessionChatInfoKey)
		chat := openaichat.NewChat(scene, content)
		if chatInfo != nil {
			chatStr, ok := chatInfo.(string)
			if !ok {
				// 内容错误则删除
				session.Delete(SessionChatInfoKey)
			} else {
				err := json.Unmarshal([]byte(chatStr), chat)
				if err != nil {
					// 内容错误则删除
					session.Delete(SessionChatInfoKey)
					logger.Error(err.Error())
					// PASS
				} else if len(chat.HistoryMsg) >= openaichat.HistoryMsgMaxCount {
					// 超过限制条数则重新开始会话
					session.Delete(SessionChatInfoKey)
					chat.HistoryMsg = nil
				}
			}
		}
		res, err = chat.Reply()
		if err != nil {
			logger.Error(err.Error())
			return nil, nil, helper.ErrInternalServerError
		}
		// Save session
		chatByte, err := json.Marshal(chat)
		if err != nil {
			logger.Error(err.Error())
			// PASS
		} else {
			session.Set(SessionChatInfoKey, string(chatByte))
			err = session.Save()
			if err != nil {
				session.Delete(SessionChatInfoKey)
				_ = session.Save()
				logger.Error(err.Error())
				// PASS
			}
		}
	case openai.SceneTextLint:
		res, err = service.OpenAI.TextLint(content)
	case openai.SceneVariable2Name:
		res, err = service.OpenAI.CodeVariable2Name(content)
	case openai.SceneTranslate:
		res, err = service.OpenAI.Translate(content)
	case openai.SceneGoUnitTestFunc:
		res, err = service.OpenAI.CodeCreateGoUnitTestFunc(content)
	case openai.ScenePHPUnitTestFunc:
		res, err = service.OpenAI.CodeCreatePHPUnitTestFunc(content)
	default:
		// TODO: Support for more scenarios
		return nil, nil, helper.ErrInvalidParam
	}
	if err != nil {
		logger.Error(err.Error())
		return nil, nil, helper.ErrInternalServerError
	}
	return nil, res, nil
}

func actionError(c *helper.Context) (*helper.HTMLResp, helper.RespInfo, error) {
	return nil, nil, helper.ErrInternalServerError
}
