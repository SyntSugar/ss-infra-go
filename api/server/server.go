package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/SyntSugar/ss-infra-go/api/server/handlers"
	"github.com/SyntSugar/ss-infra-go/api/server/middleware"
	"github.com/SyntSugar/ss-infra-go/log"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"go.opentelemetry.io/otel"
)

type Server struct {
	config       *Config
	logger       *log.Logger
	accessLogger *middleware.AccessLogger

	apiEngine   *gin.Engine
	adminEngine *gin.Engine
	apiServer   *http.Server
	adminServer *http.Server
}

// New would create server which contains api and admin api server
func New(config *Config, logger *log.Logger) (*Server, error) {
	if config == nil {
		config = DefaultConfig()
	}
	config.init()
	if err := config.Validate(); err != nil {
		return nil, err
	}
	srv := &Server{
		config: config,
		logger: logger,
	}
	if err := srv.setup(); err != nil {
		return nil, err
	}
	return srv, nil
}

func (srv *Server) setup() error {
	gin.SetMode(gin.ReleaseMode)
	if srv.config.API != nil {
		srv.apiEngine = gin.New()
		srv.apiServer = &http.Server{
			Addr:    srv.config.API.Addr,
			Handler: srv.apiEngine,
		}
		srv.setupAPIDefaultHandlers()
	}
	if srv.config.Admin != nil {
		srv.adminEngine = gin.New()
		srv.adminServer = &http.Server{
			Addr:    srv.config.Admin.Addr,
			Handler: srv.adminEngine,
		}
		srv.setupAdminDefaultHandlers()
	}
	return srv.setupMiddlewares()
}

func (srv *Server) setupMiddlewares() error {
	var err error

	accessLogger, err := middleware.NewAccessLogger(os.Stdout, srv.config.AccessLog.Pattern)
	if err != nil {
		return err
	}
	srv.accessLogger = accessLogger
	if srv.config.AccessLog.Enabled {
		srv.accessLogger.Enabled()
	}

	if srv.apiEngine != nil {
		srv.apiEngine.Use(
			middleware.CORSMiddleware(),
			middleware.DynamicDebugLogging,
			middleware.PanicRecovery(srv.logger),
			middleware.AccessLog(accessLogger),
			middleware.CollectMetrics,
		)

		if srv.config.OpenTelemetry != nil {
			srv.apiEngine.Use(middleware.NewOpenTelemetryTracing(
				srv.config.OpenTelemetry.Exporter,
				otel.GetTextMapPropagator(),
				otel.GetTracerProvider()))
		}

		srv.apiEngine.NoRoute(func(c *gin.Context) {
			c.String(http.StatusNotFound, "not found")
		})
	}

	if srv.adminEngine != nil {
		srv.adminEngine.Use(
			middleware.CORSMiddleware(),
			middleware.DynamicDebugLogging,
			middleware.PanicRecovery(srv.logger),
			middleware.AccessLog(accessLogger),
		)
	}

	return nil
}

func (srv *Server) setupAPIDefaultHandlers() {
	srv.apiEngine.NoRoute(handlers.NoRoute)
	srv.apiEngine.NoMethod(handlers.NoMethod)
	srv.apiEngine.GET(srv.config.API.BasePath+"/whoami", handlers.Whoami)
}

func (srv *Server) setupAdminDefaultHandlers() {
	srv.adminEngine.NoRoute(handlers.NoRoute)
	srv.adminEngine.NoMethod(handlers.NoMethod)
	statusGroup := srv.adminEngine.Group("/devops/status")
	{
		statusGroup.HEAD("", handlers.GetDevopsStatus)
		statusGroup.GET("", handlers.GetDevopsStatus)
		statusGroup.POST("", handlers.UpdateDevopsStatus)
		statusGroup.PUT("", handlers.UpdateDevopsStatus)
	}
	accessLog := srv.adminEngine.Group("/access_log")
	{
		accessLog.GET("/status", func(c *gin.Context) {
			c.String(http.StatusOK, srv.accessLogger.Status())
		})
		accessLog.POST("/status/:status", func(c *gin.Context) {
			if strings.ToLower(c.Param("status")) == "disabled" {
				srv.accessLogger.Disabled()
			} else {
				srv.accessLogger.Enabled()
			}
			c.String(http.StatusOK, srv.accessLogger.Status())
		})
		accessLog.POST("/slow_request_log/threshold/:duration", func(c *gin.Context) {
			if d, err := strconv.Atoi(c.Param("duration")); err == nil {
				srv.accessLogger.SetSlowRequestThreshold(time.Duration(d) * time.Millisecond)
				c.String(http.StatusOK, fmt.Sprintf("duration threshold of slow request's access_log is updated, new threshold is %d milliseconds.", d))
			}
			c.String(http.StatusUnprocessableEntity, "duration threshold(ms) param is invalid.")
		})
	}
	srv.adminEngine.GET(srv.config.Admin.BasePath+"/whoami", handlers.Whoami)
	srv.adminEngine.Any("/debug/pprof/*profile", handlers.PProf)
	srv.adminEngine.GET("/metrics", gin.WrapH(promhttp.Handler()))
}

// GetAPIRouteGroup return api's gin engine that user can add api handlers
func (srv *Server) GetAPIRouteGroup() *gin.RouterGroup {
	return srv.apiEngine.Group(srv.config.API.BasePath)
}

func (srv *Server) GetAPIEngine() *gin.Engine {
	return srv.apiEngine
}

func (srv *Server) GetAdminEngine() *gin.Engine {
	return srv.adminEngine
}

// GetAdminRouteGroup return admin's gin engine that user can add admin handlers
func (srv *Server) GetAdminRouteGroup() *gin.RouterGroup {
	return srv.adminEngine.Group("")
}

// Run would setup api/admin api server and async listening on api ports
func (srv *Server) Run() error {
	for _, server := range []*http.Server{srv.apiServer, srv.adminServer} {
		go func(httpServer *http.Server) {
			if httpServer == nil {
				return
			}
			if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				fatalf("Failed to setup api httpServer, err: %s\n", err.Error())
			}
		}(server)
	}
	return nil
}

// ServeHTTP uses to parse request from http.Server
func (srv *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	srv.apiEngine.ServeHTTP(w, r)
}

// Shutdown would graceful shutdown the api server
func (srv *Server) Shutdown() error {
	// Admin server does not need to be stopped gracefully
	srv.adminServer.Close()
	if srv.apiServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), srv.config.Shutdown)
		defer cancel()
		return srv.apiServer.Shutdown(ctx)
	}
	return nil
}

func fatalf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
	os.Exit(1)
}
