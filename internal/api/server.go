package api

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	cartService "github.com/ingenieria-del-software-2/kiosko-fiuba-shopping-experience/internal/cart/app/services"
	cartHttp "github.com/ingenieria-del-software-2/kiosko-fiuba-shopping-experience/internal/cart/infrastructure/http"
	cartRepo "github.com/ingenieria-del-software-2/kiosko-fiuba-shopping-experience/internal/cart/infrastructure/postgresql"
	checkoutService "github.com/ingenieria-del-software-2/kiosko-fiuba-shopping-experience/internal/checkout/app/services"
	checkoutHttp "github.com/ingenieria-del-software-2/kiosko-fiuba-shopping-experience/internal/checkout/infrastructure/http"
	checkoutRepo "github.com/ingenieria-del-software-2/kiosko-fiuba-shopping-experience/internal/checkout/infrastructure/postgresql"
)

// Server represents the API server
type Server struct {
	server *http.Server
}

// NewServer creates a new API server with all dependencies wired up
func NewServer(db *sql.DB) *Server {
	router := mux.NewRouter()

	// Initialize repositories
	cartRepository := cartRepo.NewPostgreSQLCartRepository(db)
	checkoutRepository := checkoutRepo.NewPostgreSQLCheckoutRepository(db)
	shippingRepository := checkoutRepo.NewPostgreSQLShippingRepository(db)

	// Initialize services
	cartSvc := cartService.NewCartService(cartRepository)
	checkoutSvc := checkoutService.NewCheckoutService(checkoutRepository, shippingRepository)
	shippingSvc := checkoutService.NewShippingService(shippingRepository)

	// Initialize handlers
	cartHandler := cartHttp.NewCartHandler(cartSvc)
	checkoutHandler := checkoutHttp.NewCheckoutHandler(checkoutSvc)
	shippingHandler := checkoutHttp.NewShippingHandler(shippingSvc)

	// Register routes
	RegisterRoutes(router, cartHandler, checkoutHandler, shippingHandler)

	// Create HTTP server
	httpServer := &http.Server{
		Addr:         ":8001",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	return &Server{
		server: httpServer,
	}
}

// Start starts the server
func (s *Server) Start() error {
	log.Printf("Server starting on %s", s.server.Addr)
	return s.server.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
