package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestFetchNews_Success(t *testing.T) {
	// Mock Finnhub server
	handler := func(w http.ResponseWriter, r *http.Request) {
		articles := []FinnhubArticle{
			{Category: "general", Datetime: time.Now().Unix(), Headline: "Test Headline", ID: 10, Image: "", Related: "", Source: "UnitTest", Summary: "Summary", URL: "https://example.com"},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(articles)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	// Set env to point to mock server
	os.Setenv("FINNHUB_API_KEY", "dummy")
	os.Setenv("FINNHUB_BASE_URL", server.URL)

	// Ensure cache is empty
	cacheMu.Lock()
	cache = map[string]cachedEntry{}
	cacheMu.Unlock()

	articles, err := fetchNews("general")
	if err != nil {
		t.Fatalf("fetchNews error: %v", err)
	}
	if len(articles) != 1 {
		t.Fatalf("expected 1 article, got %d", len(articles))
	}
	if articles[0].Title != "Test Headline" {
		t.Fatalf("unexpected title: %s", articles[0].Title)
	}
}

func TestArticlesHandler_OK(t *testing.T) {
	// Mock fetchNews by pointing to test server as above
	handler := func(w http.ResponseWriter, r *http.Request) {
		articles := []FinnhubArticle{
			{Category: "general", Datetime: time.Now().Unix(), Headline: "Handler Headline", ID: 11, Image: "", Related: "", Source: "UnitTest", Summary: "Summary", URL: "https://example.com"},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(articles)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	os.Setenv("FINNHUB_API_KEY", "dummy")
	os.Setenv("FINNHUB_BASE_URL", server.URL)

	// Ensure cache is empty so handler fetches from mock server
	cacheMu.Lock()
	cache = map[string]cachedEntry{}
	cacheMu.Unlock()

	req := httptest.NewRequest("GET", "/api/articles", nil)
	rw := httptest.NewRecorder()
	articlesHandler(rw, req)
	res := rw.Result()
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", res.StatusCode)
	}
	var got []Article
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(got) != 1 || got[0].Title != "Handler Headline" {
		t.Fatalf("unexpected response: %+v", got)
	}
}

func TestCacheExpiry(t *testing.T) {
	// set small TTL
	oldTTL := cacheTTL
	cacheTTL = 500 * time.Millisecond
	defer func() { cacheTTL = oldTTL }()

	// prime cache
	articles := []Article{{ID: 1, Title: "C1", Summary: "", Author: "a", Date: time.Now(), URL: ""}}
	setCached("general", articles)
	c := getCached("general")
	if !reflect.DeepEqual(c, articles) {
		t.Fatalf("cached mismatch before expiry")
	}
	// wait for expiry
	time.Sleep(600 * time.Millisecond)
	c2 := getCached("general")
	if c2 != nil {
		t.Fatalf("expected cache to expire")
	}
}
