package routing

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/cpekyaman/goits/framework/services"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type RoutingEngine struct {
	router *chi.Mux
}

type engineRequestBinder struct{}

func (this engineRequestBinder) PathParam(r *http.Request, key string) string {
	return chi.URLParam(r, key)
}

func (this engineRequestBinder) IdPathParam(r *http.Request, key string) (uint64, error) {
	idstr := this.PathParam(r, "id")
	id, err := strconv.ParseUint(idstr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("binding: %s", err.Error())
	}
	return id, nil
}

func (this engineRequestBinder) BindFunc(r *http.Request) services.ObjectBinder {
	return services.ObjectBinderFunc(func(target interface{}) error {
		if r == nil || r.Body == nil {
			return fmt.Errorf("binding: empty request")
		}
		err := render.DecodeJSON(r.Body, target)
		if err != nil {
			return fmt.Errorf("binding: %s", err.Error())
		}
		return nil
	})
}

type engineResponseRenderer struct{}

func (this engineResponseRenderer) JSON(w http.ResponseWriter, r *http.Request, data interface{}) {
	render.JSON(w, r, data)
}

var engine RoutingEngine

func InitRouting() {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(Monitor)
	r.Use(middleware.RealIP)
	r.Use(Logger)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	engine = RoutingEngine{r}
	binder = engineRequestBinder{}
	renderer = engineResponseRenderer{}
}

func routeWith(e *chi.Mux) {
	engine = RoutingEngine{e}
}

func Engine() RoutingEngine {
	return engine
}

func (this RoutingEngine) Router() *chi.Mux {
	return this.router
}

func (this RoutingEngine) Register(resource ApiResource) {
	Register(this.router, resource)
}

func (this RoutingEngine) RegisterPath(path string, h http.Handler) {
	this.router.Handle(path, h)
}

func Register(r *chi.Mux, resource ApiResource) {
	r.Route("/"+resource.path, func(r chi.Router) {
		r.Get("/", MonitoredHandler(resource.name, "getAll", resource.GetAll))
		r.Post("/", MonitoredHandler(resource.name, "create", resource.Create))

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", MonitoredHandler(resource.name, "getById", resource.GetById))
			r.Put("/", MonitoredHandler(resource.name, "update", resource.Update))
		})
	})
}
