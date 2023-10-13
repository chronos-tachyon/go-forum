package main

import (
	"context"
	"flag"
	"io/fs"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/chronos-tachyon/go-autolog"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"

	"github.com/chronos-tachyon/go-forum/data"
)

var (
	gBuiltInFS fs.FS  = data.FS()
	gDataFS    fs.FS  = nil
	gDataDir   string = "<built-in>"
	gConfig    Config
	gServer    http.Server
	gRouter    *mux.Router
)

func main() {
	autolog.Init()
	defer func() {
		err := autolog.Done()
		if err != nil {
			panic(err)
		}
	}()

	ctx := context.Background()

	var (
		flagListenApp     string
		flagListenMetrics string
		flagDataDir       string
	)

	flag.StringVar(&flagListenApp, "listen-app", ":8080", "ip:port to listen on (application)")
	flag.StringVar(&flagListenMetrics, "listen-mon", ":8081", "ip:port to listen on (monitoring/metrics)")
	flag.StringVar(&flagDataDir, "data-dir", "", "path to directory containing custom HTML templates and static files")
	flag.Parse()

	if flagDataDir != "" {
		gDataDir = filepath.Clean(flagDataDir)
		gDataFS = os.DirFS(gDataDir)
	}

	if err := gConfig.Load(gDataDir, gDataFS); err != nil {
		log.Logger.Fatal().
			Err(err).
			Interface("fs", gDataDir).
			Msg("failed to load templates")
		panic(nil)
	}

	gRouter = mux.NewRouter().StrictSlash(true)
	AddRoutes(gRouter)

	gServer = http.Server{
		Addr:           flagListenApp,
		Handler:        gRouter,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Logger.Info().
		Str("listen", flagListenApp).
		Msg("starting server")

	reloadCh := make(chan os.Signal, 1)
	shutdownCh := make(chan os.Signal, 1)
	go func() {
		for range reloadCh {
			OnReload(ctx)
		}
	}()
	go func() {
		<-shutdownCh

		signal.Stop(reloadCh)
		close(reloadCh)

		signal.Stop(shutdownCh)
		close(shutdownCh)

		inner, cancel := context.WithTimeout(ctx, 15*time.Second)
		defer cancel()
		OnShutdown(inner)
	}()

	signal.Notify(reloadCh, reloadSignals...)
	signal.Notify(shutdownCh, shutdownSignals...)

	err := gServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Logger.Fatal().
			Err(err).
			Msg("failed to start")
		panic(nil)
	}

	log.Logger.Info().Msg("OK")
}

func mustSubFS(outer fs.FS, dir string) fs.FS {
	inner, err := fs.Sub(outer, dir)
	if err != nil {
		panic(err)
	}
	return inner
}

func OnReload(ctx context.Context) {
	log.Logger.Info().Msg("reloading")

	if err := autolog.Rotate(); err != nil {
		log.Logger.Error().
			Err(err).
			Msg("failed to rotate log file")
	}
}

func OnShutdown(ctx context.Context) {
	log.Logger.Info().Msg("shutting down")

	err := gServer.Shutdown(ctx)
	if err != nil {
		log.Logger.Error().
			Err(err).
			Msg("failed to shut down server")
	}
}

func T(key string) string {
	return gConfig.Translate(key)
}

var (
	shutdownSignals = []os.Signal{os.Interrupt, syscall.SIGINT, syscall.SIGTERM}
	reloadSignals   = []os.Signal{syscall.SIGHUP}
)
