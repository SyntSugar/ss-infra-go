package middleware

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/SyntSugar/ss-infra-go/consts"
	"github.com/gin-gonic/gin"
)

var result any

func setupTest() (*gin.Context, *gin.Engine, *httptest.ResponseRecorder) {
	gin.SetMode(gin.ReleaseMode)
	writer := httptest.NewRecorder()
	ctx, engine := gin.CreateTestContext(writer)
	ctx.Request, _ = http.NewRequest(http.MethodGet, "http://localhost/ping", nil)

	return ctx, engine, writer
}

func BenchmarkCreateLogItem(b *testing.B) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		resp := fmt.Sprintf("<html><body>%s - Hello World!</body></html>", method)
		io.WriteString(w, resp)
	}

	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()
	handler(w, req)
	writer := NewResponseWriter(w)

	// Run the benchmark loop
	b.ResetTimer() // Reset the timer to ignore setup time
	for i := 0; i < b.N; i++ {
		logItem := createLogItem(req, writer, time.Now(), time.Second)
		// Use the logItem in a way that prevents the compiler from optimizing the function call away
		result = logItem
	}
}

func BenchmarkAccessLog(b *testing.B) {
	logger, _ := NewAccessLogger(io.Discard, "")

	handle := func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		resp := fmt.Sprintf("<html><body>%s - Hello World!</body></html>", method)
		io.WriteString(w, resp)
	}

	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()
	handle(w, req)
	writer := NewResponseWriter(w)

	logItem := createLogItem(req, writer, time.Now(), time.Second)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.record(&logItem)
	}
}

// TODO: check the log print right or wrong.
func TestAccessLog(t *testing.T) {
	testData := []struct {
		name      string
		enabled   bool
		threshold time.Duration
		sleep     time.Duration
	}{
		{"DisableAccessLogTraceSlowRequestExceedThreshold", false, 0, 1600 * time.Millisecond},
		{"DisableAccessLogTraceSlowRequestNotExceedThreshold", false, 0, 1000 * time.Millisecond},
		{"EnableAccessLogTraceSlowRequestNotExceedThreshold", true, 0, 1000 * time.Millisecond},
		{"EnableAccessLogTraceSlowRequestExceedThreshold", true, 0, 1700 * time.Millisecond},
		{"AccessLogTraceSlowRequestWithWrongThreshold", false, -300 * time.Millisecond, 300 * time.Millisecond},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			ctx, engine, writer := setupTest()

			// Create logger
			logger, _ := NewAccessLogger(io.Discard, consts.JSONAccessLogPattern)
			if tt.enabled {
				logger.Enabled()
			}
			logger.SetSlowRequestThreshold(tt.threshold)

			// Setup request handler
			engine.Handle(http.MethodGet, "/ping", AccessLog(logger), func(ctx *gin.Context) {
				time.Sleep(tt.sleep)
				ctx.Status(http.StatusOK)
			})

			// Send request
			engine.HandleContext(ctx)

			// Check response
			resp := writer.Result()
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Errorf("expected status OK; got %v", resp.StatusCode)
			}
		})
	}
}
