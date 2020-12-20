package routing

import (
	"net/http"
	"strconv"
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/cpekyaman/goits/framework/monitoring"
	"github.com/go-chi/chi/middleware"
)

var durationHist monitoring.HistogramBundle
var sizeHist monitoring.HistogramBundle
var reqCnt monitoring.CounterBundle

func init() {
	labels := []string{"resource", "operation", "status", "error"}
	durationHist = monitoring.NewHistogram("http_request_duration_ms", labels)
	sizeHist = monitoring.NewHistogram("http_response_size_bytes", labels)
	reqCnt = monitoring.NewCounter("http_request_processed", labels)
}

// Monitor returns a handler that sets up MonitoringContext with proper initial values for further usage down the pipeline.
func Monitor(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		cid := r.Header.Get(HDR_CorrelationID)
		if cid == "" {
			cid = uuid.NewV4().String()
		}

		rid := r.Header.Get(HDR_RequestID)
		if rid == "" {
			rid = uuid.NewV4().String()
		}

		next.ServeHTTP(w, r.WithContext(monitoring.WithMonitoringContext(ctx, cid, rid)))
	}
	return http.HandlerFunc(fn)
}

// Logger logs the response statistics by using the Logger from MonitoringContext.
func Logger(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		start := time.Now()
		defer func() {
			mctx, ok := monitoring.GetMonitoringContext(r.Context())
			if ok {
				duration := time.Since(start)

				errType := "None"
				if mctx.HasError() {
					errType = mctx.GetError().ErrorType.String()
				}

				if errType != "None" {
					mctx.Logger().
						With(monitoring.IntLogField(monitoring.LF_HttpStatus, ww.Status()),
							monitoring.Int64LogField(monitoring.LF_Ms, duration.Milliseconds()),
							monitoring.IntLogField(monitoring.LF_Bytes, ww.BytesWritten()),
							monitoring.StrLogField(monitoring.LF_Error, mctx.GetError().Message),
							monitoring.StrLogField(monitoring.LF_ErrorType, errType),
							monitoring.StrLogField(monitoring.LF_ErrorCause, mctx.GetError().Cause)).
						Warn("request failed")
				} else {
					mctx.Logger().
						With(monitoring.IntLogField(monitoring.LF_HttpStatus, ww.Status()),
							monitoring.Int64LogField(monitoring.LF_Ms, duration.Milliseconds()),
							monitoring.IntLogField(monitoring.LF_Bytes, ww.BytesWritten())).
						Info("request completed")
				}

				lv := map[string]string{
					"resource":  mctx.Resource(),
					"operation": mctx.Operation(),
					"status":    strconv.Itoa(ww.Status()),
					"error":     errType,
				}

				durationHist.With(lv).Record(float64(duration.Milliseconds()))
				sizeHist.With(lv).Record(float64(ww.BytesWritten()))
				reqCnt.With(lv).Incr()
			}
		}()

		next.ServeHTTP(ww, r)
	}
	return http.HandlerFunc(fn)
}

// MonitoredHandler wraps the delegate handler to set api resource related values on current MonitoringContext.
func MonitoredHandler(resource string, operationName string, h http.HandlerFunc) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		mctx, ok := monitoring.GetMonitoringContext(r.Context())
		if ok {
			mctx.SetResource(resource)
			mctx.SetOperation(operationName)
		}

		h(w, r)
	}
	return http.HandlerFunc(fn)
}
