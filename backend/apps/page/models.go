package page

import (
	"encoding/json"
	"errors"

	"github.com/avelino/slugify"

	"labix.org/v2/mgo/bson"
)

type Page struct {
	ID      bson.ObjectId `json:"_id" bson:"_id,omitempty"`
	Title   string        `json:"title"`
	Content string        `json:"content"`
	Active  bool          `json:"active"`
	Slug    string        `json:"slug"`
}

func NewPage() *Page {
	return &Page{}
}

func (p *Page) Create() error {

	if p.Title == "" {
		return errors.New("Field title is required.")
	}

	p.ID = bson.NewObjectId()

	if p.Slug == "" {
		p.Slug = slugify.Slugify(p.Title)
	}

	err := Db.Insert("page", p)
	if err != nil {
		return err
	}

	return nil
}

func (p *Page) Read() (res []Page, err error) {

	query := bson.M{}
	if p.ID != "" {
		query["_id"] = p.ID
	}

	if p.Slug != "" {
		query["slug"] = p.Slug
	}

	if p.Active == true {
		query["active"] = true
	}

	s, err := Db.Find("page", query)
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
