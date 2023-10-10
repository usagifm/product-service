package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/usagifm/product-service/config"
	"github.com/usagifm/product-service/dependency"
	"github.com/usagifm/product-service/middleware/request"
	"github.com/usagifm/product-service/utils/logger"
)

func StartAPIServer(ctx context.Context) {
	address := fmt.Sprintf(":%d", config.Get().BindAddress)
	logger.GetLogger(ctx).Infof("Starting server on %s", address)

	r := chi.NewRouter()
	r.Use(chimiddleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Use(request.RequestIDContext(request.DefaultGenerator))
	r.Use(request.RequestAttributesContext)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Timeout(config.Get().APITimeout))

	deps := dependency.InstantiateAPIDependencies(ctx)

	Router(r, deps)

	server := &http.Server{Addr: address, Handler: r}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.GetLogger(ctx).Info("Server Started")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.GetLogger(ctx).Fatalf("Server err: %s\n", err)
		}
	}()

	<-done
	logger.GetLogger(ctx).Info("Stopping Server")

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer func() {
		cancel()
	}()

	if err := server.Shutdown(ctx); err != nil {
		logger.GetLogger(ctx).Fatalf("Server Shutdown Failed:%+v", err)
	}
	logger.GetLogger(ctx).Info("Server Exited Properly")
}
