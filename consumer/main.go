package main

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"html"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Request struct {
	URL string `json:"url" binding:"required"`
}

type Feed struct {
	Entries []Entry `xml:"entry"`
}

type Entry struct {
	Link struct {
		Href string `xml:"href,attr"`
	} `xml:"link"`
	Thumbnail struct {
		URL string `xml:"url,attr"`
	} `xml:"thumbnail"`
	Title string `xml:"title"`
	Content struct {
		Type string `xml:"type,attr"`
		Body string `xml:",innerxml"`
	} `xml:"content"`
}

var mongoClient *mongo.Client


func GetFeedEntries(url string) ([]Entry, error) {
	client := &http.Client{Timeout: 15 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36 Edg/138.0.0.0")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	byteValue, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	
	var feed Feed
	if err := xml.Unmarshal(byteValue, &feed); err != nil {
		return nil, err
	}

	// Extract thumbnail URL from content if missing
	imgTagRe := regexp.MustCompile(`(?i)<img[^>]+src=["']([^"']+)["']`)
	// fallback: any http(s) URL ending with common image extensions
	imgLinkRe := regexp.MustCompile(`(?i)https?://[^\s"'<>]+?\.(?:jpg|jpeg|png|gif|webp)`)

	for i := range feed.Entries {
		// if Thumnail already present, keep it
		if strings.TrimSpace(feed.Entries[i].Thumbnail.URL) != "" {
			continue
		}
		// decode HTML entities then search
		content := html.UnescapeString(feed.Entries[i].Content.Body)

		if m := imgTagRe.FindStringSubmatch(content); len(m) == 2{
			feed.Entries[i].Thumbnail.URL = m[1]
			continue
		}
		if m := imgLinkRe.FindString(content); m != ""{
			feed.Entries[i].Thumbnail.URL = m
		}

	}

	return feed.Entries, nil
}

func setEnv(key string) string {
	v := os.Getenv(key)
	if strings.TrimSpace(v) == "" {
		log.Fatalf("%s is not set", key)
	}
	return v
}

func main () {

	// Mongo connection
	mongoURI := setEnv("MONGO_URI")
	dbName := setEnv("MONGO_DATABASE")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	mongoClient, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}
	if err := mongoClient.Ping(ctx, nil); err != nil {
		log.Fatal(err)
	}

	collection := mongoClient.Database(dbName).Collection("recipes_rss")
	// ensure unique index on url once
	if _, err := collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:	bson.D{{Key: "url", Value: 1}},
		Options: options.Index().SetUnique(true),
	}); err != nil {
		log.Fatal("create index:", err)
	}

	log.Println("Mongo connected and index ensured")

	// RabbitMQ connection
	amqpConnection, err := amqp.Dial(os.Getenv("RABBITMQ_URI"))
	if err != nil {
		panic(err)
	}
	defer amqpConnection.Close()
	//
	channelAmqp, _ := amqpConnection.Channel()
	defer channelAmqp.Close()
	

	forever := make(chan bool)
	queue := setEnv("RABBITMQ_QUEUE")
	msgs, err := channelAmqp.Consume(
		queue,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil { log.Fatal("amqp consume:", err)}
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)

			var request Request
			json.Unmarshal(d.Body, &request)

			log.Println("RSS URL:", request.URL)

			// Fetch entries
			entries, err := GetFeedEntries(request.URL)
			if err != nil {
				log.Printf("GetFeedEntries error: %v", err)
				_ = d.Nack(false, true) // transient error, requeue
				continue
			}
			if len(entries) == 0 {
				// Nothing to do for this message
				_= d.Ack(false)
				continue
			}

			// Per-message context for DB writes
      		writeCtx, _ := context.WithTimeout(context.Background(), 10*time.Second)
			// Don't skip entries unless you really want
			inserted := 0
			for _, entry := range entries {
				doc := bson.M{
					"title": entry.Title,
					"thumbnail": entry.Thumbnail.URL,
					"url": entry.Link.Href,
					"createdAt": time.Now(),
				}
				if _, err := collection.InsertOne(writeCtx, doc); err != nil {
					continue
				}
				inserted++
			}

		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}