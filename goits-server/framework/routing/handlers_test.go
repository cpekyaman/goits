package routing

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"

	"github.com/cpekyaman/goits/framework/commons"
	"github.com/cpekyaman/goits/framework/services"
	"github.com/cpekyaman/goits/framework/testlib/mocking"
	"github.com/go-chi/chi"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

const (
	rootUrl = "http://localhost"
)

// a service that does not ımplement anything
type NoOpService struct {
}

type TestEntity struct {
	Id   uint64 `json:"id"`
	Name string `json:"name"`
}

type ApiResponse struct {
	Success bool
	Error   commons.AppError
	Data    json.RawMessage
}

func TestGetAll_NotImplemented(t *testing.T) {
	// given
	assertNotImplemented(t, func(api ApiResource) ApiHandler {
		return ApiHandlerFunc(api.GetAll)
	})
}

func TestGetAll_Error(t *testing.T) {
	getAllVariant_Error(t, rootUrl+"/test", func(m *mocking.MockGetAllService) {
		m.EXPECT().GetAll(goContextMatcher()).Return(nil, fmt.Errorf("service error"))
	})
}

func TestGetAll_Paged_Error(t *testing.T) {
	getAllVariant_Error(t, rootUrl+"/test?page=2", func(m *mocking.MockGetAllService) {
		m.EXPECT().GetAllPaged(goContextMatcher(), uint(20), uint64(20)).Return(nil, fmt.Errorf("service error"))
	})
}

func getAllVariant_Error(t *testing.T, url string, mocker func(*mocking.MockGetAllService)) {
	// given
	ctrl := gomock.NewController(t)

	rw := httptest.NewRecorder()
	req, err := http.NewRequest("GET", url, nil)
	assert.Nil(t, err, "could not create request")

	svc := mocking.NewMockGetAllService(ctrl)
	mocker(svc)
	_, r := newTestApiResource(svc)

	// when
	r.ServeHTTP(rw, req)

	// then
	assert.Equal(t, http.StatusInternalServerError, rw.Result().StatusCode, "status code is not correct")
	assertErrorInResponse(t, rw, "could not list resource", "service error")
}

func TestGetAll_Success(t *testing.T) {
	getAllVariant_Success(t, rootUrl+"/test", func(m *mocking.MockGetAllService, expectedList interface{}) {
		m.EXPECT().GetAll(goContextMatcher()).Return(expectedList, nil)
	})
}

func TestGetAll_Paged_Success(t *testing.T) {
	getAllVariant_Success(t, rootUrl+"/test?page=1", func(m *mocking.MockGetAllService, expectedList interface{}) {
		m.EXPECT().GetAllPaged(goContextMatcher(), uint(20), uint64(0)).Return(expectedList, nil)
	})
}

func getAllVariant_Success(t *testing.T, url string, mocker func(*mocking.MockGetAllService, interface{})) {
	// given
	ctrl := gomock.NewController(t)

	expectedList := []TestEntity{
		{Id: 1, Name: "First"},
		{Id: 2, Name: "Second"},
	}

	rw := httptest.NewRecorder()
	req, err := http.NewRequest("GET", url, nil)
	assert.Nil(t, err, "could not create request")

	svc := mocking.NewMockGetAllService(ctrl)
	mocker(svc, expectedList)
	_, r := newTestApiResource(svc)

	// when
	r.ServeHTTP(rw, req)

	// then
	assert.Equal(t, http.StatusOK, rw.Result().StatusCode, "status should be success")

	response := ApiResponse{}
	err = json.Unmarshal(rw.Body.Bytes(), &response)
	assert.Nil(t, err, "error in unmarshal response")
	actualList := []TestEntity{}
	err = json.Unmarshal(response.Data, &actualList)
	assert.Nil(t, err, "error in unmarshal data")
	assert.Equal(t, len(expectedList), len(actualList), "number elements are not the same")

	var actual TestEntity
	var expected TestEntity
	for i := 0; i < len(expectedList); i++ {
		actual = actualList[i]
		expected = expectedList[i]
		assert.Equal(t, expected.Id, actual.Id, "id of element is not the same")
		assert.Equal(t, expected.Name, actual.Name, "name of element is not the same")

	}
}

func TestGetById_NotImplemented(t *testing.T) {
	assertNotImplemented(t, func(api ApiResource) ApiHandler {
		return ApiHandlerFunc(api.GetById)
	})
}

