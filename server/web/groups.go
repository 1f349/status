package web

import (
	"code.mrmelon54.com/melon/status/server/structure"
	"code.mrmelon54.com/melon/status/server/utils"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func (w *Web) allGroups(rw http.ResponseWriter, _ *http.Request, _ *utils.State, isAdmin bool) {
	// Find groups
	var groups []structure.Group
	r := w.engine.NewSession()
	if !isAdmin {
		r = r.Where("public = 1")
	}
	err := r.Asc("order").Find(&groups)
	if err != nil {
		log.Println(err)
		http.Error(rw, "Failed to read from database", http.StatusInternalServerError)
		return
	}
	if groups == nil {
		groups = []structure.Group{}
	}

	// Encode output
	encoder := json.NewEncoder(rw)
	err = encoder.Encode(groups)
	if err != nil {
		http.Error(rw, "Failed to encode data", http.StatusInternalServerError)
		return
	}
}

func (w *Web) getGroup(rw http.ResponseWriter, req *http.Request, _ *utils.State, isAdmin bool) {
	// Parse input
	q := req.URL.Query()
	groupId, err := strconv.Atoi(q.Get("group"))
	if err != nil {
		http.Error(rw, "Invalid group parameter", http.StatusBadRequest)
		return
	}

	// Find group
	var group structure.Group
	r := w.engine.NewSession()
	if !isAdmin {
		r = r.Where("public = 1")
	}
	b, err := r.Where("id = ?", groupId).Asc("order").Get(&group)
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
	err = encoder.Encode(group)
	if err != nil {
		http.Error(rw, "Failed to encode data", http.StatusInternalServerError)
		return
	}
}
