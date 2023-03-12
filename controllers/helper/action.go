package helper

import (
	"github.com/fish/ai-tools/controllers/helper/action"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Action struct {
	method  action.Method
	handler gin.HandlerFunc
}

type ActionFunc func(c *Context) (*HTMLResp, RespInfo, error)

type HTMLResp struct {
	Path string
	Obj  gin.H
}

func (a *Action) GetHandler() gin.HandlerFunc {
	return a.handler
}

func (a *Action) GetMethod() action.Method {
	return a.method
}

func NewAction(method action.Method, actionFunc ActionFunc) *Action {
	action := &Action{
		method: method,
		handler: func(c *gin.Context) {
			// TODO: 实现 Request ID
			requestID := "233"
			var err error
			ctx := &Context{
				C: c,
			}

			var r RespInfo
			var htmlresp *HTMLResp
			if err == nil {
				htmlresp, r, err = actionFunc(ctx)
			}

			if err != nil {
				switch v := err.(type) {
				case HTTPError:
					AbortWithError(c, v, requestID)
				default:
					AbortWithError(c, ErrInternalServerError, requestID)
				}

			} else if htmlresp != nil {
				c.HTML(http.StatusOK, htmlresp.Path, htmlresp.Obj)
			} else {
				c.JSON(http.StatusOK, gin.H{
					"code": CodeOK,
					"info": r,
				})
			}
		},
	}
	return action
}

func AbortWithError(c *gin.Context, err HTTPError, requestID string) {
	c.AbortWithStatusJSON(err.Status(), gin.H{
		"code": err.Code(),
		// "request_id": requestID,
		"message": err.Error(),
	})
}
