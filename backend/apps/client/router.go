package client

import (
	"github.com/gorilla/mux"
	"github.com/nuveo/utils/adapter"
)

const (
	path = "/client/"
)

var Db adapter.Driver

func routes(r *mux.Router) {
	// r.PathPrefix(path).Handler(
	r.HandleFunc(path, List).Methods("GET")
	r.HandleFunc(path, Create).Methods("POST")
	r.HandleFunc(path+"{id}", Read).Methods("GET")
	// r.HandleFunc(path+"{id}", Update).Methods("POST")
	// r.HandleFunc(path+"{id}", Delete).Methods("DELETE")
}

func Register(r *mux.Router, db adapter.Driver) {
	routes(r)
	Db = db.Copy()
}