func TestGetById_InvalidId_Error(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)

	rw := httptest.NewRecorder()
	req, err := http.NewRequest("GET", rootUrl+"/test/garbage", nil)
	assert.Nil(t, err, "could not create request")

	svc := mocking.NewMockGetByIdService(ctrl)
	svc.EXPECT().GetById(goContextMatcher(), gomock.Any()).Times(0)

	_, r := newTestApiResource(svc)

	// when
	r.ServeHTTP(rw, req)

	// then
	assert.Equal(t, http.StatusBadRequest, rw.Result().StatusCode)
}

func TestGetById_Error(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)

	rw := httptest.NewRecorder()
	var id uint64 = 5
	req, err := http.NewRequest("GET", rootUrl+"/test/"+strconv.FormatUint(id, 10), nil)
	assert.Nil(t, err, "could not create request")

	svc := mocking.NewMockGetByIdService(ctrl)
	svc.EXPECT().GetById(goContextMatcher(), id).Times(1).Return(nil, fmt.Errorf("service error"))
	_, r := newTestApiResource(svc)

	// when
	r.ServeHTTP(rw, req)

	// then
	assert.Equal(t, http.StatusInternalServerError, rw.Result().StatusCode)
	assertErrorInResponse(t, rw, "resource not found", "service error")
}

func TestGetById_Success(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)

	rw := httptest.NewRecorder()
	var id uint64 = 5
	req, err := http.NewRequest("GET", rootUrl+"/test/"+strconv.FormatUint(id, 10), nil)
	assert.Nil(t, err, "could not create request")

	expected := TestEntity{
		Id:   id,
		Name: "Another One",
	}

	svc := mocking.NewMockGetByIdService(ctrl)
	svc.EXPECT().GetById(goContextMatcher(), id).Times(1).Return(expected, nil)
	_, r := newTestApiResource(svc)

	// when
	r.ServeHTTP(rw, req)

	// then
	assert.Equal(t, http.StatusOK, rw.Result().StatusCode)

	response := ApiResponse{}
	err = json.Unmarshal(rw.Body.Bytes(), &response)
	assert.Nil(t, err, "error in unmarshal response")

	actual := TestEntity{}
	err = json.Unmarshal(response.Data, &actual)

	assert.Equal(t, expected.Id, actual.Id, "id of element is not the same")
	assert.Equal(t, expected.Name, actual.Name, "name of element is not the same")
}

func TestCreate_NotImplemented(t *testing.T) {
	assertNotImplemented(t, func(api ApiResource) ApiHandler {
		return ApiHandlerFunc(api.Create)
	})
}

func TestCreate_EmptyRequest_Error(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)

	rw := httptest.NewRecorder()
	req, err := http.NewRequest("POST", rootUrl+"/test", nil)
	assert.Nil(t, err, "could not create request")

	var te TestEntity

	svc := mocking.NewMockCreateService(ctrl)
	svc.EXPECT().Create(goContextMatcher(), gomock.Any()).
		Times(1).
		DoAndReturn(func(ctx context.Context, binding services.ObjectBinder) error {
			return binding.BindTo(&te)
		})
	_, r := newTestApiResource(svc)

	// when
	r.ServeHTTP(rw, req)

	// then
	assert.Equal(t, http.StatusBadRequest, rw.Result().StatusCode)
	assertErrorInResponse(t, rw, "could not create resource", "invalid request")
}

func TestCreate_Success(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)

	reqEntity := TestEntity{
		Name: "New Name",
	}
	jsonBytes, err := json.Marshal(reqEntity)
	if err != nil {
		assert.Fail(t, "could not mock request body")
	}
	body := bytes.NewReader(jsonBytes)

	rw := httptest.NewRecorder()
	req, err := http.NewRequest("POST", rootUrl+"/test", body)
	assert.Nil(t, err, "could not create request")

	var created TestEntity

	svc := mocking.NewMockCreateService(ctrl)
	svc.EXPECT().Create(goContextMatcher(), gomock.Any()).
		Times(1).
		DoAndReturn(func(ctx context.Context, binding services.ObjectBinder) error {
			return binding.BindTo(&created)
		})
	_, r := newTestApiResource(svc)

	// when
	r.ServeHTTP(rw, req)

	// then
	assert.Equal(t, http.StatusOK, rw.Result().StatusCode)
	assert.Equal(t, reqEntity.Name, created.Name, "name is not bind")
}

