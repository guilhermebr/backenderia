package news

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
	news := NewNews()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err := json.NewDecoder(r.Body).Decode(&news); err != nil {
		json.NewEncoder(w).Encode(Response{
			Status: "error",
			Error: bson.M{
				"code":    422,
				"message": err.Error(),
			},
		})
		return
	}

	err := news.Create()
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
			"news": news,
		},
	})
}

func List(w http.ResponseWriter, r *http.Request) {
	n := NewNews()

	news, err := n.Read()

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
			"news": news,
		},
	})
}

func Read(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	n := NewNews()
	n.ID = bson.ObjectIdHex(vars["id"])
	news, err := n.Read()

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
			"news": news[0],
		},
	})
}
