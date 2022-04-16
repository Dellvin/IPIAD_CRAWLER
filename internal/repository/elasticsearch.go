package repository

import (
	"IPIAD_DZ/config"
	"IPIAD_DZ/internal/model"
	"context"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic/v7"
)

func SaveNews(es *elastic.Client, news []byte, hash string) error {
	old, _ := GetBy(es, "ID", hash)
	if len(old) > 0 {
		return fmt.Errorf("this news is already in database")
	}

	ctx := context.Background()
	_, err := es.Index().
		Index(config.IndexName).
		Id(hash).
		Type("_doc").
		BodyJson(string(news)).
		Do(ctx)

	if err != nil {
		return err
	}
	return nil
}

func GetAll(es *elastic.Client) {
	searchSource := elastic.NewSearchSource()
	searchSource.Query(elastic.NewMatchAllQuery())

	searchService := es.Search().Index(config.IndexName).SearchSource(searchSource)

	searchResult, err := searchService.Do(context.Background())
	if err != nil {
		fmt.Println("[ProductsES][GetPIds]Error=", err)
		return
	}
	var news []model.News
	for _, hit := range searchResult.Hits.Hits {
		var n model.News
		err := json.Unmarshal(hit.Source, &n)
		if err != nil {
			fmt.Println("[Getting News][Unmarshal] Err=", err)
		}

		news = append(news, n)
	}

	for _, s := range news {
		fmt.Printf("%v \n", s)
	}

}

func GetBy(es *elastic.Client, field string, value interface{}) ([]model.News, error) {
	searchSource := elastic.NewSearchSource()
	searchSource.Query(elastic.NewMatchQuery(field, value))

	searchService := es.Search().Index(config.IndexName).SearchSource(searchSource)

	searchResult, err := searchService.Do(context.Background())
	if err != nil {
		return nil, fmt.Errorf("[ProductsES][GetPIds]Error=%s", err.Error())
	}
	var news []model.News
	for _, hit := range searchResult.Hits.Hits {
		var n model.News
		err := json.Unmarshal(hit.Source, &n)
		if err != nil {
			return nil, fmt.Errorf("[Getting News][Unmarshal] Err=%s", err.Error())
		}
		news = append(news, n)
	}

	return news, nil
}

func GetAggregate(es *elastic.Client, field string) {
	search := es.Search().
		Index(config.IndexName).
		Query(elastic.NewMatchAllQuery()).
		Aggregation("aggregation", elastic.NewDateHistogramAggregation().Field("PublishedAt").CalendarInterval("day")).
		Sort("PublishedAt", true).
		Pretty(true)

	searchResult, err := search.Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}

	var news []model.News
	if searchResult.Hits == nil {
		return
	}
	for _, hit := range searchResult.Hits.Hits {
		var n model.News
		err := json.Unmarshal(hit.Source, &n)
		if err != nil {
			fmt.Println("[Getting News][Unmarshal] Err=", err)
		}

		news = append(news, n)
	}

	for _, s := range news {
		fmt.Printf("%v \n", s)
	}
}
