package main

import (
	"context"
	"crypto/tls"
	"flag"
	"forum-app/app"
	"forum-app/database"
	"forum-app/helpers/envutil"
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
	// Parse flags
	addr := flag.String("addr", ":8080", "HTTP network address")
	dbName := flag.String("db", "app.db", "Database file name sqlite3")
	flag.Parse()

	// Check and create uploads directory
	err := ensureUploadsDir()
	if err != nil {
		panic(err)
	}

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

	err = envutil.LoadEnv(".env")

	if err != nil {
		logger.Error("Application runtime error", "error", err)
		os.Exit(1)
		return
	}

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
	// Cipher Suite explanation:
	// ECDHE for key exchange
	// RSA for authentication
	// AES-128-GCM for encryption
	// SHA256 for integrity check
	minVersion := tls.VersionTLS12
	var tlsConfig *tls.Config
	var httpHandler http.Handler = routes.Web(app)
	// TLS configuration based on environment (production vs development)
	if envutil.GetEnvString("IS_PRODUCTION") == "true" {
		cacheDir := envutil.GetEnvString("CACHE_DIR")
		hostnames := strings.Split(envutil.GetEnvString("HOSTNAMES"), ",")
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
		certFile := envutil.GetEnvString("CERT_FILE")
		keyFile := envutil.GetEnvString("KEY_FILE")
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
		Addr:      ":8443",
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

func ensureUploadsDir() error {
	const uploadsDir = "uploads"
	if _, err := os.Stat(uploadsDir); os.IsNotExist(err) {
		err := os.Mkdir(uploadsDir, 0755)
		if err != nil {
			return err
		}
	}
	return nil
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

// redirect handles HTTP to HTTPS redirects.
// It's used for both development and production environments.
func redirect(w http.ResponseWriter, r *http.Request) {
	target := "https://" + r.Host + r.URL.RequestURI()
	http.Redirect(w, r, target, http.StatusMovedPermanently)
}
