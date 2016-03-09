package page

import (
	"github.com/gorilla/mux"
	"github.com/nuveo/utils/adapter"
)

const (
	path = "/page/"
)

var Db adapter.Driver

func routes(r *mux.Router) {
	r.HandleFunc(path, List).Methods("GET")
	r.HandleFunc(path, Create).Methods("POST")
	// r.HandleFunc(path+"{ID}", Read).Methods("GET")
	// r.HandleFunc(path+"{ID}", Update).Methods("POST")
	// r.HandleFunc(path+"{ID}", Delete).Methods("DELETE")
}

func Register(r *mux.Router, db adapter.Driver) {
	routes(r)
	Db = db.Copy()
}
