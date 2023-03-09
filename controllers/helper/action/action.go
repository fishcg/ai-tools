package action

import (
	"github.com/gin-gonic/gin"
)

type Method uint8

const (
	GET Method = 1 << iota
	POST
	PUT
	DELETE
	PATCH
	HEAD
)

type Action interface {
	GetMethod() Method
	GetHandler() gin.HandlerFunc
}

type RespInfo interface{}
type ActionFunc func(c *gin.Context) (bool, RespInfo, error)
type Actions map[string]Action
