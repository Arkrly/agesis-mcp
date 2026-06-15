package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yourorg/aegis-mcp/internal/api"
	"github.com/yourorg/aegis-mcp/internal/config"
	"github.com/yourorg/aegis-mcp/internal/runtime"
)

func main() {
	cfg, err := config.LoadFromEnv()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	deps, err := runtime.Build(cfg)
	if err != nil {
		log.Fatalf("build runtime: %v", err)
	}
	defer deps.Audit.Close()

	handler := api.NewServer(cfg, deps.Logger, deps.Metrics, deps.Audit, deps.Decision, deps.Upstream, deps.Policy, deps.Inspector)

	srv := &http.Server{
		Addr:              cfg.ListenAddr,
		Handler:           handler,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		ReadTimeout:       cfg.ReadTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
	}

	go func() {
		deps.Logger.Info("server listening", map[string]any{
			"listen_addr": cfg.ListenAddr,
			"upstream":    cfg.UpstreamURL,
		})
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		deps.Logger.Error("shutdown error", map[string]any{"error": err.Error()})
	}
}
