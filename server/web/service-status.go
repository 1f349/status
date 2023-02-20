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

func (w *Web) getServiceStatus(rw http.ResponseWriter, req *http.Request) {
	// Parse input
	serviceId, err := strconv.Atoi(mux.Vars(req)["service"])
	if err != nil {
		http.Error(rw, "Invalid service parameter", http.StatusBadRequest)
		return
	}
	days, err := strconv.Atoi(req.URL.Query().Get("days"))
	if err != nil {
		http.Error(rw, "Invalid days parameter", http.StatusBadRequest)
		return
	}
	if days != 90 && days != 60 && days != 30 {
		http.Error(rw, "Invalid days parameter, must be 30, 60 or 90", http.StatusBadRequest)
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
	err = w.engine.SQL("SELECT * FROM hit JOIN (SELECT MAX(id) as max FROM hit WHERE service = ? GROUP BY DATE(`created_at`)) last ON hit.id = last.max ORDER BY `created_at` ASC", service.Id).Find(&hits)
	if err != nil {
		log.Println(err)
		http.Error(rw, "Failed to read from database", http.StatusInternalServerError)
		return
	}
	if hits == nil {
		hits = []structure.Hit{}
	}

	// Find failures
	var failures []structure.Failure
	err = w.engine.SQL("SELECT * FROM failure JOIN (SELECT MAX(id) as max FROM failure WHERE service = ? GROUP BY DATE(`created_at`)) last ON failure.id = last.max ORDER BY `created_at` ASC", service.Id).Find(&failures)
	if err != nil {
		log.Println(err)
		http.Error(rw, "Failed to read from database", http.StatusInternalServerError)
		return
	}
	if failures == nil {
		failures = []structure.Failure{}
	}

	// Calculate beans
	var current structure.Bean
	beans := make([]structure.Bean, days)
	t := time.Now().Add(time.Duration(-days) * 24 * time.Hour)
	for _, hit := range hits {
		day := int(hit.CreatedAt.Sub(t).Hours() / 24)
		if day >= 0 && day < days {
			beans[day] = structure.Bean{State: structure.BeanStateHit, CreatedAt: hit.CreatedAt, Time: hit.CreatedAt.Unix()}
			if beans[day].CreatedAt.After(current.CreatedAt) {
				current = beans[day]
			}
		}
	}
	for _, failure := range failures {
		day := int(failure.CreatedAt.Sub(t).Hours() / 24)
		if day >= 0 && day < days {
			beans[day] = structure.Bean{State: structure.BeanStateFailure, CreatedAt: failure.CreatedAt, Time: failure.CreatedAt.Unix()}
			if current.State == structure.BeanStateUnknown || beans[day].CreatedAt.After(current.CreatedAt) {
				current = beans[day]
			}
		}
	}

	// Encode output
	encoder := json.NewEncoder(rw)
	err = encoder.Encode(struct {
		Current structure.Bean   `json:"current"`
		Beans   []structure.Bean `json:"beans"`
	}{
		Current: current,
		Beans:   beans,
	})
	if err != nil {
		http.Error(rw, "Failed to encode data", http.StatusInternalServerError)
		return
	}
}
