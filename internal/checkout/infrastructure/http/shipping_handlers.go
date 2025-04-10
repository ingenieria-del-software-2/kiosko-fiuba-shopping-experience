package http

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ingenieria-del-software-2/kiosko-fiuba-shopping-experience/internal/checkout/app/services"
	"github.com/ingenieria-del-software-2/kiosko-fiuba-shopping-experience/internal/checkout/app/services/dto"
	"github.com/ingenieria-del-software-2/kiosko-fiuba-shopping-experience/internal/common/errors"
)

// ShippingHandler handles HTTP requests for shipping operations
type ShippingHandler struct {
	shippingService *services.ShippingService
}

// NewShippingHandler creates a new shipping handler
func NewShippingHandler(shippingService *services.ShippingService) *ShippingHandler {
	return &ShippingHandler{
		shippingService: shippingService,
	}
}

// RegisterRoutes registers the shipping routes on the given router
func (h *ShippingHandler) RegisterRoutes(router *mux.Router) {
	// Create a subrouter for shipping routes
	shippingRouter := router.PathPrefix("/shipping").Subrouter()

	// Register routes
	shippingRouter.HandleFunc("/addresses", h.AddShippingAddress).Methods("POST")
	shippingRouter.HandleFunc("/addresses", h.GetUserShippingAddresses).Methods("GET")
	shippingRouter.HandleFunc("/addresses/{addressId}", h.GetShippingAddress).Methods("GET")
	shippingRouter.HandleFunc("/addresses/{addressId}", h.UpdateShippingAddress).Methods("PUT")
	shippingRouter.HandleFunc("/addresses/{addressId}", h.DeleteShippingAddress).Methods("DELETE")
	shippingRouter.HandleFunc("/methods", h.GetShippingMethods).Methods("GET")
}

// AddShippingAddress handles the request to add a shipping address
func (h *ShippingHandler) AddShippingAddress(w http.ResponseWriter, r *http.Request) {
	var req dto.ShippingAddressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	address, err := h.shippingService.AddShippingAddress(r.Context(), &req)
	if err != nil {
		errors.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(address)
}

// GetShippingAddress handles the request to get a shipping address by ID
func (h *ShippingHandler) GetShippingAddress(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	addressID := vars["addressId"]

	address, err := h.shippingService.GetShippingAddress(r.Context(), addressID)
	if err != nil {
		if err.Error() == "shipping address not found" {
			errors.WriteErrorResponse(w, http.StatusNotFound, "Shipping address not found")
		} else {
			errors.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(address)
}

// GetUserShippingAddresses handles the request to get all shipping addresses for a user
func (h *ShippingHandler) GetUserShippingAddresses(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userId")
	if userID == "" {
		errors.WriteErrorResponse(w, http.StatusBadRequest, "User ID is required")
		return
	}

	addresses, err := h.shippingService.GetUserShippingAddresses(r.Context(), userID)
	if err != nil {
		errors.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(addresses)
}

// UpdateShippingAddress handles the request to update a shipping address
func (h *ShippingHandler) UpdateShippingAddress(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	addressID := vars["addressId"]

	var req dto.ShippingAddressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	address, err := h.shippingService.UpdateShippingAddress(r.Context(), addressID, &req)
	if err != nil {
		if err.Error() == "shipping address not found" {
			errors.WriteErrorResponse(w, http.StatusNotFound, err.Error())
		} else if err.Error() == "shipping address does not belong to the user" {
			errors.WriteErrorResponse(w, http.StatusForbidden, err.Error())
		} else {
			errors.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(address)
}

// DeleteShippingAddress handles the request to delete a shipping address
func (h *ShippingHandler) DeleteShippingAddress(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	addressID := vars["addressId"]

	err := h.shippingService.DeleteShippingAddress(r.Context(), addressID)
	if err != nil {
		if err.Error() == "shipping address not found" {
			errors.WriteErrorResponse(w, http.StatusNotFound, err.Error())
		} else {
			errors.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetShippingMethods handles the request to get all shipping methods
func (h *ShippingHandler) GetShippingMethods(w http.ResponseWriter, r *http.Request) {
	methods, err := h.shippingService.GetShippingMethods(r.Context())
	if err != nil {
		errors.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(methods)
}
