package routing

import (
	"net/http"
	"strconv"

	"github.com/cpekyaman/goits/framework/services"
)

const (
	rowsPerPage = 20
)

type ApiHandler interface {
	Handle(w http.ResponseWriter, r *http.Request)
}

type ApiHandlerFunc func(w http.ResponseWriter, r *http.Request)

func (hf ApiHandlerFunc) Handle(w http.ResponseWriter, r *http.Request) {
	hf(w, r)
}

func (this ApiResource) GetAll(w http.ResponseWriter, r *http.Request) {
	si, ok := this.service.(services.GetAllService)
	if !ok {
		this.notImplementedResponse(w, r)
		return
	}

	var payload interface{}
	var err error

	page, err := strconv.ParseUint(r.URL.Query().Get("page"), 10, 64)
	if page > 0 {
		payload, err = si.GetAllPaged(r.Context(), rowsPerPage, (page-1)*rowsPerPage)
	} else {
		payload, err = si.GetAll(r.Context())
	}

	if err != nil {
		this.errorResponse(w, r, "could not list resource", err)
	} else {
		this.successResponse(w, r, payload)
	}
}

func (this ApiResource) GetById(w http.ResponseWriter, r *http.Request) {
	si, ok := this.service.(services.GetByIdService)
	if !ok {
		this.notImplementedResponse(w, r)
		return
	}

	id, err := this.binder.IdPathParam(r, "id")
	if err != nil {
		this.errorResponse(w, r, "invalid input", err)
		return
	}

	payload, err := si.GetById(r.Context(), id)
	if err != nil {
		this.errorResponse(w, r, "could not get resource", err)
	} else {
		this.successResponse(w, r, payload)
	}
}

func (this ApiResource) Create(w http.ResponseWriter, r *http.Request) {
	si, ok := this.service.(services.CreateService)
	if !ok {
		this.notImplementedResponse(w, r)
		return
	}

	ob := this.binder.BindFunc(r)
	err := si.Create(r.Context(), ob)
	if err != nil {
		this.errorResponse(w, r, "could not create resource", err)
	} else {
		this.successResponse(w, r, nil)
	}
}

func (this ApiResource) Update(w http.ResponseWriter, r *http.Request) {
	si, ok := this.service.(services.UpdateService)
	if !ok {
		this.notImplementedResponse(w, r)
		return
	}

	id, err := this.binder.IdPathParam(r, "id")
	if err != nil {
		this.errorResponse(w, r, "invalid input", err)
		return
	}

	ob := this.binder.BindFunc(r)
	err = si.Update(r.Context(), id, ob)
	if err != nil {
		this.errorResponse(w, r, "could not update resource", err)
	} else {
		this.successResponse(w, r, nil)
	}
}
