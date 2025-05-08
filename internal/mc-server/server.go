package server

import (
	"context"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/candbright/go-log/log"
	"github.com/candbright/go-server/internal/mc-server/route"
	"github.com/candbright/go-server/pkg/config"
	"github.com/candbright/go-server/pkg/rest/handler"
	"github.com/gin-gonic/gin"
)

type Server struct {
	g sync.WaitGroup
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Serve() {
	ctx, cancel := context.WithCancel(context.Background())
	s.StartHTTPServer()
	s.StartWatcher(ctx)
	s.WaitSignal()
	cancel()
	s.g.Wait()
}

func (s *Server) StartWatcher(ctx context.Context) {
	log.Info("Starting Watcher...")
	s.g.Add(1)
	go func() {
		defer s.g.Done()
		testTicker := time.NewTicker(time.Hour * 1)
		defer testTicker.Stop()

		for {
			select {
			case <-testTicker.C:
				testLoop()
			case <-ctx.Done():
				log.Info("Watcher context cancelled, shutting down...")
				return
			}
		}
	}()
}

func (s *Server) StartHTTPServer() {
	log.Info("Start HTTP Server...")
	go func() {
		engine := gin.New()
		engine.Use(gin.BasicAuth(
			map[string]string{
				config.Global.Get("server.username"): config.Global.Get("server.password"),
			},
		))
		engine.Use(handler.LogHandler())
		engine.Use(gin.Recovery())

		route.Init()
		route.Incubate(engine)
		_ = engine.Run(":" + strconv.Itoa(config.Global.GetInt("server.port")))
		log.Warn("Exit HTTP server!")
	}()
}

func (s *Server) WaitSignal() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	for {
		sig := <-signals
		switch sig {
		case syscall.SIGINT:
			log.Warn("Received SIGINT: Initiating graceful shutdown...")
			return
		case syscall.SIGTERM:
			log.Warn("Received SIGTERM: Initiating graceful shutdown...")
			return
		default:
			log.Warnf("Received unexpected signal: %v", sig)
		}
	}
}