func TestUpdate_NotImplemented(t *testing.T) {
	assertNotImplemented(t, func(api ApiResource) ApiHandler {
		return ApiHandlerFunc(api.Update)
	})
}

func TestUpdate_InvalidId_Error(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)

	rw := httptest.NewRecorder()
	req, err := http.NewRequest("PUT", rootUrl+"/test/garbage", nil)
	assert.Nil(t, err, "could not create request")

	svc := mocking.NewMockUpdateService(ctrl)
	svc.EXPECT().Update(goContextMatcher(), gomock.Any(), gomock.Any()).Times(0)
	_, r := newTestApiResource(svc)

	// when
	r.ServeHTTP(rw, req)

	// then
	assert.Equal(t, http.StatusBadRequest, rw.Result().StatusCode)
}

func TestUpdate_EmptyRequest_Error(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)

	rw := httptest.NewRecorder()
	var id uint64 = 5
	req, err := http.NewRequest("PUT", rootUrl+"/test/"+strconv.FormatUint(id, 10), nil)
	assert.Nil(t, err, "could not create request")

	var te TestEntity

	svc := mocking.NewMockUpdateService(ctrl)
	svc.EXPECT().Update(goContextMatcher(), gomock.Eq(id), gomock.Any()).
		Times(1).
		DoAndReturn(func(ctx context.Context, id uint64, binding services.ObjectBinder) error {
			return binding.BindTo(&te)
		})
	_, r := newTestApiResource(svc)

	// when
	r.ServeHTTP(rw, req)

	// then
	assert.Equal(t, http.StatusBadRequest, rw.Result().StatusCode)
	assertErrorInResponse(t, rw, "could not update resource", "invalid request")
}

func TestUpdate_Success(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)

	var id uint64 = 5
	reqEntity := TestEntity{
		Id:   id,
		Name: "Updated Name",
	}
	jsonBytes, err := json.Marshal(reqEntity)
	if err != nil {
		assert.Fail(t, "could not mock request body")
	}
	body := bytes.NewReader(jsonBytes)

	rw := httptest.NewRecorder()
	req, err := http.NewRequest("PUT", rootUrl+"/test/"+strconv.FormatUint(id, 10), body)

	existing := TestEntity{
		Id:   id,
		Name: "Current Name",
	}

	svc := mocking.NewMockUpdateService(ctrl)
	svc.EXPECT().Update(goContextMatcher(), gomock.Eq(id), gomock.Any()).
		Times(1).
		DoAndReturn(func(ctx context.Context, id uint64, binding services.ObjectBinder) error {
			return binding.BindTo(&existing)
		})
	_, r := newTestApiResource(svc)

	// when
	r.ServeHTTP(rw, req)

	// then
	assert.Equal(t, http.StatusOK, rw.Result().StatusCode)
	assert.Equal(t, reqEntity.Name, existing.Name, "name is not bind")
}

func assertNotImplemented(t *testing.T, handlerFunc func(ApiResource) ApiHandler) {
	// given
	rw := httptest.NewRecorder()
	api := NewApiResource("Test", "test", NoOpService{})
	req, err := http.NewRequest("GET", rootUrl+"/test", nil)
	assert.Nil(t, err, "could not create request")

	// when
	handlerFunc(api).Handle(rw, req)

	// then
	assert.Equal(t, http.StatusNotImplemented, rw.Result().StatusCode)
}

func assertErrorInResponse(t *testing.T, rw *httptest.ResponseRecorder, msg string, cause string) {
	response := ApiResponse{}
	err := json.Unmarshal(rw.Body.Bytes(), &response)
	assert.Nil(t, err, "error in unmarshal response")

	assert.Equal(t, msg, response.Error.Message, "error message not correct")
	assert.Equal(t, cause, response.Error.Cause, "error cause not correct")
}

func newTestApiResource(svc interface{}) (ApiResource, *chi.Mux) {
	api := NewCustomApiResource("Test", "test", engineRequestBinder{}, engineResponseRenderer{}, svc)
	r := chi.NewRouter()
	Register(r, api)
	return api, r
}

func goContextMatcher() gomock.Matcher {
	return gomock.AssignableToTypeOf(goContextType())
}

func goContextType() reflect.Type {
	return reflect.TypeOf((*context.Context)(nil)).Elem()
}
