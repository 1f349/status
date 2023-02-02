package web

import (
	"code.mrmelon54.com/melon/status/server/structure"
	"code.mrmelon54.com/melon/status/server/utils"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func (w *Web) allHits(rw http.ResponseWriter, req *http.Request, _ *utils.State, isAdmin bool) {
	// Parse input
	q := req.URL.Query()
	serviceId, err := strconv.Atoi(q.Get("service"))
	if err != nil {
		http.Error(rw, "Invalid group parameter", http.StatusBadRequest)
		return
	}

	// Check service is public
	var service structure.Service
	r := w.engine.NewSession()
	if !isAdmin {
		r = r.Where("public = 1")
	}
	foundService, err := r.Where("id = ?", serviceId).Get(&service)
	if err != nil {
		log.Println(err)
		http.Error(rw, "Failed to read from database", http.StatusInternalServerError)
		return
	}
	if !foundService {
		http.Error(rw, "404 Not Found", http.StatusNotFound)
		return
	}

	// Check group is public
	var group structure.Group
	r = w.engine.NewSession()
	if !isAdmin {
		r = r.Where("public = 1")
	}
	foundGroup, err := r.Where("id = ?", service.GroupId).Get(&group)
	if err != nil {
		log.Println(err)
		http.Error(rw, "Failed to read from database", http.StatusInternalServerError)
		return
	}
	if !foundGroup {
		http.Error(rw, "404 Not Found", http.StatusNotFound)
		return
	}

	// Find hits
	var hits []structure.Hit
	err = w.engine.Where("service = ?", serviceId).Asc("created_at").Find(&hits)
	if err != nil {
		log.Println(err)
		http.Error(rw, "Failed to read from database", http.StatusInternalServerError)
		return
	}
	if hits == nil {
		hits = []structure.Hit{}
	}

	// Encode output
	encoder := json.NewEncoder(rw)
	err = encoder.Encode(hits)
	if err != nil {
		http.Error(rw, "Failed to encode data", http.StatusInternalServerError)
		return
	}
}
