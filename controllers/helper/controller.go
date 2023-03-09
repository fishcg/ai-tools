package helper

import (
	"github.com/gin-gonic/gin"

	"github.com/fish/ai-tools/controllers/helper/action"
)

type Controller struct {
	Path        string
	Middlewares gin.HandlersChain
	Actions     action.Actions
}

func (c *Controller) RegistAction(r *gin.Engine) {
	g := r.Group(c.Path)
	g.Use(c.Middlewares...)
	for relativePath, a := range c.Actions {
		if a.GetMethod()&action.GET != 0 {
			g.GET(relativePath, a.GetHandler())
		}
		if a.GetMethod()&action.POST != 0 {
			g.POST(relativePath, a.GetHandler())
		}
	}
}
