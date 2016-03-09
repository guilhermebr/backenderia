package page

import (
	"encoding/json"
	"net/http"

	"labix.org/v2/mgo/bson"
)

type Response struct {
	Status string      `json:"status,omitempty"`
	Data   interface{} `json:"data,omitempty"`
	Error  interface{} `json:"error,omitempty"`
}

func Create(w http.ResponseWriter, r *http.Request) {
	page := NewPage()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err := json.NewDecoder(r.Body).Decode(&page); err != nil {
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

	if page.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Status: "fail",
			Error: bson.M{
				"code":    http.StatusBadRequest,
				"message": "Field title is required.",
			},
		})
		return
	}

	err := page.Create()
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
			"page": page,
		},
	})
}

func List(w http.ResponseWriter, r *http.Request) {

	page := NewPage()

	pages, err := page.Read()

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
			"pages": pages,
		},
	})
}
