package main

import (
	"flag"
	"forum-app/server"
	"log"
	"runtime/debug"
)

func main() {
	defer HandlePanic()

	// Load environment variables from .env if present
	server.LoadEnvironment(".env")

	// Parse command line flags
	addr := flag.String("addr", ":8888", "HTTP network address")
	dbName := flag.String("db", "app.db", "Database file name sqlite3")
	flag.Parse()

	// Initialize the application with all dependencies
	application, httpHandler, http80Handler, tlsConfig := server.InitializeApp(*dbName)

	// Serve HTTP on :8080 for ACME + redirect
	go server.StartHTTP(":8080", http80Handler, application.Logger)

	// Serve HTTPS with appropriate config and graceful shutdown
	server.StartHTTPS(*addr, httpHandler, application, tlsConfig)
}

func HandlePanic() {
	if rec := recover(); rec != nil {
		log.Printf("Critical server error: %v\nStack Trace:\n%s", rec, debug.Stack())
		log.Println("Restarting server...")
		go main() // Relaunch main on panic
	}
}
