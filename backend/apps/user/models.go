package user

import (
	"encoding/json"

	"golang.org/x/crypto/scrypt"

	"labix.org/v2/mgo/bson"
)

type User struct {
	ID       bson.ObjectId `bson:"_id,omitempty" json:"_id"`
	ClientID bson.ObjectId `bson:"client_id" json:"client_id"`
	Name     string        `bson:"name"`
	Username string        `bson:"username"`
	Email    string        `bson:"email"`
	Secret   []byte        `bson:"secret,omitempty" json:"-"`
}

func NewUser() *User {
	return &User{}
}

func (u *User) encryptSecret() {
	salt := make([]byte, 16)
	dk, _ := scrypt.Key(u.Secret, salt, 16384, 8, 1, 32)
	u.Secret = dk
}

func (u *User) Create() error {
	u.ID = bson.NewObjectId()
	u.encryptSecret()

	err := Db.Insert("user", u)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) Read() (res []User, err error) {
	query := bson.M{}
	if u.ID != "" {
		query["_id"] = u.ID
	}

	if u.ClientID != "" {
		query["client_id"] = u.ClientID
	}

	if u.Username != "" {
		query["username"] = u.Username
	}

	if u.Email != "" {
		query["email"] = u.Email
	}

	s, err := Db.Find("user", query)
	if err != nil {
		return nil, err
	}

	x, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(x, &res)
	if err != nil {
		return nil, err
	}

	return
}
