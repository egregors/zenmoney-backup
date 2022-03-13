package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/egregors/zenmoney-backup/cmd"
	"github.com/egregors/zenmoney-backup/store"
	log "github.com/go-pkgz/lgr"
	"github.com/jessevdk/go-flags"
)

// Opts is App settings (from cli args or ENV)
type Opts struct {
	ZenUsername string `short:"l" long:"zen_username" env:"ZEN_USERNAME" description:"Your zenmoney username"`
	ZenPassword string `short:"p" long:"zen_password" env:"ZEN_PASSWORD" description:"Your zenmoney password"`
	SleepTime   string `short:"t" long:"sleep_time" env:"SLEEP_TIME" default:"24h" description:"Backup performs every SLEEP_TIME minutes"`

	Dbg bool `long:"dbg" env:"DEBUG" description:"Debug mode"`
}

var revision = "unknown"

func main() {
	fmt.Printf("zenmoney-backup %s\n~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=[,,_,,]:3\n", revision)

	var opts Opts
	p := flags.NewParser(&opts, flags.PrintErrors|flags.PassDoubleDash|flags.HelpFlag)
	if _, err := p.Parse(); err != nil {
		if err.(*flags.Error).Type != flags.ErrHelp {
			log.Printf("[ERROR] cli error: %v", err)
		}
		os.Exit(2)
	}

	setupLog(opts.Dbg)

	srv, err := makeServer(opts)
	if err != nil {
		log.Printf("[FATAL] can't make server: %s", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())

	go func(cancel context.CancelFunc) {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

		<-stop
		cancel()
		log.Printf("[INFO] shutting down")
	}(cancel)

	srv.Run(ctx)
}

func setupLog(dbg bool) {
	if dbg {
		log.Setup(log.Debug, log.CallerFile, log.CallerFunc, log.Msec, log.LevelBraces)
		return
	}
	log.Setup(log.Msec, log.LevelBraces)
}

func makeServer(opts Opts) (*cmd.Server, error) {
	d, err := time.ParseDuration(opts.SleepTime)
	if err != nil {
		return nil, err
	}
	return cmd.NewServer(opts.ZenUsername, opts.ZenPassword, d, store.LocalFs{}), nil
}
