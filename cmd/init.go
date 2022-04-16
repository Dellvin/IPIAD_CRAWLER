package main

import (
	"IPIAD_DZ/config"
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
)

func CreateConnection() (*elastic.Client, error) {
	client, err := elastic.NewClient(elastic.SetURL("http://elasticsearch:9200"),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false))
	if err != nil {
		return nil, err
	}

	exists, err := client.IndexExists(config.IndexName).Do(context.Background())
	if err != nil {

		return nil, err

	} else if exists {
		return client, nil
	}

	_, err = client.CreateIndex(config.IndexName).Body(config.Mapping).Do(context.Background())
	if err != nil {
		return nil, err
	}
	fmt.Println("ES initialized...")

	return client, err

}
