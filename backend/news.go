package main

import "time"

type Article struct {
	ID      int       `json:"id"`
	Title   string    `json:"title"`
	Summary string    `json:"summary"`
	Author  string    `json:"author"`
	Date    time.Time `json:"date"`
	URL     string    `json:"url"`
}

func sampleArticles() []Article {
	now := time.Now()
	return []Article{
		{ID: 1, Title: "Go 1.20 Released", Summary: "Go 1.20 includes many improvements.", Author: "Go Team", Date: now.AddDate(0, 0, -2), URL: "https://golang.org"},
		{ID: 2, Title: "TypeScript 5.0 Announced", Summary: "New TS features for better DX.", Author: "TS Team", Date: now.AddDate(0, 0, -7), URL: "https://www.typescriptlang.org"},
		{ID: 3, Title: "Vite for Fast Frontends", Summary: "Vite continues to lead modern tooling.", Author: "Frontend Weekly", Date: now.AddDate(0, 0, -1), URL: "https://vitejs.dev"},
	}
}
