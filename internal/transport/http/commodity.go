package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/FairleyC/space-sim-service/internal/commodity"
	"github.com/FairleyC/space-sim-service/internal/data"
	"github.com/gorilla/mux"
)

type CommodityResponse struct {
	Commodities []commodity.Commodity `json:"commodities"`
	Pagination  data.Pagination       `json:"pagination"`
}

func (h *Handler) GetCommodities(w http.ResponseWriter, r *http.Request) {
	log.Println("REQUEST: GetCommodities")

	pagination := data.GetPagination(r)

	commodities, err := h.Service.GetAllCommodity(r.Context(), pagination)
	if err != nil {
		log.Println("Error getting commodities", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(CommodityResponse{
		Commodities: commodities,
		Pagination:  pagination,
	}); err != nil {
		log.Println("Error encoding commodities", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetCommodity(w http.ResponseWriter, r *http.Request) {
	log.Println("REQUEST: GetCommodity")
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		log.Println("ID was missing from request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	foundCommodity, err := h.Service.GetCommodity(r.Context(), id)
	if err != nil {
		if errors.Is(err, commodity.ErrCommodityNotFound) {
			log.Println("Commodity not found", err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		log.Println("Error getting commodity", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(foundCommodity); err != nil {
		log.Println("Error encoding commodity", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type CommodityJson struct {
	ID         string
	Name       string
	UnitMass   float64
	UnitVolume float64
}

func (h *Handler) PostCommodity(w http.ResponseWriter, r *http.Request) {
	log.Println("REQUEST: PostCommodity")
	var commodityJson CommodityJson
	if err := json.NewDecoder(r.Body).Decode(&commodityJson); err != nil {
		log.Println("Error decoding commodity", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	commodity := commodity.Commodity{
		ID:         commodityJson.ID,
		Name:       commodityJson.Name,
		UnitMass:   commodityJson.UnitMass,
		UnitVolume: commodityJson.UnitVolume,
	}

	commodity, err := h.Service.CreateCommodity(r.Context(), commodity)

	if err != nil {
		log.Println("Error creating commodity", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(commodity); err != nil {
		log.Println("Error encoding commodity", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) DeleteCommodity(w http.ResponseWriter, r *http.Request) {
	log.Println("REQUEST: DeleteCommodity")
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := h.Service.RemoveCommodity(r.Context(), id)
	if err != nil {
		log.Println("Error deleting commodity", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
