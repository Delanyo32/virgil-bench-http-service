package api

import (
	"net/http"
)

// Router wraps http.ServeMux with middleware support.
type Router struct {
	mux        *http.ServeMux
	middleware []func(http.Handler) http.Handler
}

// NewRouter creates a new Router with the given handler.
func NewRouter(h *Handler) *Router {
	r := &Router{
		mux: http.NewServeMux(),
	}

	// Register routes
	r.mux.HandleFunc("/api/v1/orders", func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPost:
			h.CreateOrder(w, req)
		case http.MethodGet:
			h.ListOrders(w, req)
		default:
			WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	})

	r.mux.HandleFunc("/api/v1/orders/detail", func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		h.GetOrder(w, req)
	})

	r.mux.HandleFunc("/api/v1/inventory", func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			h.ListInventory(w, req)
		case http.MethodPut:
			h.UpdateInventory(w, req)
		default:
			WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	})

	return r
}

// Use adds a middleware to the chain.
func (r *Router) Use(mw func(http.Handler) http.Handler) {
	r.middleware = append(r.middleware, mw)
}

// HandleFunc registers a handler function for the given pattern.
func (r *Router) HandleFunc(pattern string, handler http.HandlerFunc) {
	r.mux.HandleFunc(pattern, handler)
}

// ServeHTTP implements http.Handler, applying middleware chain.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var handler http.Handler = r.mux
	for i := len(r.middleware) - 1; i >= 0; i-- {
		handler = r.middleware[i](handler)
	}
	handler.ServeHTTP(w, req)
}

// ListInventory handles GET /api/v1/inventory.
func (h *Handler) ListInventory(w http.ResponseWriter, r *http.Request) {
	products, err := h.inventorySvc.ListProducts()
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "failed to list inventory")
		return
	}
	WriteJSON(w, http.StatusOK, products)
}

// UpdateInventory handles PUT /api/v1/inventory.
func (h *Handler) UpdateInventory(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ProductID int `json:"product_id"`
		Quantity  int `json:"quantity"`
	}

	if err := decodeJSON(r, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.inventorySvc.UpdateStock(req.ProductID, req.Quantity); err != nil {
		WriteError(w, http.StatusInternalServerError, "failed to update stock")
		return
	}

	WriteJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}
