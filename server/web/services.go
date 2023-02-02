package web

import (
	"code.mrmelon54.com/melon/status/server/structure"
	"code.mrmelon54.com/melon/status/server/utils"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

func (w *Web) allServices(rw http.ResponseWriter, req *http.Request, _ *utils.State, isAdmin bool) {
	// Parse input
	q := req.URL.Query()
	groupId, err := strconv.Atoi(q.Get("group"))
	if err != nil {
		http.Error(rw, "Invalid group parameter", http.StatusBadRequest)
		return
	}

	// Check group is public
	var group structure.Group
	r := w.engine.NewSession()
	if !isAdmin {
		r = r.Where("public = 1")
	}
	foundGroup, err := r.Where("id = ?", groupId).Get(&group)
	if err != nil {
		log.Println(err)
		http.Error(rw, "Failed to read from database", http.StatusInternalServerError)
		return
	}
	if !foundGroup {
		http.Error(rw, "404 Not Found", http.StatusNotFound)
		return
	}

	// Find services
	var services []structure.Service
	r = w.engine.NewSession()
	if !isAdmin {
		r = r.Where("public = 1")
	}
	err = r.Where("group_id = ?", groupId).Asc("order").Find(&services)
	if err != nil {
		log.Println(err)
		http.Error(rw, "Failed to read from database", http.StatusInternalServerError)
		return
	}
	if services == nil {
		services = []structure.Service{}
	}

	// Encode output
	encoder := json.NewEncoder(rw)
	err = encoder.Encode(services)
	if err != nil {
		http.Error(rw, "Failed to encode data", http.StatusInternalServerError)
		return
	}
}

func (w *Web) getService(rw http.ResponseWriter, req *http.Request, _ *utils.State, isAdmin bool) {
	// Parse input
	q := req.URL.Query()
	groupId, err := strconv.Atoi(q.Get("service"))
	if err != nil {
		http.Error(rw, "Invalid group parameter", http.StatusBadRequest)
		return
	}

	// Find group
	var service structure.Service
	r := w.engine.NewSession()
	if !isAdmin {
		r = r.Where("public = 1")
	}
	b, err := r.Where("id = ?", groupId).Asc("order").Get(&service)
	if err != nil {
		log.Println(err)
		http.Error(rw, "Failed to read from database", http.StatusInternalServerError)
		return
	}
	if !b {
		http.Error(rw, "404 Not Found", http.StatusNotFound)
		return
	}

	// Encode output
	encoder := json.NewEncoder(rw)
	err = encoder.Encode(service)
	if err != nil {
		http.Error(rw, "Failed to encode data", http.StatusInternalServerError)
		return
	}
}

func (w *Web) getServiceStatus(rw http.ResponseWriter, req *http.Request, _ *utils.State, isAdmin bool) {
	// Parse input
	q := req.URL.Query()
	serviceId, err := strconv.Atoi(q.Get("service"))
	if err != nil {
		http.Error(rw, "Invalid service parameter", http.StatusBadRequest)
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
	var beans [90]structure.Bean
	t := time.Now().Add(-90 * 24 * time.Hour)
	for _, hit := range hits {
		day := int(hit.CreatedAt.Sub(t).Hours() / 24)
		if day >= 0 && day < 90 {
			beans[day] = structure.Bean{State: structure.BeanStateHit, CreatedAt: hit.CreatedAt, Time: hit.CreatedAt.Unix()}
			if beans[day].CreatedAt.After(current.CreatedAt) {
				current = beans[day]
			}
		}
	}
	for _, failure := range failures {
		day := int(failure.CreatedAt.Sub(t).Hours() / 24)
		if day >= 0 && day < 90 {
			beans[day] = structure.Bean{State: structure.BeanStateFailure, CreatedAt: failure.CreatedAt, Time: failure.CreatedAt.Unix()}
			if current.State == structure.BeanStateUnknown || beans[day].CreatedAt.After(current.CreatedAt) {
				current = beans[day]
			}
		}
	}

	// Encode output
	encoder := json.NewEncoder(rw)
	err = encoder.Encode(struct {
		Current structure.Bean     `json:"current"`
		Beans   [90]structure.Bean `json:"beans"`
	}{
		Current: current,
		Beans:   beans,
	})
	if err != nil {
		http.Error(rw, "Failed to encode data", http.StatusInternalServerError)
		return
	}
}
