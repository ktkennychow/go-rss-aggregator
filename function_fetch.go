package main

import (
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/ktkennychow/go-rss-aggregator/internal/database"
)

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Version string  `xml:"version,attr"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Title         string  `xml:"title"`
	Link          Link  `xml:"link"`
	Description   string  `xml:"description"`
	Generator     string  `xml:"generator"`
	Language      string  `xml:"language"`
	LastBuildDate string  `xml:"lastBuildDate"`
	Items         []Item  `xml:"item"`
}

type Link struct {
	Rel      string `xml:"rel,attr,omitempty"`
	Href     string `xml:"href,attr"`
	Type     string `xml:"type,attr,omitempty"`
	HrefLang string `xml:"hreflang,attr,omitempty"`
	Title    string `xml:"title,attr,omitempty"`
	Length   uint   `xml:"length,attr,omitempty"`
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	PubDate     string `xml:"pubDate"`
	Guid        string `xml:"guid"`
	Description string `xml:"description"`
}

type fetchResult struct {
	feed RSS
	err  error
}

func (cfg *apiConfig) fetchAFeed(url string) (RSS, error) {
	var rssFeed RSS
	
	resp, err := http.Get(url)
	if err != nil {
		log.Println("faild to fetch RSS feed", err)
		return rssFeed, err
	}
	defer resp.Body.Close()

	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("faild to read RSS feed response body", err)
		return rssFeed, err
	}
	
	err = xml.Unmarshal(dat, &rssFeed)
	if err != nil {
		log.Println("faild to unmarshal RSS feed", err)
		return rssFeed, err
	}

	return rssFeed, nil
}

func (cfg *apiConfig) fetchNFeedsContinuously(n int32, interval int) {
	for {
		urlToFeedID := map[string]uuid.UUID{}

		nFeeds, err := cfg.Queries.ReadNFeedsByLastFetchedAt(cfg.Ctx, n)
		if err != nil {
			log.Println(err.Error())
		}
		wg := sync.WaitGroup{}
		resultChan := make(chan fetchResult)
	
		for _, feed := range nFeeds {
			urlToFeedID[feed.Url] = feed.ID
			wg.Add(1)
			
			go func(url string) {
				defer wg.Done()
	
				rssFeed, err := cfg.fetchAFeed(url)
				if err != nil {
					log.Println(err.Error())
					return
				}
				resultChan <- fetchResult{rssFeed, err}
			} (feed.Url)
		}
	
		go func() {
			wg.Wait()
			close(resultChan)
		}()
		
		for result := range resultChan {
			if result.err != nil {
				log.Println("Error fetching RSS feed", result.err)
				} else {
					for _, post := range result.feed.Channel.Items {
						parsedTime, err := time.Parse("Mon, 02 Jan 2006 15:04:05 -0700", post.PubDate)
						if err != nil {
							log.Println("Error parsing date:", err)
						}
						createdPost, err := cfg.Queries.CreatePost(cfg.Ctx, database.CreatePostParams{
							ID: uuid.New(),
							Title: post.Title,
							Url: post.Link,
							Description: post.Description,
							PublishedAt: parsedTime,
							FeedID: urlToFeedID[result.feed.Channel.Link.Href],
						})
						if err != nil {
							if !strings.Contains(err.Error(), "unique_url_published_at") {
								log.Println("Error creating post in db:", err)
							}
							} else {
							log.Println("Saved a New/Updated Post: ", createdPost.Title)
						}
					}
				}
		}
	
		log.Println("All feeds fetched!")
		time.Sleep(time.Duration(interval) * time.Second)
	}
}