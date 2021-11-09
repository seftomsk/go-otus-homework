package main

import (
	"context"
	"errors"
	"flag"
	"io"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	"github.com/seftomsk/go-otus-homework/hw12_13_14_15_calendar/cmd"
	"github.com/seftomsk/go-otus-homework/hw12_13_14_15_calendar/internal/app"
	"github.com/seftomsk/go-otus-homework/hw12_13_14_15_calendar/internal/logger"
	"github.com/seftomsk/go-otus-homework/hw12_13_14_15_calendar/internal/server/web"
	"github.com/seftomsk/go-otus-homework/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/seftomsk/go-otus-homework/hw12_13_14_15_calendar/internal/storage/postgres"
	_ "github.com/seftomsk/go-otus-homework/hw12_13_14_15_calendar/migrations"
)

func main() {
	cfg := cmd.Execute()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	logsPath := "logs"
	if _, err := os.Stat(logsPath); errors.Is(err, os.ErrNotExist) {
		err = os.Mkdir(logsPath, os.ModePerm)
		if err != nil {
			_, _ = os.Stderr.Write([]byte(err.Error()))
			os.Exit(1)
		}
	}
	f, err := os.OpenFile(path.Join(logsPath, cfg.Logger.FileName), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o666)
	if err != nil {
		_, _ = os.Stderr.Write([]byte(err.Error()))
		os.Exit(1)
	}
	defer func() {
		_ = f.Close()
	}()

	var writers io.Writer
	writers = f

	if cfg.Logger.StdErr {
		writers = io.MultiWriter(writers, os.Stderr)
	}

	logg := logger.New(cfg.Logger.LogLevel, writers)

	var repository app.Storage
	if cfg.Storage == "memory" {
		repository = memory.New()
	}
	if cfg.Storage == "postgres" {
		repository, err = sqlstorage.New(cfg.Database.Username, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)
		if err != nil {
			logg.Error("failed to create postgres storage: " + err.Error())
		}
	}

	calendar := app.New(logg, repository)

	server := web.NewServer(logg, calendar, cfg)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop web server: " + err.Error())
		}

		postgres, ok := repository.(*sqlstorage.InPostgres)
		if ok {
			if err := postgres.Close(); err != nil {
				logg.Error("failed to close database: " + err.Error())
			}
		}
	}()

	logg.Info("calendar is running on port: " + cfg.Server.Port + "...")

	postgres, ok := repository.(*sqlstorage.InPostgres)
	if ok {
		if err = postgres.Connect(ctx); err != nil {
			logg.Error("failed to connect to database: " + err.Error())
			cancel()
		}
	}

	if err = server.Start(ctx); err != nil {
		logg.Error("failed to start web server: " + err.Error())
		cancel()
	}
}
