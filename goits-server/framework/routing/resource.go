//go:generate mockgen -source=resource.go -destination=api_mock.go -package=mocking
package routing

import (
	"net/http"

	"github.com/cpekyaman/goits/framework/commons"
	"github.com/cpekyaman/goits/framework/monitoring"
	"github.com/cpekyaman/goits/framework/services"
)

var binder RequestBinder
var renderer ResponseRenderer

// RequestBinder represents the request processing functionality we explicitly use from actual routing engine.
// It is used from within handlers of ApiResource to minimize direct dependency on third party router.
type RequestBinder interface {
	PathParam(r *http.Request, key string) string
	IdPathParam(r *http.Request, key string) (uint64, error)
	BindFunc(r *http.Request) services.ObjectBinder
}

// ResponseRenderer represents response rendering functionality we explicitly use from actual routing engine.
// It is used from within handlers of ApiResource to minimize direct dependency on third party router.
type ResponseRenderer interface {
	JSON(w http.ResponseWriter, r *http.Request, data interface{})
}

// ApiResource represents a rest CRUD resource for a specific entity.
type ApiResource struct {
	name    string
	path    string
	service interface{}
	binder  RequestBinder
	render  ResponseRenderer
}

// NewApiResource creates a new api resource by using engine provided defaults for binder and renderer.
func NewApiResource(name string, path string, svc interface{}) ApiResource {
	return ApiResource{name, path, svc, binder, renderer}
}

// NewCustomApiResource creates a new api resource by using provided binder and renderer.
func NewCustomApiResource(name string, path string, b RequestBinder, r ResponseRenderer, svc interface{}) ApiResource {
	return ApiResource{name, path, svc, b, r}
}

// Register registers the api resource with routing engine making it available to be used via rest.
func (this ApiResource) Register() {
	engine.Register(this)
}

func (this ApiResource) successResponse(w http.ResponseWriter, r *http.Request, result interface{}) {
	this.render.JSON(w, r, map[string]interface{}{
		"success": true,
		"data":    result,
	})
}

func (this ApiResource) notImplementedResponse(w http.ResponseWriter, r *http.Request) {
	appErr := commons.AppError{
		ErrorType: commons.ErrClient,
		Message:   "invalid request",
		Cause:     "operation not supported",
	}

	monitoring.SetContextError(r.Context(), appErr)

	w.WriteHeader(http.StatusNotImplemented)
}

func (this ApiResource) errorResponse(w http.ResponseWriter, r *http.Request, msg string, err error) {
	code := http.StatusInternalServerError
	errType := commons.DetermineErrorType(err)
	if errType == commons.ErrClient || errType == commons.ErrValidation {
		code = http.StatusBadRequest
	} else if errType == commons.ErrNotFound {
		code = http.StatusNotFound
	}

	appErr := commons.AppError{
		ErrorType: errType,
		Message:   msg,
		Cause:     err.Error(),
	}

	monitoring.SetContextError(r.Context(), appErr)

	w.WriteHeader(code)
	this.render.JSON(w, r, map[string]interface{}{
		"success": false,
		"error":   appErr,
	})
}

func (this ApiResource) logger(r *http.Request) monitoring.Logger {
	return monitoring.GetContextLogger(r.Context())
}
