package server

import (
	"context"
	"crypto/tls"
	"forum-app/app"
	"forum-app/database"
	"forum-app/environment"
	"forum-app/middleware"
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

func LoadEnvironment(envFile string) {
	if err := environment.LoadDotEnv(envFile); err != nil {
		log.Fatalf("loading .env: %v", err)
	}
}

func InitializeApp(dbName string) (*app.Application, http.Handler, http.Handler, *tls.Config) {
	isProd := os.Getenv("IS_PRODUCTION") == "true"
	certFile := os.Getenv("CERT_FILE")
	keyFile := os.Getenv("KEY_FILE")
	cacheDir := os.Getenv("CERT_CACHE_DIR")
	hostnames := strings.Split(os.Getenv("HOSTNAMES"), ",")
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := initDatabase(dbName, logger)
	if err != nil {
		logger.Error("Database initialization error", "error", err)
		os.Exit(1)
	}
	deferFunc := func() { db.DB.Close() }
	_ = deferFunc // For linter. This close should be placed in the real main context.

	sessionStore := session.NewSessionStore(1*time.Hour, 1*time.Hour)
	rl := ratelimiter.NewRateLimiter(100, 1*time.Minute)

	application := &app.Application{
		DB:          db,
		Logger:      logger,
		Session:     sessionStore,
		RateLimiter: rl,
	}

	cipherSuites := []uint16{
		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
	}
	minVersion := tls.VersionTLS12

	tlsConfig, http80Handler := buildTLSConfig(isProd, certFile, keyFile, cacheDir, hostnames, cipherSuites, uint16(minVersion), logger)
	httpHandler := middleware.RecoverMiddleware(routes.Web(application))

	return application, httpHandler, http80Handler, tlsConfig
}

func initDatabase(dbName string, logger *slog.Logger) (*database.Connection, error) {
	db, err := database.NewConnection(dbName)
	if err != nil {
		return nil, err
	}
	logger.Info("Database connected", "dbName", dbName)
	return db, nil
}

func buildTLSConfig(isProd bool, certFile, keyFile, cacheDir string, hostnames []string, cipherSuites []uint16, minVersion uint16, logger *slog.Logger) (*tls.Config, http.Handler) {
	httpsPort := os.Getenv("HTTPS_PORT")
	if httpsPort == "" {
		httpsPort = ":8888" // Default fallback you prefer
	} else if httpsPort[0] != ':' {
		httpsPort = ":" + httpsPort
	}

	redirectHandler := NewRedirectHandler(isProd, hostnames, httpsPort)
	if isProd {
		m := &autocert.Manager{
			Cache:      autocert.DirCache(cacheDir),
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(hostnames...),
		}
		h := m.HTTPHandler(redirectHandler)
		return &tls.Config{
			GetCertificate: m.GetCertificate,
			MinVersion:     minVersion,
			CipherSuites:   cipherSuites,
		}, h
	} else {
		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			log.Fatalf("loading self-signed cert: %v", err)
		}
		h := redirectHandler
		return &tls.Config{
			Certificates: []tls.Certificate{cert},
			MinVersion:   minVersion,
			CipherSuites: cipherSuites,
		}, h
	}
}

// NewRedirectHandler returns an http.HandlerFunc that dynamically builds the target redirect URL
func NewRedirectHandler(isProd bool, hostnames []string, httpsPort string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var target string

		if isProd && len(hostnames) > 0 {
			// Always use main hostname (production)
			target = "https://" + hostnames[0] + r.URL.RequestURI()
		} else {
			// Development: use localhost + the HTTPS port
			// Normalize the httpsPort in form ":8888"
			target = "https://localhost" + httpsPort + r.URL.RequestURI()
		}

		http.Redirect(w, r, target, http.StatusMovedPermanently)
	}
}

func StartHTTP(addr string, handler http.Handler, logger *slog.Logger) {
	logger.Info("HTTP (for ACME/redirect) listening on", "addr", addr)
	log.Fatal(http.ListenAndServe(addr, handler))
}

func StartHTTPS(addr string, handler http.Handler, app *app.Application, tlsConfig *tls.Config) {
	server := &http.Server{
		Addr:         addr,
		TLSConfig:    tlsConfig, // Unwrap if needed in your actual code
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
		Handler:      handler,
	}
	go func() {
		app.Logger.Info("HTTPS listening on", "addr", addr)
		if err := server.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
			app.Logger.Error("HTTPS server error", "error", err)
		}
	}()
	waitForShutdown(server, app.Logger)
}

func waitForShutdown(server *http.Server, logger *slog.Logger) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	logger.Info("shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server shutdown error", "error", err)
	} else {
		logger.Info("Server stopped gracefully")
	}
}
