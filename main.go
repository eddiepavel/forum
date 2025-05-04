package main

import (
	"context"
	"crypto/tls"
	"flag"
	"forum-app/app"
	"forum-app/database"
	"forum-app/environment"
	"forum-app/ratelimiter"
	"forum-app/routes"
	"forum-app/session"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"golang.org/x/crypto/acme/autocert"
)

func main() {

	// first, load .env if present
	if err := environment.LoadDotEnv(".env"); err != nil {
		log.Fatalf("loading .env: %v", err)
	}

	// Read the  config from os.Getenv
	isProd := os.Getenv("IS_PRODUCTION") == "true"
	certFile := os.Getenv("CERT_FILE")
	keyFile := os.Getenv("KEY_FILE")
	cacheDir := os.Getenv("CERT_CACHE_DIR")
	hostnames := strings.Split(os.Getenv("HOSTNAMES"), ",") // e.g. "example.com,www.example.com"  replace with real domain

	// Command line flags for configuration
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

	// Session manager with 1-hour timeout for both session and idle time
	session := session.NewSessionStore(1*time.Hour, 1*time.Hour)

	// Rate limiter allowing 100 requests per minute per client
	rl := ratelimiter.NewRateLimiter(100, 1*time.Minute)

	// Application context holding all dependencies
	app := &app.Application{
		DB:          db,
		Logger:      logger,
		Session:     session,
		RateLimiter: rl,
	}

	// Security configuration for TLS
	cipherSuites := []uint16{
		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
	}
	// Cipher Suite explaination:
	// ECDHE for key exchange
	// RSA for authentication
	// AES-128-GCM for encryption
	// SHA256 for integrity check
	minVersion := tls.VersionTLS12
	var tlsConfig *tls.Config
	var httpHandler http.Handler = routes.Web(app)
	// TLS configuration based on environment (production vs development)
	if isProd {
		// ==== Production: autocert + Let's Encrypt ====
		m := &autocert.Manager{
			Cache:      autocert.DirCache(cacheDir), // cert cache dir
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(hostnames...),
		}

		// Handle HTTP challenges for Let's Encrypt and redirect HTTP to HTTPS
		go func() {
			srv := &http.Server{
				Addr:    ":80",
				Handler: m.HTTPHandler(http.HandlerFunc(redirect)),
			}
			log.Fatal(srv.ListenAndServe())
		}()

		tlsConfig = &tls.Config{
			GetCertificate: m.GetCertificate,
			MinVersion:     uint16(minVersion),
			CipherSuites:   cipherSuites,
		}

	} else {
		// ==== Development: self-signed cert.pem/key.pem ====
		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			log.Fatalf("loading self-signed cert: %v", err)
		}
		tlsConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
			MinVersion:   uint16(minVersion),
			CipherSuites: cipherSuites,
		}
	}

	// Configure and start the HTTPS server
	server := &http.Server{
		Addr:      ":443",
		TLSConfig: tlsConfig,
		Handler:   httpHandler,
	}
	// Start HTTPS server in a separate goroutine
	go func() {
		logger.Info("starting server", "addr", *addr)
		if err := server.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
			logger.Error("HTTP server error", "error", err)
			os.Exit(1)
		}
	}()

	// Start HTTP redirect server (redirects all HTTP traffic to HTTPS)
	go func() {
		http.ListenAndServe(":80", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "https://"+r.Host+r.URL.String(), http.StatusMovedPermanently)
		}))
	}()
	// Graceful shutdown
	waitForShutdown(server, logger)
}

// initDatabase creates and initializes the database connection.
// It returns a database connection wrapper and any error encountered.

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

// redirect handles HTTP to HTTPS redirects.
// It's used for both development and production environments.
func redirect(w http.ResponseWriter, r *http.Request) {
	target := "https://" + r.Host + r.URL.RequestURI()
	http.Redirect(w, r, target, http.StatusMovedPermanently)
}
