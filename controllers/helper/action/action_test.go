package action

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConst(t *testing.T) {
	assert := assert.New(t)
	methods := [...]Method{GET, POST, PUT, DELETE, PATCH, HEAD}
	for i, method := range methods {
		assert.Equal(method, Method(1<<uint(i)))
	}
}
