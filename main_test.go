package main

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestFetchHandler(t *testing.T) {
    reqBody := []byte(`{"urls": ["https://httpbin.org/get"]}`)
    req := httptest.NewRequest(http.MethodPost, "/fetch", bytes.NewBuffer(reqBody))
    req.Header.Set("Content-Type", "application/json")

    rr := httptest.NewRecorder()
    fetchHandler(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    var result FetchResult
    if err := json.NewDecoder(rr.Body).Decode(&result); err != nil {
        t.Errorf("failed to decode response: %v", err)
    }

    if _, ok := result.Results["https://httpbin.org/get"]; !ok {
        t.Errorf("expected result for https://httpbin.org/get, got none")
    }
}