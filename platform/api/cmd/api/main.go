package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"eisparser/pkg/collect"
	"zakupkiplatform/internal/db"
	"zakupkiplatform/internal/httpapi"
	"zakupkiplatform/internal/ingest"
	"zakupkiplatform/internal/sources"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	pool, err := db.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	store := db.NewStore(pool)
	if n, err := store.RequeueStuckRunning(ctx); err != nil {
		log.Printf("requeue stuck: %v", err)
	} else if n > 0 {
		log.Printf("requeued %d stuck running ingest items", n)
	}
	srv := httpapi.New(store)

	col := collect.New(true)
	if col.Client != nil && col.Client.HTTP != nil {
		col.Client.HTTP.Timeout = 45 * time.Second
	}
	pipe := &sources.Pipeline{EIS: col}
	w := &ingest.Worker{Store: store, Pipeline: pipe, Log: log.Default()}
	go w.Run(ctx)

	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	httpSrv := &http.Server{
		Addr:              addr,
		Handler:           httpapi.WithCORS(srv.Mux),
		ReadHeaderTimeout: 10 * time.Second,
	}
	go func() {
		log.Printf("api listening on %s", addr)
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()
	shutdownCtx, c := context.WithTimeout(context.Background(), 10*time.Second)
	defer c()
	_ = httpSrv.Shutdown(shutdownCtx)
}
