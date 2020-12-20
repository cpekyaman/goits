package monitoring

import (
	"context"

	"github.com/cpekyaman/goits/framework/commons"
)

const (
	MonitoringCtxKey = "monitoringContext"
)

// MonitoringContext represents a context component which provides monitoring related helpers and identifiers.
type MonitoringContext interface {
	CID() string
	REQID() string
	Resource() string
	SetResource(name string)
	Operation() string
	SetOperation(name string)
	SetError(err commons.AppError)
	GetError() commons.AppError
	HasError() bool
	Logger() Logger
}

// defaultMonitoringContext, as its name suggests, is the concrete implementation of MonitoringContext.
type defaultMonitoringContext struct {
	cid       string
	reqid     string
	resource  string
	operation string
	err       commons.AppError
	hasError  bool
	logger    *Logger
}

// CID returns the correlation id set for current context.
func (this *defaultMonitoringContext) CID() string {
	return this.cid
}

// REQID returns the request it set for the current context.
func (this *defaultMonitoringContext) REQID() string {
	return this.reqid
}

// Resource returns the name of of api resource that started request processing.
func (this *defaultMonitoringContext) Resource() string {
	return this.resource
}

// SetResource is used by api handlers to set resource name of the context.
func (this *defaultMonitoringContext) SetResource(name string) {
	this.resource = name
}

// Operation returns the name of api operation that is currently processing.
func (this *defaultMonitoringContext) Operation() string {
	return this.operation
}

// SetOperation is used by api handlers to set currently executing operation of the context.
func (this *defaultMonitoringContext) SetOperation(name string) {
	this.operation = name
}

// SetError sets an error raised during request processing for current context.
func (this *defaultMonitoringContext) SetError(err commons.AppError) {
	this.err = err
	this.hasError = true
}

// GetError returns any error that is raised during request processing.
func (this *defaultMonitoringContext) GetError() commons.AppError {
	return this.err
}

// HasError returns true if current context has some associated error.
func (this *defaultMonitoringContext) HasError() bool {
	return this.hasError
}

// Logger creates a Logger that has its context keys pre-set from this context.
func (this *defaultMonitoringContext) Logger() Logger {
	if this.logger == nil {
		l := NewLoggerWithMonitoring(this)
		this.logger = &l
	}
	return *this.logger
}

// WithMonitoringContext creates a new MonitoringContext and returns a child context for request that includes it.
func WithMonitoringContext(ctx context.Context, cid string, reqid string) context.Context {
	mctx := &defaultMonitoringContext{
		cid:       cid,
		reqid:     reqid,
		resource:  "None",
		operation: "None",
	}

	return context.WithValue(ctx, MonitoringCtxKey, mctx)
}

// GetMonitoringContext gets the MonitoringContext from given request context.
func GetMonitoringContext(ctx context.Context) (MonitoringContext, bool) {
	mValue := ctx.Value(MonitoringCtxKey)
	if mValue == nil {
		return nil, false
	}

	mctx, ok := mValue.(MonitoringContext)
	if !ok {
		return nil, false
	}
	return mctx, true
}

// GetContextLogger gets or creates the logger associated with current context.
func GetContextLogger(ctx context.Context) Logger {
	mctx, ok := GetMonitoringContext(ctx)
	if ok {
		return mctx.Logger()
	}
	return NewLogger()
}

// SetContextError gets the current monitoring context and sets the active error on it.
func SetContextError(ctx context.Context, err commons.AppError) {
	mctx, ok := GetMonitoringContext(ctx)
	if ok {
		mctx.SetError(err)
	}
}
