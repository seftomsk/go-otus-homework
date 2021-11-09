package web

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/seftomsk/go-otus-homework/hw12_13_14_15_calendar/internal/app"
)

type Server struct {
	host   string
	port   string
	app    Application
	logger Logger
	server *http.Server
}

type Logger interface {
	Info(msg string)
	Error(msg string)
	Warn(msg string)
	Debug(msg string)
}

type Application interface {
	CreateEvent(ctx context.Context, dto app.EventDTO) (app.Event, error)
	UpdateEvent(ctx context.Context, dto app.UpdateEventDTO) (app.Event, error)
	DeleteEvent(ctx context.Context, id string) error
	EventsForDay(ctx context.Context, date time.Time) (map[string]app.Event, error)
	EventsForWeek(ctx context.Context, date time.Time) (map[string]app.Event, error)
	EventsForMonth(ctx context.Context, date time.Time) (map[string]app.Event, error)
}

type ServerConfig interface {
	Host() string
	Port() string
}

func NewServer(logger Logger, app Application, cfg ServerConfig) *Server {
	return &Server{
		host:   cfg.Host(),
		port:   cfg.Port(),
		app:    app,
		logger: logger,
	}
}

func indexHandler(_ http.ResponseWriter, _ *http.Request) {
	time.Sleep(time.Second)
}

func (s *Server) Start(ctx context.Context) error {
	var err error

	go func() {
		server := &http.Server{
			Addr:    net.JoinHostPort(s.host, s.port),
			Handler: nil,
		}
		s.server = server

		http.Handle("/", loggerMiddleware(indexHandler, s.logger))

		err = server.ListenAndServe()
	}()

	<-ctx.Done()

	if err != nil {
		return err
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if s.server == nil {
		return nil
	}

	if err := s.server.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}
