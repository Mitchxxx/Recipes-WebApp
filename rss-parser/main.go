package main

import (
	"context"
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/html"
)

var client *mongo.Client
var ctx context.Context

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

type Request struct {
	URL string `json:"url"`
}

func GetFeedEntries(url string) ([]Entry, error) {
	client := &http.Client{}
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


func ParseHandler(c *gin.Context) {
    var request Request
    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    entries, err := GetFeedEntries(request.URL)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    if len(entries) == 0 {
        c.JSON(http.StatusOK, gin.H{"message": "no entries parsed"})
        return
    }

    dbName := os.Getenv("MONGO_DATABASE")
    if strings.TrimSpace(dbName) == "" {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "MONGO_DATABASE not set"})
        return
    }

    // derive context with timeout from request
    reqCtx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
    defer cancel()

    col := client.Database(dbName).Collection("recipes_rss")

    // Optional: ensure unique index on url to avoid duplicates
    _, _ = col.Indexes().CreateOne(reqCtx, mongo.IndexModel{
        Keys:    bson.D{{Key: "url", Value: 1}},
        Options: options.Index().SetUnique(true),
    })

    // Don’t skip entries unless you really want to
    inserted := 0
    for i, entry := range entries {
        doc := bson.M{
            "title":     entry.Title,
            "thumbnail": entry.Thumbnail.URL,
            "url":       entry.Link.Href,
            "createdAt": time.Now(),
        }
        if _, err := col.InsertOne(reqCtx, doc); err != nil {
            // Log and continue so you see why a doc failed
            log.Printf("insert failed (i=%d url=%s): %v", i, entry.Link.Href, err)
            continue
        }
        inserted++
    }

    c.JSON(http.StatusOK, gin.H{
        "parsed":   len(entries),
        "inserted": inserted,
        "entries":  entries, // or return a DTO if large
    })
}


func init() {
	ctx = context.Background()
	// Connect to MongoDB
	uri := os.Getenv("MONGO_URI")
	if uri == ""{
		panic("MONGO_URI not set")
	}
	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	if err = client.Ping(ctx, nil); err != nil {
		panic(err)
	}
	log.Println("Connected to MongoDB")
}


func main() {
	router := gin.Default()
	router.POST("/parse", ParseHandler)
	router.Run(":5050")
}