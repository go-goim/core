package mid

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/go-goim/core/pkg/types"
)

func TestNewJwtToken(t *testing.T) {
	var uid int64 = 1
	jwt, err := NewJwtToken(types.ID(uid))
	assert.Nil(t, err)
	assert.NotEmpty(t, jwt)

	c, err := ParseJwtToken(jwt)
	assert.Nil(t, err)
	assert.Equal(t, uid, c.UserID)
}
