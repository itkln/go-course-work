package main

import (
	"github.com/gorilla/mux"
	"log"
	"movieapp/metadata/internal/controller/metadata"
	httphandler "movieapp/metadata/internal/handler/http"
	"movieapp/metadata/internal/repository/memory"
	"net/http"
)

func main() {
	log.Println("Starting the movie metadata service")
	repo := memory.New()
	ctrl := metadata.New(repo)
	h := httphandler.New(ctrl)
	router := mux.NewRouter()
	router.Handle("/metadata", http.HandlerFunc(h.GetMetadata))
	if err := http.ListenAndServe(":8081", nil); err != nil {
		panic(err)
	}
}
