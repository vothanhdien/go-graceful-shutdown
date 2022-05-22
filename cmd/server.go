package cmd

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/robfig/cron/v3"
)

func NewService() *Service {
	return &Service{
		// allocates 0 byte
		running: make(chan interface{}),
	}
}

type Service struct {
	// running channel prevents service exist unexpected
	running chan interface{}

	c    *cron.Cron
	http *http.Server
}

func (s *Service) Start() {
	quit := make(chan os.Signal, 1) // we need to reserve to buffer size 1, so the notifier are not blocked
	// syscall.SIGINT on Ctrl+C
	// syscall.SIGTERM is the usual signal for termination on docker containers, which is also used by kubernetes.
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		defer func() {
			// now service can stop
			close(s.running)
		}()

		s.stop()
	}()

	// start http server
	s.startHTTP()

	// start cron
	s.startCron()

	<-s.running
	log.Println("service stopped")
}

func (s *Service) stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.http.Shutdown(ctx); err != nil {
		log.Printf("http close error %v\n", err)
	}

	// don't need to separate the reason
	reasonCareless, cronCancel := context.WithTimeout(s.c.Stop(), 30*time.Second)
	defer cronCancel()
	<-reasonCareless.Done()
}

func (s *Service) startHTTP() {
	s.http = &http.Server{
		Addr: ":8080",
		Handler: http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			time.Sleep(30 * time.Second)
			writer.WriteHeader(http.StatusOK)
			_, _ = writer.Write([]byte("hello"))
		}),
	}

	// start goroutine avoid blocking
	go func() {
		log.Printf("server listening %s\n", s.http.Addr)
		if err := s.http.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server closed unexpect %v", err)
		}
	}()
}

func (s *Service) startCron() {
	s.c = cron.New()
	_, err := s.c.AddFunc("* * * * *", func() {
		t := time.Now()
		log.Printf("now : %v\n", t.Format(time.RFC3339))
	})
	if err != nil {
		panic(err)
	}

	// start goroutine avoid blocking
	go func() {
		log.Println("cron starting")
		s.c.Run()
	}()
}
