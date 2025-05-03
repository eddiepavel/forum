package main

import (
	"context"
	"crypto/tls"
	"flag"
	"forum-app/app"
	"forum-app/database"
	"forum-app/ratelimiter"
	"forum-app/routes"
	"forum-app/session"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Parse flags
	addr := flag.String("addr", ":8080", "HTTP network address")
	dbName := flag.String("db", "app.db", "Database file name sqlite3")
	flag.Parse()

	// Initialize logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Initialize dependencies
	db, err := initDatabase(*dbName, logger)
	if err != nil {
		logger.Error("Database initialization error", "error", err)
		os.Exit(1)
	}
	defer db.DB.Close()

	session := session.NewSessionStore(1*time.Hour, 1*time.Hour)

	rl := ratelimiter.NewRateLimiter(100, 1*time.Minute)

	app := &app.Application{
		DB:          db,
		Logger:      logger,
		Session:     session,
		RateLimiter: rl,
	}
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
	}

	// Create HTTP server
	server := &http.Server{
		Addr:      *addr,
		TLSConfig: tlsConfig,
		Handler:   routes.Web(app),
	}

	// Start server in a goroutine
	go func() {
		logger.Info("starting server", "addr", *addr)
		if err := server.ListenAndServeTLS("cert.pem", "key.pem"); err != nil && err != http.ErrServerClosed {
			logger.Error("HTTP server error", "error", err)
			os.Exit(1)
		}
	}()
	// Added redirect from http to https
	go func() {
		http.ListenAndServe(":80", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "https://"+r.Host+r.URL.String(), http.StatusMovedPermanently)
		}))
	}()
	// Graceful shutdown
	waitForShutdown(server, logger)
}

func initDatabase(dbName string, logger *slog.Logger) (*database.Connection, error) {
	db, err := database.NewConnection(dbName)
	if err != nil {
		return nil, err
	}
	logger.Info("Database connected", "dbName", dbName)
	return db, nil
}

func waitForShutdown(server *http.Server, logger *slog.Logger) {
	// Listen for termination signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	logger.Info("shutting down server")

	// Gracefully shut down the server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server shutdown error", "error", err)
	} else {
		logger.Info("Server stopped gracefully")
	}
}
