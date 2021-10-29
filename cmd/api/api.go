package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const version = "1.0.0"

// config a structure which contains data about our app
type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
	stripe struct {
		secret string
		key    string
	}
}

// application is structure which contains config and additional configurations for our app
type application struct {
	config     config
	infoLogger *log.Logger
	errLogger  *log.Logger
	version    string
}

// serve configure & run server and multiplexer
func (app *application) serve() error {
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", app.config.port),
		Handler:           app.routes(),
		TLSConfig:         nil,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       20 * time.Second,
	}

	app.infoLogger.Println("Starting Back-End Server on port ", app.config.env, " : ", app.config.port)

	sigC := make(chan os.Signal)
	signal.Notify(sigC, os.Interrupt)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			app.errLogger.Fatal(err.Error())
			return
		}
	}()

	// release the program and return
	<-sigC

	return nil
}

func main() {
	cfg := config{}

	flag.IntVar(&cfg.port, "port", 4000, "Server Port to Listen On")
	flag.StringVar(&cfg.env, "env", "development", "Application Environment {development|production|maintenance}")

	flag.Parse()

	cfg.stripe.secret = os.Getenv("STRIPE_SECRET")
	cfg.stripe.key = os.Getenv("STRIPE_KEY")

	infoLogger := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errLogger := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		config:     cfg,
		infoLogger: infoLogger,
		errLogger:  errLogger,
		version:    version,
	}

	_ = app.serve()

	return
}
