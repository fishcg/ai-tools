package chat

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/fish/ai-tools/controllers/helper"
	"github.com/fish/ai-tools/controllers/helper/action"
	"github.com/fish/ai-tools/service"
	"github.com/fish/ai-tools/service/openai"
)

func NewController() *helper.Controller {
	return &helper.Controller{
		Path:        "chat",
		Middlewares: gin.HandlersChain{},
		Actions: action.Actions{
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
		res, err = service.OpenAI.FreeChat(content)
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
		fmt.Println(err)
		return nil, nil, helper.ErrInternalServerError
	}
	return nil, res, nil
}

func actionError(c *helper.Context) (*helper.HTMLResp, helper.RespInfo, error) {
	return nil, nil, helper.ErrInternalServerError
}
