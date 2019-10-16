package controller

import (
	"encoding/json"
	"github.com/Albert221/UnbottledApi/repository"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type PointController struct {
	points repository.PointRepository
}

func NewPointController(points repository.PointRepository) *PointController {
	return &PointController{
		points: points,
	}
}

func (p *PointController) GetPointsHandler(w http.ResponseWriter, r *http.Request) {
	lat, lng, radius, errors := p.parseLatLngRadiusVars(r)
	if len(errors) > 1 {
		w.WriteHeader(http.StatusBadRequest)
		var message string
		for i, err := range errors {
			message += err.Error()
			if i-1 < len(errors) {
				message += "; "
			}
		}
		_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
		return
	}

	points, err := p.points.InArea(lat, lng, radius)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]interface{}{"points": points})
}

func (PointController) parseLatLngRadiusVars(r *http.Request) (float32, float32, float32, []error) {
	vars := mux.Vars(r)

	var errors []error
	lat, err := strconv.ParseFloat(vars["lat"], 32)
	if err != nil {
		errors = append(errors, err)
	}
	lng, err := strconv.ParseFloat(vars["lng"], 32)
	if err != nil {
		errors = append(errors, err)
	}
	radius, err := strconv.ParseFloat(vars["radius"], 32)
	if err != nil {
		errors = append(errors, err)
	}

	return float32(lat), float32(lng), float32(radius), errors
}
