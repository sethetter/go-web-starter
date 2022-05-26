package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/sethetter/go-web-starter/pkg/config"
	"github.com/sethetter/go-web-starter/pkg/server"
	"github.com/sethetter/go-web-starter/pkg/services"

	_ "github.com/lib/pq"
)

func main() {
	rand.Seed(time.Now().Unix())

	log.SetFlags(log.Flags() | log.Lshortfile)
	log.SetOutput(os.Stderr)

	if err := run(); err != nil {
		log.Fatalln("main failed to run:", err)
	}

	log.Println("sucessful shutdown")
}

func run() error {
	c, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to LoadConfig: %w", err)
	}

	// get our database connection
	db, err := sql.Open("postgres", c.DatabaseURL)
	if err != nil {
		return fmt.Errorf("failed to sqlx.Open: %w", err)
	}

	conf := &server.ServerConfig{
		Config:       c,
		DB:           db,
		TemplatePath: "./templates",
	}

	if c.Email.SMTPHost != "" {
		conf.EmailService = &services.EmailService{Conf: c.Email}
	}

	server, err := server.NewServer(conf)
	if err != nil {
		return fmt.Errorf("failed to create server: %w", err)
	}

	serverErrors := make(chan error, 1)

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Printf("server listening on port %s", server.Addr)
		serverErrors <- server.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case <-serverErrors:
		return fmt.Errorf("received server error: %w", err)

	case sig := <-shutdown:
		log.Printf("received shutdown signal %q", sig)

		if err := server.Shutdown(context.Background()); err != nil {
			return fmt.Errorf("failed to server.Shutdown: %w", err)
		}
	}

	wg.Wait()

	return nil
}
