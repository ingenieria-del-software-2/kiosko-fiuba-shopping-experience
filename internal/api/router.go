package api

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/ingenieria-del-software-2/kiosko-fiuba-shopping-experience/docs" // Import generated Swagger docs
	cartHttp "github.com/ingenieria-del-software-2/kiosko-fiuba-shopping-experience/internal/cart/infrastructure/http"
	checkoutHttp "github.com/ingenieria-del-software-2/kiosko-fiuba-shopping-experience/internal/checkout/infrastructure/http"
	httpSwagger "github.com/swaggo/http-swagger"
)

// RegisterRoutes registers all API routes
func RegisterRoutes(
	router *mux.Router,
	cartHandler *cartHttp.CartHandler,
	checkoutHandler *checkoutHttp.CheckoutHandler,
	shippingHandler *checkoutHttp.ShippingHandler,
) {
	// Create an API subrouter
	apiRouter := router.PathPrefix("/api").Subrouter()

	// Health check endpoint
	// @Summary Health check
	// @Description Get service health status
	// @Tags health
	// @Produce json
	// @Success 200 {object} map[string]string "Service is healthy"
	// @Router /api/health [get]
	apiRouter.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","service":"shopping-experience"}`))
	}).Methods("GET")

	// Get API path prefix from environment or default to empty string
	apiPathPrefix := os.Getenv("API_PATH_PREFIX")
	if apiPathPrefix == "" {
		apiPathPrefix = ""
	}

	// Swagger documentation - add at root router level for better visibility
	router.PathPrefix("/api/docs/").Handler(httpSwagger.Handler(
		httpSwagger.URL(apiPathPrefix+"/api/docs/doc.json"), // The URL pointing to API definition
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	))

	// Register routes for each handler
	cartHandler.RegisterRoutes(apiRouter)
	checkoutHandler.RegisterRoutes(apiRouter)
	shippingHandler.RegisterRoutes(apiRouter)
}
