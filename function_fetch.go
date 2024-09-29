package main

import (
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Title         string  `xml:"title"`
	Link          string  `xml:"link"`
	Description   string  `xml:"description"`
	Generator     string  `xml:"generator"`
	Language      string  `xml:"language"`
	LastBuildDate string  `xml:"lastBuildDate"`
	AtomLink      AtomLink `xml:"http://www.w3.org/2005/Atom link"` // Namespaced element
	Items         []Item  `xml:"item"`
}

type AtomLink struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr"`
	Type string `xml:"type,attr"`
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
		nFeeds, err :=cfg.Queries.ReadNFeedsByLastFetchedAt(cfg.Ctx, n)
		if err != nil {
			log.Println(err.Error())
		}
		wg := sync.WaitGroup{}
		resultChan := make(chan fetchResult)
	
		for _, feed := range nFeeds {
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
						log.Println("Fetched Post Title: ", post.Title)
					}
				}
		}
	
		log.Println("All feeds fetched!")
		time.Sleep(time.Duration(interval) * time.Second)
	}
}