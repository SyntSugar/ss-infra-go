package middleware

import (
	"context"

	"github.com/SyntSugar/ss-infra-go/consts"
	"github.com/gin-gonic/gin"
)

// DynamicDebugLogging open the debug level logging in dynamically
func DynamicDebugLogging(c *gin.Context) {
	if c.GetHeader(consts.HeaderEnableDebugLogging) != "" {
		c.Set(string(consts.ContextKeyEnableDebugLogging), true)
		ctx := context.WithValue(c.Request.Context(), consts.ContextKeyEnableDebugLogging, true)
		c.Request = c.Request.WithContext(ctx)
	}
	c.Next()
}
