package main

import (
	"IPIAD_DZ/internal/parser"
	"IPIAD_DZ/internal/repository"
	"fmt"
	"log"
)

func main() {
	es, err := CreateConnection()
	if err != nil {
		log.Fatalf("connection error: %s", err.Error())
	}
	_, err = parser.Process(es)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\"        GET ALL ITEMS       \"\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n\n\n")
	repository.GetAll(es)

	fmt.Println("\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\"        GET ITEMS BY(link)      \"\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n\n\n")
	news, err := repository.GetBy(es, "Link", "https://upravavernadskogo.ru/sotrudniki-detskoj-biblioteki-215-posetili-biblioteku-pri-rospatente")
	if err == nil {
		for _, n := range news {
			fmt.Printf("%v\n", n)
		}
	}

	fmt.Println("\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\"        GET AGGREGATED(PublishedAt) ITEMS       \"\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\n\n\n")

	repository.GetAggregate(es, "PublishedAt")

}
