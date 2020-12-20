package routing

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cpekyaman/goits/framework/monitoring"
	"github.com/stretchr/testify/assert"
)

func TestMonitor_ShouldUseCIDHeader_WhenExists(t *testing.T) {
	// given
	req, rw := setup(t)

	fakeCID := "123456"
	req.Header.Add(HDR_CorrelationID, fakeCID)

	var updatedReq *http.Request
	m := Monitor(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		updatedReq = r
	}))

	// when
	m.ServeHTTP(rw, req)

	// then
	mctx, ok := monitoring.GetMonitoringContext(updatedReq.Context())
	assert.True(t, ok, "could not get monitoring context")
	assert.Equal(t, fakeCID, mctx.CID(), "cid is not correct")
}

func TestMonitor_ShouldCreateCID_NoHeaderExists(t *testing.T) {
	// given
	req, rw := setup(t)

	var updatedReq *http.Request
	m := Monitor(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		updatedReq = r
	}))

	// when
	m.ServeHTTP(rw, req)

	// then
	mctx, ok := monitoring.GetMonitoringContext(updatedReq.Context())
	assert.True(t, ok, "could not get monitoring context")
	assert.NotEmpty(t, mctx.CID(), "cid is not correct")
}

func TestMonitoredHandler_ShouldSetResource(t *testing.T) {
	// given
	req, rw := setup(t)
	reqWithContext := req.WithContext(monitoring.WithMonitoringContext(req.Context(), "cid", "rid"))

	resource := "demo"
	operation := "test_operation"
	m := MonitoredHandler(resource, operation, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))

	// when
	m.ServeHTTP(rw, reqWithContext)

	// then
	mctx, ok := monitoring.GetMonitoringContext(reqWithContext.Context())
	assert.True(t, ok, "could not get monitoring context")
	assert.Equal(t, resource, mctx.Resource(), "resource is not set")
	assert.Equal(t, operation, mctx.Operation(), "operation is not set")
}

func setup(t *testing.T) (*http.Request, *httptest.ResponseRecorder) {
	rw := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "http://localhost:8080/test", nil)
	assert.Nil(t, err, "could not create request")

	return req, rw
}
