package http

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ingenieria-del-software-2/kiosko-fiuba-shopping-experience/internal/checkout/app/services"
	"github.com/ingenieria-del-software-2/kiosko-fiuba-shopping-experience/internal/checkout/app/services/dto"
	"github.com/ingenieria-del-software-2/kiosko-fiuba-shopping-experience/internal/common/errors"
)

// CheckoutHandler handles HTTP requests for checkout operations
type CheckoutHandler struct {
	checkoutService *services.CheckoutService
}

// NewCheckoutHandler creates a new checkout handler
func NewCheckoutHandler(checkoutService *services.CheckoutService) *CheckoutHandler {
	return &CheckoutHandler{
		checkoutService: checkoutService,
	}
}

// RegisterRoutes registers the checkout routes on the given router
func (h *CheckoutHandler) RegisterRoutes(router *mux.Router) {
	// Create a subrouter for checkout routes
	checkoutRouter := router.PathPrefix("/checkout").Subrouter()

	// Register routes
	checkoutRouter.HandleFunc("/init", h.InitiateCheckout).Methods("POST")
	checkoutRouter.HandleFunc("/{checkoutId}", h.GetCheckout).Methods("GET")
	checkoutRouter.HandleFunc("/{checkoutId}/shipping", h.UpdateShipping).Methods("PUT")
	checkoutRouter.HandleFunc("/{checkoutId}/payment-method", h.SetPaymentMethod).Methods("PUT")
	checkoutRouter.HandleFunc("/{checkoutId}/complete", h.CompleteCheckout).Methods("POST")
}

// InitiateCheckout handles the request to initialize a checkout
func (h *CheckoutHandler) InitiateCheckout(w http.ResponseWriter, r *http.Request) {
	var req dto.CheckoutInitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	checkout, err := h.checkoutService.InitiateCheckout(r.Context(), &req)
	if err != nil {
		errors.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(checkout)
}

// GetCheckout handles the request to get a checkout by ID
func (h *CheckoutHandler) GetCheckout(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	checkoutID := vars["checkoutId"]

	checkout, err := h.checkoutService.GetCheckout(r.Context(), checkoutID)
	if err != nil {
		if err.Error() == "checkout not found" {
			errors.WriteErrorResponse(w, http.StatusNotFound, "Checkout not found")
		} else {
			errors.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(checkout)
}

// UpdateShipping handles the request to update shipping details
func (h *CheckoutHandler) UpdateShipping(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	checkoutID := vars["checkoutId"]

	var req dto.ShippingDetailsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	checkout, err := h.checkoutService.UpdateShipping(r.Context(), checkoutID, &req)
	if err != nil {
		if err.Error() == "checkout not found" || err.Error() == "shipping address not found" || err.Error() == "shipping method not found" {
			errors.WriteErrorResponse(w, http.StatusNotFound, err.Error())
		} else if err.Error() == "shipping address does not belong to the user" {
			errors.WriteErrorResponse(w, http.StatusForbidden, err.Error())
		} else if err.Error() == "cannot update a cancelled checkout" {
			errors.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		} else {
			errors.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(checkout)
}

// SetPaymentMethod handles the request to set the payment method
func (h *CheckoutHandler) SetPaymentMethod(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	checkoutID := vars["checkoutId"]

	var req dto.PaymentMethodRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	checkout, err := h.checkoutService.SetPaymentMethod(r.Context(), checkoutID, &req)
	if err != nil {
		if err.Error() == "checkout not found" {
			errors.WriteErrorResponse(w, http.StatusNotFound, err.Error())
		} else if err.Error() == "cannot update a cancelled checkout" || err.Error() == "shipping option must be selected before payment" {
			errors.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		} else {
			errors.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(checkout)
}

// CompleteCheckout handles the request to complete a checkout
func (h *CheckoutHandler) CompleteCheckout(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	checkoutID := vars["checkoutId"]

	checkout, err := h.checkoutService.CompleteCheckout(r.Context(), checkoutID)
	if err != nil {
		if err.Error() == "checkout not found" {
			errors.WriteErrorResponse(w, http.StatusNotFound, err.Error())
		} else if err.Error() == "cannot complete a cancelled checkout" || err.Error() == "payment method must be selected before completing checkout" {
			errors.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		} else {
			errors.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(checkout)
}
