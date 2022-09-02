package mid

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/go-goim/core/pkg/types"
)

func TestNewJwtToken(t *testing.T) {
	var uid = types.ID(1)
	jwt, err := NewJwtToken(uid)
	assert.Nil(t, err)
	assert.NotEmpty(t, jwt)

	c, err := ParseJwtToken(jwt)
	assert.Nil(t, err)
	assert.Equal(t, uid, c.UserID)
}
