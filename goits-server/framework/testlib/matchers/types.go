package matchers

import (
	"context"
	"reflect"

	"github.com/golang/mock/gomock"
)

// String returns a matcher that matches arguments of type string.
func String() gomock.Matcher {
	return gomock.AssignableToTypeOf("")
}

// GoContext returns a matcher that matches arguments of type context.Context.
func GoContext() gomock.Matcher {
	t := reflect.TypeOf((*context.Context)(nil)).Elem()
	return gomock.AssignableToTypeOf(t)
}
