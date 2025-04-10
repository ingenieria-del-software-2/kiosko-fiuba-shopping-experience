package api

import (
	"net/http"

	"github.com/gorilla/mux"
	cartHttp "github.com/kiosko-fiuba/shopping-experience/internal/cart/infrastructure/http"
	checkoutHttp "github.com/kiosko-fiuba/shopping-experience/internal/checkout/infrastructure/http"
)

// RegisterRoutes registers all API routes
func RegisterRoutes(
	router *mux.Router,
	cartHandler *cartHttp.CartHandler,
	checkoutHandler *checkoutHttp.CheckoutHandler,
	shippingHandler *checkoutHttp.ShippingHandler,
) {
	// Health check endpoint
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","service":"shopping-experience"}`))
	}).Methods("GET")

	// Register routes for each handler
	cartHandler.RegisterRoutes(router)
	checkoutHandler.RegisterRoutes(router)
	shippingHandler.RegisterRoutes(router)
}
