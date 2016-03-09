package mongo

import (
	"time"

	"github.com/nuveo/utils/adapter"

	"log"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type Mongo struct {
	uri       string
	database  string
	session   *mgo.Session
	pageLimit int
	timeout   int
}

func New(uri, database string) *Mongo {
	return &Mongo{uri, database, nil, 20, 30}
}

func (m *Mongo) SetLimit(limit int) {
	m.pageLimit = limit
}

func (m *Mongo) SetTimeout(timeout int) {
	m.timeout = timeout
}

func (m *Mongo) Conn() error {
	dialInfo := &mgo.DialInfo{
		Addrs:    []string{m.uri},
		Timeout:  time.Duration(m.timeout) * time.Second,
		Database: m.database,
	}

	sess, err := mgo.DialWithInfo(dialInfo)

	if err != nil {
		return err
	}

	m.session = sess
	return nil
}

func (m *Mongo) Copy() adapter.Driver {
	sessionCopy := m.session.Copy()
	return &Mongo{m.uri, m.database, sessionCopy, m.pageLimit, m.timeout}
}

func (m *Mongo) Insert(collection string, model interface{}) error {
	coll := m.session.DB(m.database).C(collection)
	err := coll.Insert(model)

	if err != nil {
		return err
	}

	return nil
}

func (m *Mongo) Update(collection string, where bson.M, model interface{}) error {
	coll := m.session.DB(m.database).C(collection)
	err := coll.Update(where, model)

	if err != nil {
		return err
	}

	return nil
}

func (m *Mongo) UpdateAll(collection string, where bson.M, model interface{}) error {
	coll := m.session.DB(m.database).C(collection)
	info, err := coll.UpdateAll(where, model)

	if err != nil {
		return err
	}

	log.Println(info)

	return nil
}

func (m *Mongo) Find(collection string, query bson.M, params ...int) ([]map[string]interface{}, error) {
	var objects []map[string]interface{}

	coll := m.session.DB(m.database).C(collection)

	skipCount := 0

	if len(params) >= 1 && params[0] > 1 {
		skipCount = (params[0] - 1) * m.pageLimit
	}

	err := coll.Find(query).Skip(skipCount).Limit(m.pageLimit).All(&objects)

	if err != nil {
		return nil, err
	}

	return objects, nil
}

func (m *Mongo) DeleteID(collection string, id bson.ObjectId) error {
	coll := m.session.DB(m.database).C(collection)
	err := coll.RemoveId(id)

	if err != nil {
		return err
	}

	return nil
}

func (m *Mongo) DropCollection(collection string) error {
	err := m.session.DB(m.database).C(collection).DropCollection()
	if err != nil {
		return err
	}
	return nil

}

func (m *Mongo) Index(collection string, index string) error {
	c := m.session.DB(m.database).C(collection)

	i := mgo.Index{
		Key:        []string{index},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	err := c.EnsureIndex(i)
	if err != nil {
		return err
	}

	return nil
}

func (m *Mongo) Close() {
	m.session.Close()
}
