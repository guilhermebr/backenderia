package elastic

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/olivere/elastic"
)

const (
	pageLimit = 20
)

// ElasticSearch .
type ElasticSearch struct {
	uri    string
	index  string
	client *elastic.Client
}

// NewElasticSearchSession .
func New(uri, index string) *ElasticSearch {
	return &ElasticSearch{uri, index, nil}
}

func (es *ElasticSearch) Conn() error {
	client, err := elastic.NewClient(
		elastic.SetURL(es.uri),
		elastic.SetSniff(false),
		elastic.SetHealthcheckInterval(10*time.Second),
		elastic.SetMaxRetries(5))

	if err != nil {
		// Handle error
		return err
	}

	exists, err := client.IndexExists(es.index).Do()
	if err != nil {
		// Handle error
		return err
	}
	if !exists {
		// Create a new index.
		createIndex, err := client.CreateIndex(es.index).Do()
		if err != nil {
			// Handle error
			return err
		}
		if !createIndex.Acknowledged {
			// Not acknowledged
		}
	}

	es.client = client
	return nil
}

func (es *ElasticSearch) Find(query elastic.Query, table string, params ...int) ([]interface{}, int64, error) {
	var objects []interface{}

	skipCount := 0

	if len(params) >= 1 {
		if params[0] > 1 {
			skipCount = (params[0] - 1) * pageLimit
		}
	}

	searchResult, err := es.client.Search().
		Index(es.index).
		Type(table).
		Query(query).
		From(skipCount).Size(pageLimit).
		Pretty(true).
		Do()
	if err != nil {
		return nil, 0, err
	}

	log.Printf("Query took %d milliseconds\n", searchResult.TookInMillis)

	if searchResult.Hits != nil {
		// Iterate through results
		for _, hit := range searchResult.Hits.Hits {
			var model interface{}

			err := json.Unmarshal(*hit.Source, &model)
			if err != nil {
				return nil, 0, err
			}

			objects = append(objects, model)

		}
	}

	return objects, searchResult.Hits.TotalHits, nil
}

func (es *ElasticSearch) DeleteIndex() {
	es.client.DeleteIndex(es.index).Do()
}

func (es *ElasticSearch) Insert(table string, model interface{}) error {
	_, err := es.client.Index().
		Index(es.index).
		Type(table).
		BodyJson(model).
		Do()

	if err != nil {
		return err
	}

	// Flush data
	_, err = es.client.Flush().Index(es.index).Do()

	if err != nil {
		return err
	}

	return nil
}

func (es *ElasticSearch) InsertByID(table, id string, model interface{}) error {
	_, err := es.client.Index().
		Index(es.index).
		Type(table).
		Id(id).
		BodyJson(model).
		Do()

	if err != nil {
		return err
	}

	// Flush data
	_, err = es.client.Flush().Index(es.index).Do()

	if err != nil {
		return err
	}

	return nil
}

func (es *ElasticSearch) Delete(table string, query elastic.Query) error {
	res, err := es.client.DeleteByQuery().Index(es.index).Type(table).Query(query).Do()
	if err != nil {
		return err
	}

	if res == nil {
		return errors.New("response is nil")
	}

	_, found := res.Indices[es.index]
	if !found {
		log.Printf("expected Found = true; got: %v", found)
	}

	_, err = es.client.Flush().Index(es.index).Do()
	if err != nil {
		return err
	}
	return nil
}

func (es *ElasticSearch) DeleteID(table string, id string) error {
	res, err := es.client.Delete().Index(es.index).Type(table).Id(id).Do()
	if err != nil {
		return err
	}

	if res.Found != true {
		return errors.New("document not found")
	}

	_, err = es.client.Flush().Index(es.index).Do()
	if err != nil {
		return err
	}
	return nil
}

func (es *ElasticSearch) Update(table string, id string, data map[string]interface{}) error {
	_, err := es.client.Update().Index(es.index).Type(table).Id(id).Doc(data).Do()

	if err != nil {
		return err
	}

	return nil
}

//Suggester and Completion

func (es *ElasticSearch) Suggester(tsName, field, text string) []elastic.SearchSuggestion {
	ts := elastic.NewTermSuggester(tsName)
	ts = ts.Text(text)
	ts = ts.Field(field)

	return es.suggester(tsName, ts)
}

func (es *ElasticSearch) Completion(tsName, field, text string) []elastic.SearchSuggestion {
	ts := elastic.NewCompletionSuggester(tsName)
	ts = ts.Text(text)
	ts = ts.Field(field)

	return es.suggester(tsName, ts)
}

func (es *ElasticSearch) suggester(tsName string, term elastic.Suggester) []elastic.SearchSuggestion {
	all := elastic.NewMatchAllQuery()

	searchResult, err := es.client.Search().
		Index(es.index).
		Query(all).
		Suggester(term).
		Do()
	if err != nil {
		log.Fatal(err)
	}

	result, _ := searchResult.Suggest[tsName]

	return result
}

//bulk methods

func (es *ElasticSearch) NewBulk() *elastic.BulkService {
	return es.client.Bulk()
}

func (es *ElasticSearch) AddToBulk(bulk *elastic.BulkService, table string, model interface{}, id string) {
	bulk.Add(elastic.NewBulkIndexRequest().Index(es.index).Type(table).Doc(model).Id(id))
}

func (es *ElasticSearch) SendBulk(bulk *elastic.BulkService) {
	bulk.Do()
}
