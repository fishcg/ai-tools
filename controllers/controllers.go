package controllers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"


	"github.com/fish/ai-tools/controllers/chat"
	"github.com/fish/ai-tools/controllers/helper"
)

// Mount mount server level middlerwares and controllers
type Mount func(r *gin.Engine)

// MountChat mount all controllers for game server
func MountChat(r *gin.Engine) {
	// 聊天记录暂时放 session
	// TODO: middlerware、config secret
	chatStore := memstore.NewStore([]byte("test_secret"))
	r.Use(sessions.Sessions("chat", chatStore))

	controllers := []*helper.Controller{
		chat.NewController(),
	}
	// 此处可以添加全局的 middlewares
	for _, c := range controllers {
		c.RegistAction(r)
	}
}
