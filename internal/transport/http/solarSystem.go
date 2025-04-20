package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/FairleyC/space-sim-service/internal/data"
	"github.com/FairleyC/space-sim-service/internal/services/solarSystem"
	"github.com/gorilla/mux"
)

type SolarSystemResponse struct {
	SolarSystems []solarSystem.SolarSystem `json:"solarSystems"`
	Pagination   data.Pagination           `json:"pagination"`
}

func (h *Handler) GetSolarSystems(w http.ResponseWriter, r *http.Request) {
	log.Println("REQUEST: GetSolarSystems")

	pagination := data.GetPagination(r)

	solarSystems, err := h.SolarSystemService.FindAllSolarSystems(r.Context(), pagination)
	if err != nil {
		log.Println("Error getting solar systems", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(SolarSystemResponse{
		SolarSystems: solarSystems,
		Pagination:   pagination,
	}); err != nil {
		log.Println("Error encoding solar systems", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetSolarSystem(w http.ResponseWriter, r *http.Request) {
	log.Println("REQUEST: GetSolarSystem")
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		log.Println("ID was missing from request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	foundSolarSystem, err := h.SolarSystemService.FindSolarSystem(r.Context(), id)
	if err != nil {
		if errors.Is(err, solarSystem.ErrSolarSystemNotFound) {
			log.Println("Solar system not found", err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		log.Println("Error getting solar system", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(foundSolarSystem); err != nil {
		log.Println("Error encoding solar system", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type SolarSystemJson struct {
	ID   string
	Name string
}

func (h *Handler) PostSolarSystem(w http.ResponseWriter, r *http.Request) {
	log.Println("REQUEST: PostSolarSystem")
	var solarSystemJson SolarSystemJson
	if err := json.NewDecoder(r.Body).Decode(&solarSystemJson); err != nil {
		log.Println("Error decoding solar system", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	solarSystem := solarSystem.SolarSystem{
		ID:   solarSystemJson.ID,
		Name: solarSystemJson.Name,
	}

	solarSystem, err := h.SolarSystemService.CreateSolarSystem(r.Context(), solarSystem)

	if err != nil {
		log.Println("Error creating solar system", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(solarSystem); err != nil {
		log.Println("Error encoding solar system", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) DeleteSolarSystem(w http.ResponseWriter, r *http.Request) {
	log.Println("REQUEST: DeleteSolarSystem")
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := h.SolarSystemService.RemoveSolarSystem(r.Context(), id)
	if err != nil {
		log.Println("Error deleting solar system", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type CommodityMarketJson struct {
	BasePrice      float64
	DemandQuantity int
	CommodityID    string
}

func (h *Handler) PostCommodityMarket(w http.ResponseWriter, r *http.Request) {
	log.Println("REQUEST: PostCommodityMarket")
	vars := mux.Vars(r)

	solarSystemId := vars["solarSystemId"]
	if solarSystemId == "" {
		log.Println("Solar system ID was missing from request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var commodityMarketJson CommodityMarketJson
	if err := json.NewDecoder(r.Body).Decode(&commodityMarketJson); err != nil {
		log.Println("Error decoding commodity market", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	commodityMarket, err := h.SolarSystemService.CreateCommodityMarket(r.Context(), solarSystemId, commodityMarketJson.BasePrice, commodityMarketJson.DemandQuantity, commodityMarketJson.CommodityID)
	if err != nil {
		log.Println("Error creating commodity market", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(commodityMarket); err != nil {
		log.Println("Error encoding commodity market", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type CommodityMarketUpdateJson struct {
	BasePrice      float64
	DemandQuantity int
}

func (h *Handler) PutCommodityMarket(w http.ResponseWriter, r *http.Request) {
	log.Println("REQUEST: PutCommodityMarket")
	vars := mux.Vars(r)
	solarSystemId := vars["solarSystemId"]
	commodityMarketId := vars["commodityMarketId"]

	if solarSystemId == "" || commodityMarketId == "" {
		log.Println("ID was missing from request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	commodityMarketUpdateJson := CommodityMarketUpdateJson{}
	if err := json.NewDecoder(r.Body).Decode(&commodityMarketUpdateJson); err != nil {
		log.Println("Error decoding commodity market update", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	commodityMarketUpdate := solarSystem.CommodityMarketUpdate{
		BasePrice:      commodityMarketUpdateJson.BasePrice,
		DemandQuantity: commodityMarketUpdateJson.DemandQuantity,
	}

	commodityMarket, err := h.SolarSystemService.UpdateCommodityMarket(r.Context(), commodityMarketId, commodityMarketUpdate)
	if err != nil {
		log.Println("Error updating commodity market", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(commodityMarket); err != nil {
		log.Println("Error encoding commodity market", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) DeleteCommodityMarket(w http.ResponseWriter, r *http.Request) {
	log.Println("REQUEST: DeleteCommodityMarket")
	vars := mux.Vars(r)
	solarSystemId := vars["solarSystemId"]
	commodityMarketId := vars["commodityMarketId"]

	if solarSystemId == "" || commodityMarketId == "" {
		log.Println("ID was missing from request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := h.SolarSystemService.RemoveCommodityMarket(r.Context(), commodityMarketId)
	if err != nil {
		log.Println("Error deleting commodity market", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
