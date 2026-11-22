package main

import (
	"context"
	"errors"
	"firstgoproject/internal/database"
	"firstgoproject/internal/handler"
	"firstgoproject/internal/middleware"
	"firstgoproject/internal/repository"
	"firstgoproject/internal/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	db, err := database.ConnectPostgres()
	if err != nil {
		log.Fatal(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()

	log.Println("Connected to PostgreSQL successfully")

	// productRepo := repository.NewInMemoryProductRepository()
	productRepo := repository.NewGormProductRepository(db)
	productService := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productService)

	router := mux.NewRouter()

	// router.Use(middleware.RecoveryMiddleware)
	// router.Use(middleware.RequestIDMiddleware)
	// router.Use(middleware.LoggingMiddleware)
	handler := middleware.Chain(
		router,
		middleware.RequestIDMiddleware,
		middleware.RecoveryMiddleware,
		middleware.LoggingMiddleware,
	)

	// added for test purposes
	// router.HandleFunc("/test/panic", func(w http.ResponseWriter, r *http.Request) {
	// 	panic("intentional test panic")
	// }).Methods(http.MethodGet)
	router.HandleFunc("/test/live", TestPanic)
	router.HandleFunc("/health/live", GetStatus).
		Methods(http.MethodGet)

	router.HandleFunc("/api/v1/products", productHandler.CreateProduct).
		Methods(http.MethodPost)

	router.HandleFunc("/api/v1/products", productHandler.GetProducts).
		Methods(http.MethodGet)

	router.HandleFunc("/api/v1/products/{id}", productHandler.GetProductByID).
		Methods(http.MethodGet)

	router.HandleFunc("/api/v1/products/{id}", productHandler.UpdateProduct).
		Methods(http.MethodPut)

	router.HandleFunc("/api/v1/products/{id}", productHandler.DeleteProduct).
		Methods(http.MethodDelete)

	log.Println("API server running at :8080")

	// if err := http.ListenAndServe(":8080", router); err != nil {
	// 	log.Fatal(err)
	// }
	// graceful shutdown
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()
	// context.Background() — the root context.
	// os.Interrupt — handles Ctrl+C.
	// syscall.SIGTERM — handles a termination request on Unix-like systems.
	// stop() — releases the signal-handling resources when you're done.

	server := http.Server{
		Addr:              ":8080",
		Handler:           handler, // it starts the wrapped handler that is router
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// the below code will block the main and don't happen the graceful shutdown
	// err := server.ListenAndServe()
	// if err != nil && err != http.ErrServerClosed {
	// 	log.Fatal(err)
	// }

	// here through go routines the main() and listenserve will run concurrently adn we hav time for signal context
	serverErrors := make(chan error, 1)

	go func() {
		serverErrors <- server.ListenAndServe()
	}()

	// make(chan error, 1) creates a buffered channel that can hold one error.
	// go func() { ... }() starts the server in a separate goroutine.
	// server.ListenAndServe() starts listening for HTTP requests.
	// serverErrors <- ... sends the result into the channel when ListenAndServe() returns.

	// Wait for either a server error or an OS shutdown signal.
	select {
	case err := <-serverErrors:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server failed: %v", err)
		}
		log.Println("HTTP server stopped")

	case <-ctx.Done():
		log.Println("Shutdown signal received")
	}

	// Give in-flight requests up to 10 seconds to finish.
	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Graceful shutdown failed: %v", err)

		// Force-close connections if graceful shutdown times out or fails.
		if closeErr := server.Close(); closeErr != nil {
			log.Printf("Error closing HTTP server: %v", closeErr)
		}
	}

	log.Println("Server shutdown complete")
}

func GetStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

func TestPanic(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Processing started"))

	panic("something went wrong")
}

// 1. Why use a goroutine?
// ListenAndServe() blocks. Running it with go allows main() to wait for shutdown signals at the same time.

// 2. Why use select?
// It waits for either the server to return an error or the OS signal context to be canceled.

// 3. Why use Shutdown()?
// server.Shutdown(shutdownCtx) stops accepting new connections and waits for active HTTP requests to finish, up to the deadline.

// 4. Why use Close()?
// If graceful shutdown fails or times out, Close() force-closes active connections. It is the fallback, not the normal shutdown path.
