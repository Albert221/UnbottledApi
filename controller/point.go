package controller

import (
	"encoding/json"
	"github.com/Albert221/UnbottledApi/entity"
	"github.com/Albert221/UnbottledApi/repository"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

type PointController struct {
	points repository.PointRepository
	photos repository.PhotoRepository
}

func NewPointController(points repository.PointRepository, photos repository.PhotoRepository) *PointController {
	return &PointController{
		points: points,
		photos: photos,
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

func (p *PointController) UploadPhoto(w http.ResponseWriter, r *http.Request) {
	user := getUser(r)
	if user == nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	mime := r.Header.Get("Content-Type")
	if mime != "image/jpeg" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Only image/jpeg Content-Type is permitted"})
		return
	}

	const maxSize = (1 << 20) * 5 // 5MiB
	if r.ContentLength == -1 || r.ContentLength > maxSize {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Photos with a maximum size of 5MiB are permitted"})
		return
	}

	id := uuid.New()
	fileName := id.String() + ".jpg"

	photo := &entity.Photo{
		Base: entity.Base{
			ID: id,
		},
		AuthorID: user.ID,
		FileName: fileName,
	}

	file, err := os.Create("uploads/" + fileName)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()

	defer r.Body.Close()
	_, err = io.Copy(file, r.Body)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := p.photos.Save(photo); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"photo": photo})
}

func (p *PointController) AddHandler(w http.ResponseWriter, r *http.Request) {
	user := getUser(r)
	if user == nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var body struct {
		Latitude  float32 `json:"latitude" valid:"required,latitude"`
		Longitude float32 `json:"longitude" valid:"required,longitude"`
		PhotoId   string  `json:"photo_id" valid:"required,uuid"`
	}

	if err := decodeAndValidateBody(&body, r); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	photoId, err := uuid.Parse(body.PhotoId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Invalid photo_id"})
		return
	}

	photo := p.photos.ById(photoId)
	if photo == nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Photo with given ID does not exist"})
		return
	}

	point := &entity.Point{
		Latitude:  body.Latitude,
		Longitude: body.Longitude,
		PhotoID:   photoId,
		Photo:     *photo,
		AuthorID:  user.ID,
	}

	if err := p.points.Save(point); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"point": point})
}
