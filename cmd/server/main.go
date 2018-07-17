package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/hashicorp/logutils"
	flags "github.com/jessevdk/go-flags"
	"github.com/mbobakov/chat/api"
	"github.com/mbobakov/chat/internal/debug"
	"golang.org/x/sync/errgroup"
)

func main() {
	var opts = struct {
		GRPCListen  string `long:"grpc.listen" env:"BP_GRPC_LISTEN" default:":8080" description:"GRPC server interface"`
		DebugListen string `long:"debug.listen" env:"BP_DEBUG_LISTEN" default:":6060" description:"Interface for serve debug information(metrics/health/pprof)"`
		Verbose     bool   `long:"v" env:"BP_VERBOSE" description:"Enable Verbose log  output"`
	}{}

	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "INFO", "WARNING", "ERROR"},
		MinLevel: logutils.LogLevel("INFO"),
		Writer:   os.Stdout,
	}

	if opts.Verbose {
		filter.SetMinLevel(logutils.LogLevel("DEBUG"))
	}

	logger := log.New(filter, "", log.Ldate|log.Ltime|log.LUTC|log.Lshortfile)
	logger.Printf("[INFO] Launching Application with: %+v", opts)

	gr, ctx := errgroup.WithContext(context.Background())
	_, err := flags.Parse(&opts)
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		}
		os.Exit(1)
	}
	d := debug.New()
	gr.Go(func() error {
		return d.Serve(
			ctx,
			opts.DebugListen,
			logger,
		)
	})
	gr.Go(func() error {
		d.SetHealthOK()
		defer d.SetHealthNotOK()
		return api.Serve(
			ctx,
			opts.GRPCListen,
			api.WithLogger(logger),
		)
	})

	ErrCanceled := errors.New("Canceled")
	gr.Go(func() error {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		select {
		case <-ctx.Done():
			return nil
		case <-sigs:
			logger.Printf("[INFO] Caught stop signal. Exiting ...")
			return ErrCanceled
		}
	})

	err = gr.Wait()
	if err != nil && err != ErrCanceled {
		logger.Fatal(err)
	}
}
