package monitoring

import (
	"context"
	"testing"

	"github.com/cpekyaman/goits/framework/commons"
	"github.com/stretchr/testify/assert"
)

func TestMonitoringContext_Create(t *testing.T) {
	// given
	ctx := context.Background()
	cid := "123456"
	reqid := "654321"

	// when
	newCtx := WithMonitoringContext(ctx, cid, reqid)

	// then
	mctx, ok := GetMonitoringContext(newCtx)
	assert.True(t, ok, "could not get monitoring context")
	assert.Equal(t, cid, mctx.CID(), "cid is not set properly")
	assert.Equal(t, reqid, mctx.REQID(), "reqid is not set properly")
}

func TestMonitoringContext_SetError(t *testing.T) {
	// given
	ctx := context.Background()
	cid := "123456"
	reqid := "654321"

	newCtx := WithMonitoringContext(ctx, cid, reqid)
	mctx, ok := GetMonitoringContext(newCtx)
	assert.True(t, ok, "could not get monitoring context")

	err := commons.AppError{
		ErrorType: commons.ErrClient,
		Message:   "test",
		Cause:     "crap",
	}

	// when
	mctx.SetError(err)

	// then
	assert.True(t, mctx.HasError(), "has error should be true")
	assert.Equal(t, err, mctx.GetError(), "error object is not the same")
}
