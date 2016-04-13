package client

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"labix.org/v2/mgo/bson"
)

type Response struct {
	Status string      `json:"status,omitempty"`
	Data   interface{} `json:"data,omitempty"`
	Error  interface{} `json:"error,omitempty"`
}

func Create(w http.ResponseWriter, r *http.Request) {
	client := NewClient()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err := json.NewDecoder(r.Body).Decode(&client); err != nil {
		w.WriteHeader(422)
		json.NewEncoder(w).Encode(Response{
			Status: "error",
			Error: bson.M{
				"code":    422,
				"message": err.Error(),
			},
		})
		return
	}

	if client.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Status: "fail",
			Error: bson.M{
				"code":    http.StatusBadRequest,
				"message": "Field name is required.",
			},
		})
		return
	}

	err := client.Create()
	if err != nil {
		w.WriteHeader(422)
		json.NewEncoder(w).Encode(Response{
			Status: "error",
			Error: bson.M{
				"code":    422,
				"message": err.Error(),
			},
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(Response{
		Status: "success",
		Data: bson.M{
			"client": client,
		},
	})
}

func List(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	client := NewClient()

	clients, err := client.Read()

	if err != nil {
		w.WriteHeader(422)
		json.NewEncoder(w).Encode(Response{
			Status: "error",
			Error: bson.M{
				"code":    422,
				"message": err.Error(),
			},
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{
		Status: "success",
		Data: bson.M{
			"clients": clients,
		},
	})
}

func Read(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	vars := mux.Vars(r)

	client := NewClient()
	client.Slug = vars["slug"]
	clients, err := client.Read()

	if err != nil {
		w.WriteHeader(422)
		json.NewEncoder(w).Encode(Response{
			Status: "error",
			Error: bson.M{
				"code":    422,
				"message": err.Error(),
			},
		})
		return
	}

	if len(clients) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Status: "fail",
			Error: bson.M{
				"code":    http.StatusBadRequest,
				"message": "Client don't exist.",
			},
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{
		Status: "success",
		Data: bson.M{
			"client": clients[0],
		},
	})
}
