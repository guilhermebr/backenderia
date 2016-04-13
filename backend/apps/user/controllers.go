package user

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
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	c := map[string]string{}
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
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

	if c["username"] == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Status: "fail",
			Error: bson.M{
				"code":    http.StatusBadRequest,
				"message": "Field username is required.",
			},
		})
		return
	}

	if c["secret"] == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Status: "fail",
			Error: bson.M{
				"code":    http.StatusBadRequest,
				"message": "Field secret is required.",
			},
		})
		return
	}

	if c["email"] == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Status: "fail",
			Error: bson.M{
				"code":    http.StatusBadRequest,
				"message": "Field email is required.",
			},
		})
		return
	}

	if c["client"] == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Status: "fail",
			Error: bson.M{
				"code":    http.StatusBadRequest,
				"message": "Field client is required.",
			},
		})
		return
	}

	//Verify if username already exist for the client
	user := NewUser()
	user.Username = c["username"]
	user.ClientID = bson.ObjectIdHex(c["client"])
	users, err := user.Read()
	if len(users) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Status: "fail",
			Error: bson.M{
				"code":    http.StatusBadRequest,
				"message": "Username already exist",
			},
		})
		return
	}

	//Verify if email already exist
	user = NewUser()
	user.Email = c["email"]
	user.ClientID = bson.ObjectIdHex(c["client"])
	users, err = user.Read()
	if len(users) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Status: "fail",
			Error: bson.M{
				"code":    http.StatusBadRequest,
				"message": "E-mail already exist",
			},
		})
		return
	}

	user.Username = c["username"]
	user.Name = c["name"]
	user.Secret = []byte(c["secret"])

	err = user.Create()
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
			"user": user,
		},
	})
}

func List(w http.ResponseWriter, r *http.Request) {
	user := NewUser()

	users, err := user.Read()

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
			"users": users,
		},
	})
}

func Read(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	user := NewUser()
	user.ID = bson.ObjectIdHex(vars["id"])
	users, err := user.Read()

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
			"user": users[0],
		},
	})
}
