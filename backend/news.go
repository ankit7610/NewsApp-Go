package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type Article struct {
	ID      int       `json:"id"`
	Title   string    `json:"title"`
	Summary string    `json:"summary"`
	Author  string    `json:"author"`
	Date    time.Time `json:"date"`
	URL     string    `json:"url"`
	Image   string    `json:"image"`
}

type FinnhubArticle struct {
	Category string `json:"category"`
	Datetime int64  `json:"datetime"`
	Headline string `json:"headline"`
	ID       int    `json:"id"`
	Image    string `json:"image"`
	Related  string `json:"related"`
	Source   string `json:"source"`
	Summary  string `json:"summary"`
	URL      string `json:"url"`
}

func fetchNews(category string) ([]Article, error) {
	// Check cache first
	if cached := getCached(category); cached != nil {
		return cached, nil
	}

	if useSampleData() {
		articles := sampleArticles()
		setCached(category, articles)
		return articles, nil
	}

	token := os.Getenv("FINNHUB_API_KEY")
	if token == "" {
		return nil, fmt.Errorf("missing FINNHUB_API_KEY environment variable")
	}

	base := os.Getenv("FINNHUB_BASE_URL")
	if base == "" {
		base = "https://finnhub.io"
	}
	url := fmt.Sprintf("%s/api/v1/news?category=%s&token=%s", base, category, token)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch news: %s", resp.Status)
	}

	var finnhubArticles []FinnhubArticle
	if err := json.NewDecoder(resp.Body).Decode(&finnhubArticles); err != nil {
		return nil, err
	}

	var articles []Article
	for _, fa := range finnhubArticles {
		articles = append(articles, Article{
			ID:      fa.ID,
			Title:   fa.Headline,
			Summary: fa.Summary,
			Author:  fa.Source,
			Date:    time.Unix(fa.Datetime, 0),
			URL:     fa.URL,
			Image:   fa.Image,
		})
	}

	// store in cache
	setCached(category, articles)
	return articles, nil
}

func useSampleData() bool {
	if v := strings.ToLower(os.Getenv("USE_SAMPLE_DATA")); v == "1" || v == "true" || v == "yes" {
		return true
	}
	return os.Getenv("CI") == "true" && os.Getenv("FINNHUB_API_KEY") == ""
}

// Simple in-memory cache for fetched articles (per category)
var (
	cacheMu sync.RWMutex
	cache   = map[string]cachedEntry{}
	// TTL for cache entries (1 hour)
	cacheTTL = 1 * time.Hour
)

type cachedEntry struct {
	ts       time.Time
	articles []Article
}

func getCached(category string) []Article {
	cacheMu.RLock()
	defer cacheMu.RUnlock()
	e, ok := cache[category]
	if !ok {
		return nil
	}
	if time.Since(e.ts) > cacheTTL {
		return nil
	}
	return e.articles
}

func setCached(category string, articles []Article) {
	cacheMu.Lock()
	defer cacheMu.Unlock()
	cache[category] = cachedEntry{ts: time.Now(), articles: articles}
}

func sampleArticles() []Article {
	now := time.Now()
	return []Article{
		{ID: 1, Title: "Go 1.20 Released", Summary: "Go 1.20 includes many improvements.", Author: "Go Team", Date: now.AddDate(0, 0, -2), URL: "https://golang.org"},
		{ID: 2, Title: "TypeScript 5.0 Announced", Summary: "New TS features for better DX.", Author: "TS Team", Date: now.AddDate(0, 0, -7), URL: "https://www.typescriptlang.org"},
		{ID: 3, Title: "Vite for Fast Frontends", Summary: "Vite continues to lead modern tooling.", Author: "Frontend Weekly", Date: now.AddDate(0, 0, -1), URL: "https://vitejs.dev"},
	}
}
