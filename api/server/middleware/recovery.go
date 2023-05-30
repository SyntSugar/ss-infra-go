package middleware

import (
	"fmt"

	"net/http"
	"runtime/debug"

	rsp "github.com/SyntSugar/ss-infra-go/api/response"
	internal_metrics "github.com/SyntSugar/ss-infra-go/internal"
	"github.com/SyntSugar/ss-infra-go/log"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"

	"go.uber.org/zap"
)

// PanicRecovery uses to catch panic in api handler
func PanicRecovery(logger *log.Logger) gin.HandlerFunc {
	logError := logger != nil
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				internal_metrics.Get().Panic.With(prometheus.Labels{"type": "http"}).Inc()
				if logError {
					logger.Error("Catch panic",
						zap.Error(fmt.Errorf("panic error: %v", err)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					fmt.Printf("Catch panic: %v with stack:\n%s\n", err, string(debug.Stack()))
				}
				rsp.ResponseWithErrors(c, http.StatusInternalServerError, 0, nil)
				c.Abort()
			}
		}()
		c.Next()
	}
}
