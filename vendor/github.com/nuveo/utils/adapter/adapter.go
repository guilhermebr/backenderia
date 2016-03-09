package adapter

import "labix.org/v2/mgo/bson"

// DB Driver interface
type Driver interface {
	Conn() error
	Find(collection string, query bson.M, params ...int) ([]map[string]interface{}, error)
	Insert(collection string, model interface{}) error
	Update(collection string, where bson.M, model interface{}) error
	UpdateAll(collection string, where bson.M, model interface{}) error
	Copy() Driver
	Close()
}
