package mid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewJwtToken(t *testing.T) {
	var uid int64 = 1
	jwt, err := NewJwtToken(uid)
	assert.Nil(t, err)
	assert.NotEmpty(t, jwt)

	c, err := ParseJwtToken(jwt)
	assert.Nil(t, err)
	assert.Equal(t, uid, c.UserID)
}
