package main

import (
	"encoding/json"
	"net/http"
	"path/filepath"
)

func newRouter(a app) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/me", a.handleMe)
	mux.HandleFunc("GET /api/manga", a.handleManga)
	mux.HandleFunc("GET /api/manga/{mangaID}", a.handleMangaDetail)
	mux.HandleFunc("GET /api/manga/{mangaID}/volumes/{volumeID}", a.handleVolume)
	mux.HandleFunc("GET /api/manga/{mangaID}/volumes/{volumeID}/pages", a.handleVolumePages)
	mux.HandleFunc("GET /", servePage("index.html"))
	mux.HandleFunc("GET /library", servePage("library.html"))
	mux.HandleFunc("GET /manga/{mangaID}", servePage("manga.html"))
	mux.HandleFunc("GET /manga/{mangaID}/volumes/{volumeID}", servePage("reader.html"))
	mux.Handle("GET /assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("public/assets"))))
	return mux
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func servePage(name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join("public", name))
	}
}
