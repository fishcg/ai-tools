package action

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Default struct {
	method  Method
	handler gin.HandlerFunc
}

func (d *Default) GetMethod() Method {
	return d.method
}

func (d *Default) GetHandler() gin.HandlerFunc {
	return d.handler
}

func NewDefault(method Method, actionFunc ActionFunc) *Default {
	return &Default{
		method,
		func(c *gin.Context) {
			_, r, _ := actionFunc(c)
			c.String(http.StatusOK, fmt.Sprintf("%v", r))
		},
	}
}
