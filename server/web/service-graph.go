package web

import (
	"code.mrmelon54.com/melon/status/server/structure"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"time"
)

type GraphItem struct {
	X int64 `json:"x"`
	Y int64 `json:"y"`
}

func (w *Web) getServiceGraph(rw http.ResponseWriter, req *http.Request) {
	// Parse input
	serviceId, err := strconv.Atoi(mux.Vars(req)["service"])
	if err != nil {
		http.Error(rw, "Invalid service parameter", http.StatusBadRequest)
		return
	}

	// Check service is public
	var service structure.Service
	foundService, err := w.engine.Where("public = 1 and id = ?", serviceId).Get(&service)
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
	foundGroup, err := w.engine.Where("public = 1 and id = ?", service.GroupId).Get(&group)
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
	err = w.engine.SQL("SELECT * FROM hit WHERE service = ? and created_at > FROM_UNIXTIME(?) ORDER BY `created_at` ASC", service.Id, time.Now().Add(-24*time.Hour).Unix()).Find(&hits)
	if err != nil {
		log.Println(err)
		http.Error(rw, "Failed to read from database", http.StatusInternalServerError)
		return
	}
	if hits == nil {
		hits = []structure.Hit{}
	}

	a := make([]GraphItem, len(hits))
	for i, j := range hits {
		a[i] = GraphItem{
			X: j.CreatedAt.Unix(),
			Y: j.Latency.Nanoseconds(),
		}
	}

	// Encode output
	encoder := json.NewEncoder(rw)
	err = encoder.Encode(a)
	if err != nil {
		http.Error(rw, "Failed to encode data", http.StatusInternalServerError)
		return
	}
}
