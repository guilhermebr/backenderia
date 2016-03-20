package news

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/avelino/slugify"
	"labix.org/v2/mgo/bson"
)

type News struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	Title       string        `json:"title"`
	Hat         string        `json:"hat"`
	Description string        `json:"description"`
	Featured    string        `json:"featured"`
	Content     string        `json:"content"`
	MainImage   string        `json:"main_image"`
	Tags        []string      `json:"tags"`
	PubDate     time.Time     `json:"pub_date"`
	Published   bool          `json:"published"`
	Slug        string        `json:"slug"`
}

func NewNews() *News {
	return &News{}
}

func (n *News) Create() error {
	if n.Title == "" {
		return errors.New("Field title is required.")
	}

	n.ID = bson.NewObjectId()

	if n.Slug == "" {
		n.Slug = slugify.Slugify(n.Title)
	}

	n.PubDate = time.Now()

	err := Db.Insert("news", n)
	if err != nil {
		return err
	}

	return nil
}

func (n *News) Read() (res []News, err error) {
	query := bson.M{}
	if n.ID != "" {
		query["_id"] = n.ID
	}

	if n.Slug != "" {
		query["slug"] = n.Slug
	}

	if n.Published == true {
		query["published"] = true
	}

	s, err := Db.Find("news", query)
	if err != nil {
		return nil, err
	}

	x, err := json.Marshal(s)
	err = json.Unmarshal(x, &res)

	return
}
