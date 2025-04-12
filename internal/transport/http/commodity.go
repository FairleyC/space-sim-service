package http

import (
	"encoding/json"
	"net/http"

	"github.com/FairleyC/space-sim-service/internal/commodity"
)

func (h *Handler) GetCommodities(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) GetCommodity(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) PostCommodity(w http.ResponseWriter, r *http.Request) {
	var commodity commodity.Commodity
	if err := json.NewDecoder(r.Body).Decode(&commodity); err != nil {
		return
	}

	commodity, err := h.Service.CreateCommodity(r.Context(), commodity)
	if err != nil {
		return
	}

	if err := json.NewEncoder(w).Encode(commodity); err != nil {
		return
	}
}

func (h *Handler) PutCommodityPrice(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) DeleteCommodity(w http.ResponseWriter, r *http.Request) {}
