package parser

import (
	"IPIAD_DZ/internal/model"
	"IPIAD_DZ/internal/repository"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic/v7"
	"github.com/streadway/amqp"
	"net/http"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func getAsyncNews(es *elastic.Client) ([]model.News, error) {
	var (
		news = make([]model.News, 0)
	)

	wg := &sync.WaitGroup{}

	msgs := recvSetup()

	for m := range msgs {
		if string(m.Body) == "stop" {
			break
		}
		wg.Add(1)
		fmt.Println("++++++++++++++recv: ", string(m.Body))
		go getPageInfo(string(m.Body), wg, es)
	}

	wg.Wait()

	return news, nil
}

func getPageInfo(link string, wg *sync.WaitGroup, es *elastic.Client) {
	defer func() {
		wg.Done()
		fmt.Println("DONE")
	}()

	var news model.News

	res, err := http.Get(link)
	if err != nil {
		news.Error = err
		return
	}

	if res.StatusCode != 200 {
		news.Code = res.StatusCode
		return
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		news.Error = err
		return
	}

	// Find the review items
	news.Header = doc.Find("#vse-novosti > div > h2").Text()
	news.Body = doc.Find("#vse-novosti > div > article > div").Text()
	date := doc.Find("#vse-novosti > div > article > header > div > div > time").Text()
	news.Author = doc.Find("#vse-novosti > div > article > div > div").Text()
	news.Link = link
	news.Code = res.StatusCode
	hash := md5.Sum([]byte(news.Body))
	news.ID = hex.EncodeToString(hash[:])
	layout := "02.01.2006"
	news.PublishedAt, err = time.Parse(layout, date)
	if err != nil {
		return
	}

	body, err := json.Marshal(news)
	if err != nil {
		fmt.Println("ERROR parsing json: ", err)
		return
	}

	status := repository.SaveNews(es, body, news.ID)
	if status == nil {
		fmt.Println("OK")
	} else {
		fmt.Println(status)
	}
}

func recvSetup() <-chan amqp.Delivery {
	q, ch := createConnection()
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	return msgs
}
