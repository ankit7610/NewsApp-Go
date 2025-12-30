package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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
	token := os.Getenv("FINNHUB_API_KEY")
	if token == "" {
		token = "d50frm1r01qsabpt5oc0d50frm1r01qsabpt5ocg"
	}
	url := fmt.Sprintf("https://finnhub.io/api/v1/news?category=%s&token=%s", category, token)
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
	return articles, nil
}

func sampleArticles() []Article {
	now := time.Now()
	return []Article{
		{ID: 1, Title: "Go 1.20 Released", Summary: "Go 1.20 includes many improvements.", Author: "Go Team", Date: now.AddDate(0, 0, -2), URL: "https://golang.org"},
		{ID: 2, Title: "TypeScript 5.0 Announced", Summary: "New TS features for better DX.", Author: "TS Team", Date: now.AddDate(0, 0, -7), URL: "https://www.typescriptlang.org"},
		{ID: 3, Title: "Vite for Fast Frontends", Summary: "Vite continues to lead modern tooling.", Author: "Frontend Weekly", Date: now.AddDate(0, 0, -1), URL: "https://vitejs.dev"},
	}
}
