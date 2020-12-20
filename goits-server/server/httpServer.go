// Package server contains the main http server entry point of the goits application.
package server

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cpekyaman/goits/config"
	"github.com/cpekyaman/goits/framework/monitoring"
	"github.com/cpekyaman/goits/framework/routing"

	"github.com/cpekyaman/goits/application/project"

	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type httpConfig struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"readTimeout"`
	WriteTimeout time.Duration `mapstructure:"writeTimeout"`
}

var conf httpConfig

// Start configures the router and starts the http server.
func Start() {
	// routing engine
	routing.InitRouting()
	routing.Engine().RegisterPath("/metrics", promhttp.Handler())

	// individual routers
	project.InitProject()

	config.ReadInto("http", &conf)
	svc := createServer(routing.Engine().Router())

	startServer(svc)
	handleGracefulShutdown(svc)
}

// createServer creates a new http server by using given router as the handler.
func createServer(r http.Handler) *http.Server {
	return &http.Server{
		Addr:           fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		Handler:        r,
		ReadTimeout:    conf.ReadTimeout * time.Second,
		WriteTimeout:   conf.WriteTimeout * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}

// startServer actually starts the server.
func startServer(svc *http.Server) {
	go func() {
		monitoring.RootLogger().Info("starting server")
		if err := svc.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			monitoring.RootLogger().With(monitoring.ErrLogField(err)).Fatal("could not start server")
		}
		monitoring.RootLogger().Info("server started")
	}()
}

// handleGracefulShutdown registers necessary signal handlers to handle graceful shutdown with term / kill.
func handleGracefulShutdown(svc *http.Server) {
	ch := make(chan os.Signal)

	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	monitoring.RootLogger().Info("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := svc.Shutdown(ctx); err != nil {
		monitoring.RootLogger().With(monitoring.ErrLogField(err)).Fatal("failed to shutdown grafully")
	}
	monitoring.RootLogger().Info("server exiting")
}
