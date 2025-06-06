package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type URLRequest struct {
	URLs []string `json:"urls"`
}

type FetchResult struct {
	Results map[string]int `json:"results"`
	Errors  map[string]string `json:"errors"`
}

func fetchURL(ctx context.Context, url string, results chan<- map[string]int, errors chan<- map[string]string) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		errors <- map[string]string{url: fmt.Sprintf("request creation failed: %v", err)}
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		errors <- map[string]string{url: fmt.Sprintf("request failed: %v", err)}
		return
	}
	defer resp.Body.Close()
	results <- map[string]int{url: resp.StatusCode}
}

func fetchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var reqBody URLRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}
	if len(reqBody.URLs) == 0 {
		http.Error(w, "No URLs provided", http.StatusBadRequest)
		return
	}

	var wg sync.WaitGroup
	results := make(chan map[string]int, len(reqBody.URLs))
	errors := make(chan map[string]string, len(reqBody.URLs))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for _, url := range reqBody.URLs {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			fetchURL(ctx, url, results, errors)
		}(url)
	}

	wg.Wait()
	close(results)
	close(errors)

	fetchResult := FetchResult{
		Results: make(map[string]int),
		Errors:  make(map[string]string),
	}

	for result := range results {
		for url, status := range result {
			fetchResult.Results[url] = status
		}
	}
	for err := range errors {
		for url, msg := range err {
			fetchResult.Errors[url] = msg
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(fetchResult); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func main() {
	http.HandleFunc("/fetch", fetchHandler)
	log.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}