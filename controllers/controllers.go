package controllers

import (
	"github.com/gin-gonic/gin"

	"github.com/fish/ai-tools/controllers/chat"
	"github.com/fish/ai-tools/controllers/helper"
)

// Mount server level middlerwares and controllers
type Mount func(r *gin.Engine)

// MountGame mount all controllers for game server
func MountGame(r *gin.Engine) {
	controllers := []*helper.Controller{
		chat.NewController(),
	}
	// 此处可以添加全局的 middlewares
	for _, c := range controllers {
		c.RegistAction(r)
	}
}
