package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/FairleyC/space-sim-service/internal/data"
	"github.com/FairleyC/space-sim-service/internal/solarSystem"
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
