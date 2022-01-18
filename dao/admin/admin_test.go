package admin

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateAdmin(t *testing.T) {
	uid := int64(10)
	ctx := context.TODO()

	p, err := QueryByUID(ctx, uid)
	assert.Nil(t, err)
	assert.Equal(t, p.ID, int64(0))

	p, err = QueryByUsername(ctx, "test")
	assert.Nil(t, err)
	assert.Equal(t, p.ID, int64(0))

	p.Username = "test"
	p.Password = "123456"
	id, err := CreateAdmin(ctx, p)
	assert.Nil(t, err)
	p.ID = id

	dst, err := QueryByUID(ctx, p.ID)
	assert.Nil(t, err)
	assert.Equal(t, p.Username, dst.Username)
	assert.Equal(t, p.Password, dst.Password)

	dst, err = QueryByUsername(ctx, p.Username)
	assert.Nil(t, err)
	assert.Equal(t, p.Username, dst.Username)
	assert.Equal(t, p.Password, dst.Password)

	err = deleteByID(ctx, p.ID)
	assert.Nil(t, err)
}
