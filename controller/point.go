package controller

import (
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
	points  repository.PointRepository
	photos  repository.PhotoRepository
	ratings repository.RatingRepository
}

func NewPointController(points repository.PointRepository, photos repository.PhotoRepository, ratings repository.RatingRepository) *PointController {
	return &PointController{
		points:  points,
		photos:  photos,
		ratings: ratings,
	}
}

func (p *PointController) GetPointsHandler(w http.ResponseWriter, r *http.Request) {
	lat, lng, radius, errors := p.parseLatLngRadiusVars(r)
	if len(errors) > 1 {
		var message string
		for i, err := range errors {
			message += err.Error()
			if i-1 < len(errors) {
				message += "; "
			}
		}
		writeJSON(w, map[string]string{"error": message}, http.StatusBadRequest)
		return
	}

	points, err := p.points.InArea(lat, lng, radius)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]interface{}{"points": points}, http.StatusOK)
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

func (p *PointController) GetMyPoints(w http.ResponseWriter, r *http.Request) {
	user := getUser(r)
	if user == nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	points := p.points.ByAuthorID(user.ID)

	writeJSON(w, map[string]interface{}{"points": points}, http.StatusOK)
}

func (p *PointController) UploadPhoto(w http.ResponseWriter, r *http.Request) {
	user := getUser(r)
	if user == nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	mime := r.Header.Get("Content-Type")
	if mime != "image/jpeg" {
		writeJSON(w, map[string]string{
			"error": "Only image/jpeg Content-Type is permitted",
		}, http.StatusBadRequest)
		return
	}

	const maxSize = (1 << 20) * 5 // 5MiB
	if r.ContentLength == -1 || r.ContentLength > maxSize {
		writeJSON(w, map[string]string{
			"error": "Photos with a maximum size of 5MiB are permitted",
		}, http.StatusRequestEntityTooLarge)
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

	writeJSON(w, map[string]interface{}{"photo": photo}, http.StatusCreated)
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
		PhotoId   string  `json:"photo_id" valid:"uuid"`
	}

	if err := decodeAndValidateBody(&body, r); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	point := &entity.Point{
		Latitude:  body.Latitude,
		Longitude: body.Longitude,
		AuthorID:  user.ID,
	}

	if body.PhotoId != "" {
		photoId, err := uuid.Parse(body.PhotoId)
		if err != nil {
			writeJSON(w, map[string]string{"error": "Invalid photo_id"}, http.StatusBadRequest)
			return
		}

		photo := p.photos.ByID(photoId)
		if photo == nil {
			writeJSON(w, map[string]string{"error": "Photo with given ID does not exist"}, http.StatusBadRequest)
			return
		}

		point.PhotoID = photoId
		point.Photo = *photo
	}


	if err := p.points.Save(point); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]interface{}{"point": point}, http.StatusCreated)
}

func (p *PointController) RateHandler(w http.ResponseWriter, r *http.Request) {
	// todo(Albert221): add getting point id from link etc., complete this handler

	var body struct {
		Taste int32 `json:"taste" valid:"required,range(1,5)"`
	}

	if err := decodeAndValidateBody(&body, r); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}
}
