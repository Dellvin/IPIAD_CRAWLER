package parser

import (
	"IPIAD_DZ/internal/model"
	"fmt"
	"github.com/olivere/elastic/v7"
	"github.com/streadway/amqp"
	"log"
	"net/http"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

func Process(es *elastic.Client) ([]model.News, error) {
	res, err := http.Get("https://upravavernadskogo.ru")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var news []model.News
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		news, err = getAsyncNews(es)
		wg.Done()
	}()

	// Find the review items
	doc.Find("#content > div.blog-post.cat-nature > div > article").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title
		link, ok := s.Find("a").Attr("href")
		fmt.Println("-------------send: ", link)
		if ok {
			send(link)
		}
	})

	send("stop")

	wg.Wait()

	return news, err
}

func createConnection() (amqp.Queue, *amqp.Channel) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")
	return q, ch
}

func send(link string) {
	q, ch := createConnection()

	body := link
	err := ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
