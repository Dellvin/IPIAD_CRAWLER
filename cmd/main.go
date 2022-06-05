package main

import (
	"IPIAD_DZ/internal/parser"
	"IPIAD_DZ/internal/repository"
	"fmt"
	"github.com/rai-project/go-fasttext"
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

	news, err := repository.GetAll(es)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	m := fasttext.Open("cc.en.300.bin")
	if m == nil {
		log.Fatalf("file not found")
	}

	a, b, err := calcClusters(m)
	if err != nil {
		log.Fatalf(err.Error())
	}

	var aCounter, bCounter int
	for _, n := range news {
		cluster, e := chooseCluster(a, b, n.Body, m)
		if e != nil {
			continue
		}
		if cluster {
			aCounter++
		} else {
			bCounter++
		}
	}

	fmt.Printf("Add %d news to cluster A\n", aCounter)
	fmt.Printf("Add %d news to cluster B\n", bCounter)
}
