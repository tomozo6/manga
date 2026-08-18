package main

import (
	"net/http"
	"strconv"
)

func (a app) handleMe(w http.ResponseWriter, r *http.Request) {
	user, ok := a.authenticate(w, r)
	if !ok {
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"uid": user.UID, "name": user.Name, "email": user.Email})
}

func (a app) handleManga(w http.ResponseWriter, r *http.Request) {
	if _, ok := a.authenticate(w, r); !ok {
		return
	}
	type response struct {
		ID     string `json:"id"`
		Title  string `json:"title"`
		Author string `json:"author"`
	}
	rows, err := a.db.QueryContext(r.Context(), `SELECT id, title, author_name FROM mangas ORDER BY title, id`)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not read catalog"})
		return
	}
	defer rows.Close()
	var items []response
	for rows.Next() {
		var item response
		if err := rows.Scan(&item.ID, &item.Title, &item.Author); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not read catalog"})
			return
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not read catalog"})
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (a app) handleMangaDetail(w http.ResponseWriter, r *http.Request) {
	if _, ok := a.authenticate(w, r); !ok {
		return
	}
	item, found, err := a.findManga(r.Context(), r.PathValue("mangaID"))
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not read catalog"})
		return
	}
	if !found {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "manga not found"})
		return
	}
	type volumeResponse struct {
		ID     string `json:"id"`
		Number int    `json:"number"`
		Title  string `json:"title"`
	}
	rows, err := a.db.QueryContext(r.Context(), `SELECT id, number, title FROM volumes WHERE manga_id = ? ORDER BY number`, item.ID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not read catalog"})
		return
	}
	defer rows.Close()
	var volumes []volumeResponse
	for rows.Next() {
		var volume volumeResponse
		if err := rows.Scan(&volume.ID, &volume.Number, &volume.Title); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not read catalog"})
			return
		}
		volumes = append(volumes, volume)
	}
	if err := rows.Err(); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not read catalog"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"id": item.ID, "title": item.Title, "author": item.Author, "volumes": volumes})
}

func (a app) handleVolume(w http.ResponseWriter, r *http.Request) {
	if _, ok := a.authenticate(w, r); !ok {
		return
	}
	item, selectedVolume, found, err := a.findVolume(r.Context(), r.PathValue("mangaID"), r.PathValue("volumeID"))
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not read catalog"})
		return
	}
	if !found {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "volume not found"})
		return
	}
	pages, err := a.issuePageBatch(r.Context(), item, selectedVolume, 1)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not issue image URL"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"mangaTitle":   item.Title,
		"volumeNumber": selectedVolume.Number,
		"volumeTitle":  selectedVolume.Title,
		"pageCount":    selectedVolume.PageCount,
		"pages":        pages,
	})
}

func (a app) handleVolumePages(w http.ResponseWriter, r *http.Request) {
	if _, ok := a.authenticate(w, r); !ok {
		return
	}
	start, err := strconv.Atoi(r.URL.Query().Get("start"))
	if err != nil || start < 1 || (start-1)%pageURLBatchSize != 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "start must be a positive page URL batch boundary"})
		return
	}
	item, selectedVolume, found, err := a.findVolume(r.Context(), r.PathValue("mangaID"), r.PathValue("volumeID"))
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not read catalog"})
		return
	}
	if !found {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "volume not found"})
		return
	}
	if start > selectedVolume.PageCount {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "start is beyond the last page"})
		return
	}
	pages, err := a.issuePageBatch(r.Context(), item, selectedVolume, start)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not issue image URL"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"pages": pages})
}
