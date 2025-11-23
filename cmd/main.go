// Package main provides ZenMoney backup application.
package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/egregors/zenmoney-backup/notifier"
	"github.com/egregors/zenmoney-backup/srv"
	"github.com/egregors/zenmoney-backup/store"
	log "github.com/go-pkgz/lgr"
	"github.com/jessevdk/go-flags"
)

// Opts is App settings (from cli args or ENV).
type Opts struct {
	Token     string `short:"t" long:"zenmoney OAuth token" env:"ZEN_TOKEN" description:"Zenmoney API Token, to get it visit: https://zerro.app/token"`
	SleepTime string `short:"p" long:"sleep_time" env:"SLEEP_TIME" default:"24h" description:"Backup performs every SLEEP_TIME minutes"`
	Timeout   int    `short:"c" long:"timeout" env:"TIMEOUT" default:"10" description:"Backup request timeout in seconds"`
	NotifyURL string `short:"n" long:"notify_url" env:"NOTIFY_URL" description:"ntfy.sh notification URL (e.g., https://ntfy.sh/your_topic)"`

	Dbg bool `long:"dbg" env:"DEBUG" description:"Debug mode"`
}

var revision = "unknown"

func main() {
	fmt.Printf("zenmoney-backup %s\n~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=[,,_,,]:3\n", revision)

	var opts Opts
	p := flags.NewParser(&opts, flags.PrintErrors|flags.PassDoubleDash|flags.HelpFlag)
	if _, err := p.Parse(); err != nil {
		var flagsErr *flags.Error
		if errors.As(err, &flagsErr) && flagsErr.Type != flags.ErrHelp {
			log.Printf("[ERROR] cli error: %v", err)
		}
		os.Exit(2)
	}

	setupLog(opts.Dbg)

	s, err := makeServer(opts)
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

	s.Run(ctx)
}

func setupLog(dbg bool) {
	if dbg {
		log.Setup(log.Debug, log.CallerFile, log.CallerFunc, log.Msec, log.LevelBraces)
		return
	}
	log.Setup(log.Msec, log.LevelBraces)
}

func makeServer(opts Opts) (*srv.Server, error) {
	d, err := time.ParseDuration(opts.SleepTime)
	if err != nil {
		return nil, err
	}
	
	if opts.Timeout <= 0 {
		return nil, fmt.Errorf("timeout must be a positive integer, got %d", opts.Timeout)
	}
	
	timeout := time.Duration(opts.Timeout) * time.Second
	
	// Create notifier
	var n srv.Notifier
	if opts.NotifyURL != "" {
		n = notifier.NewNtfy(opts.NotifyURL)
		log.Printf("[INFO] notifications enabled for URL: %s", opts.NotifyURL)
	} else {
		n = notifier.NewNoop()
	}
	
	return srv.NewServer(opts.Token, d, timeout, store.LocalFs{}, n), nil
}
