package http

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ingenieria-del-software-2/kiosko-fiuba-shopping-experience/internal/cart/app/services"
	"github.com/ingenieria-del-software-2/kiosko-fiuba-shopping-experience/internal/cart/app/services/dto"
	"github.com/ingenieria-del-software-2/kiosko-fiuba-shopping-experience/internal/common/errors"
)

// CartHandler handles HTTP requests for cart operations
type CartHandler struct {
	cartService *services.CartService
}

// NewCartHandler creates a new cart handler
func NewCartHandler(cartService *services.CartService) *CartHandler {
	return &CartHandler{
		cartService: cartService,
	}
}

// RegisterRoutes registers the cart routes on the given router
func (h *CartHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/carts", h.CreateCart).Methods("POST")
	router.HandleFunc("/carts/{cartId}", h.GetCart).Methods("GET")
	router.HandleFunc("/carts/{cartId}", h.DeleteCart).Methods("DELETE")
	router.HandleFunc("/carts/{cartId}/items", h.AddCartItem).Methods("POST")
	router.HandleFunc("/carts/{cartId}/items/{itemId}", h.UpdateCartItem).Methods("PUT")
	router.HandleFunc("/carts/{cartId}/items/{itemId}", h.RemoveCartItem).Methods("DELETE")
}

// CreateCart handles the request to create a new cart
func (h *CartHandler) CreateCart(w http.ResponseWriter, r *http.Request) {
	var req dto.CartCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	cart, err := h.cartService.CreateCart(r.Context(), &req)
	if err != nil {
		errors.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(cart)
}

// GetCart handles the request to get a cart by ID
func (h *CartHandler) GetCart(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cartID := vars["cartId"]

	cart, err := h.cartService.GetCart(r.Context(), cartID)
	if err != nil {
		if err.Error() == "cart not found" {
			errors.WriteErrorResponse(w, http.StatusNotFound, "Cart not found")
		} else {
			errors.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cart)
}

// DeleteCart handles the request to delete a cart
func (h *CartHandler) DeleteCart(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cartID := vars["cartId"]

	err := h.cartService.DeleteCart(r.Context(), cartID)
	if err != nil {
		if err.Error() == "cart not found" {
			errors.WriteErrorResponse(w, http.StatusNotFound, "Cart not found")
		} else {
			errors.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// AddCartItem handles the request to add an item to a cart
func (h *CartHandler) AddCartItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cartID := vars["cartId"]

	var req dto.CartItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	cart, err := h.cartService.AddCartItem(r.Context(), cartID, &req)
	if err != nil {
		if err.Error() == "cart not found" {
			errors.WriteErrorResponse(w, http.StatusNotFound, "Cart not found")
		} else {
			errors.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cart)
}

// UpdateCartItem handles the request to update an item in a cart
func (h *CartHandler) UpdateCartItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cartID := vars["cartId"]
	itemID := vars["itemId"]

	var req dto.CartItemUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	cart, err := h.cartService.UpdateCartItem(r.Context(), cartID, itemID, &req)
	if err != nil {
		if err.Error() == "cart not found" || err.Error() == "item not found in cart" {
			errors.WriteErrorResponse(w, http.StatusNotFound, err.Error())
		} else {
			errors.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cart)
}

// RemoveCartItem handles the request to remove an item from a cart
func (h *CartHandler) RemoveCartItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cartID := vars["cartId"]
	itemID := vars["itemId"]

	err := h.cartService.RemoveCartItem(r.Context(), cartID, itemID)
	if err != nil {
		if err.Error() == "cart not found" || err.Error() == "item not found in cart" {
			errors.WriteErrorResponse(w, http.StatusNotFound, err.Error())
		} else {
			errors.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
