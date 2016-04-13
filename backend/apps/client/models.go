package client

import (
	"encoding/json"

	"github.com/avelino/slugify"

	"labix.org/v2/mgo/bson"
)

type Client struct {
	ID    bson.ObjectId `bson:"_id,omitempty" json:"_id" `
	Name  string        `bson:"name" json:"name"`
	Slug  string        `json:"slug" json:"slug"`
	Email string        `bson:"email json:"email""`
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) Create() error {
	c.ID = bson.NewObjectId()

	if c.Slug != "" {
		c.Slug = slugify.Slugify(c.Slug)
	} else {
		c.Slug = slugify.Slugify(c.Name)
	}

	err := Db.Insert("client", c)

	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Read() (res []Client, err error) {
	query := bson.M{}
	if c.ID != "" {
		query["_id"] = c.ID
	}

	if c.Slug != "" {
		query["slug"] = c.Slug
	}

	if c.Name != "" {
		query["name"] = c.Name
	}

	if c.Email != "" {
		query["email"] = c.Email
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
